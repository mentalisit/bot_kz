package config

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type SimpleHandler struct {
	writers     []Writer
	prefix      string
	skipCallers int
}

type Writer interface {
	Write(message string) error
}

func NewSimpleHandler(prefix string, writers ...Writer) *SimpleHandler {
	return &SimpleHandler{
		prefix:      fmt.Sprintf("[%s]", prefix),
		writers:     writers,
		skipCallers: 3,
	}
}
func (h *SimpleHandler) SetSkipCallers(s int) {
	h.skipCallers = s
}

func (h *SimpleHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (h *SimpleHandler) Handle(ctx context.Context, record slog.Record) error {
	// Форматируем ОДИН раз
	timestamp := record.Time.Format("2006-01-02 15:04:05")

	var filePrefix string
	var line int

	if record.PC != 0 {
		frames := runtime.CallersFrames([]uintptr{record.PC})
		frame, _ := frames.Next()
		filePrefix = filepath.Base(filepath.Dir(frame.File)) + "/" + filepath.Base(frame.File)
		line = frame.Line
	} else {
		filePrefix = "unknown"
		line = 0
	}

	message := fmt.Sprintf("%s %s %s %s:%d %s\n",
		timestamp,
		record.Level.String(),
		h.prefix,
		filePrefix,
		line,
		record.Message,
	)
	// Если уровень логирования - ошибка, добавляем стек вызовов
	if record.Level == slog.LevelError {
		stackTrace := h.getStackTrace()
		message += stackTrace
	}

	// Отправляем ВСЕМ writer-ам
	for _, writer := range h.writers {
		_ = writer.Write(message)
	}

	return nil
}

func (h *SimpleHandler) getStackTrace() string {
	// Указываем, сколько элементов стека захватывать (например, 20)
	pcs := make([]uintptr, 20)
	n := runtime.Callers(h.skipCallers, pcs) // Пропускаем первые два элемента, чтобы не включать текущие вызовы

	frames := runtime.CallersFrames(pcs[:n])
	var result []string

	// Перебираем кадры стека и фильтруем лишние
	for {
		frame, more := frames.Next()
		// Пропускаем кадры из стандартной библиотеки
		if !strings.Contains(frame.File, "runtime") && !strings.Contains(frame.File, "log/slog") {
			result = append(result, fmt.Sprintf("%s:%d %s\n", frame.File, frame.Line, frame.Function))
		}
		if !more {
			break
		}
	}

	return strings.Join(result, "\n")
}

func (h *SimpleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *SimpleHandler) WithGroup(name string) slog.Handler {
	return h
}

type StdoutWriter struct{}

func (w *StdoutWriter) Write(message string) error {
	_, err := fmt.Fprintln(os.Stdout, message)
	return err
}

type DiscordWriter struct {
	WebhookURL string
}

func (w *DiscordWriter) Write(message string) error {
	if w.WebhookURL == "" {
		return errors.New("webhook url is empty")
	}
	go sendDiscordMessage(w.WebhookURL, message)
	return nil
}

func sendDiscordMessage(webhookURL, text string) {
	message := map[string]interface{}{
		"content":    fmt.Sprintf("```%s```", text),
		"username":   "Logger",
		"avatar_url": "https://e7.pngegg.com/pngimages/836/966/png-clipart-go-programming-language-computer-programming-others-baltimore-web-application-thumbnail.png",
	}

	jsonBody, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println(err)
	}
	_ = resp.Body.Close()
}

type TelegramWriter struct {
	BotToken string
	ChatID   string
}

func (w *TelegramWriter) Write(message string) error {
	go sendTelegramMessage(w.BotToken, w.ChatID, message)
	return nil
}

func sendTelegramMessage(token, chatID, message string) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	params := map[string]string{
		"chat_id": chatID,
		"text":    message,
	}

	body := new(bytes.Buffer)
	for key, value := range params {
		body.WriteString(fmt.Sprintf("%s=%s&", key, value))
	}
	bodyStr := body.String()

	resp, err := http.Post(apiURL, "application/x-www-form-urlencoded", bytes.NewBufferString(bodyStr))
	if err != nil {
		fmt.Printf("ошибка при отправке запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("неправильный статус ответа: %s", resp.Status)
	}
}

func InitLogger(prefix, webhook, TgToken string, TgChat int64) *slog.Logger {
	handler := NewSimpleHandler(prefix,
		&StdoutWriter{},
		&DiscordWriter{WebhookURL: webhook},
		&TelegramWriter{BotToken: TgToken, ChatID: strconv.FormatInt(TgChat, 10)},
	)

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}
