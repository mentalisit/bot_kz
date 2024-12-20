package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

//
//import (
//	"context"
//	"fmt"
//	"github.com/google/generative-ai-go/genai"
//	"github.com/mentalisit/logger"
//	"google.golang.org/api/option"
//	"strings"
//	"time"
//)
//
//type Gemini struct {
//	log     *logger.Logger
//	history map[string][]*genai.Content
//}
//
//func NewGemini(log *logger.Logger) *Gemini {
//	return &Gemini{log: log, history: make(map[string][]*genai.Content)}
//}
//
//func (g *Gemini) GetAnswerGemini(qustion, name string) []string {
//	answerBig := g.getAnswerGemini(qustion, name)
//	if len(answerBig) == 0 {
//		return nil
//	}
//	return g.splitTextByDot(answerBig, 1900)
//}
//func (g *Gemini) getAnswerGemini(qustion, name string) (text string) {
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//	client, err := genai.NewClient(ctx, option.WithAPIKey("AIzaSyAhZ6vmmalmtKWuTOk62Bld9MR87aZYKy0"))
//	if err != nil {
//		g.log.ErrorErr(err)
//	}
//	defer client.Close()
//
//	model := client.GenerativeModel("gemini-1.5-flash")
//	cs := model.StartChat()
//
//	cs.History = []*genai.Content{
//		{
//			Parts: []genai.Part{
//				genai.Text("Привет меня зовут " + name + " я пишу тебе через бота "),
//				genai.Text("КЗ или RS это красная звезда, тКз или DRS это мрачная красная звезда,  Bso или бсо  это одна из топовых корпораций в игре. Это все в игре Hades Star"),
//			},
//			Role: "user",
//		},
//		{
//			Parts: []genai.Part{
//				genai.Text(fmt.Sprintf("привет %s!", name)),
//			},
//			Role: "model",
//		},
//		{
//			Parts: []genai.Part{
//				genai.Text("КЗ или RS это красная звезда, тКз или DRS это мрачная красная звезда,  Bso или бсо  это одна из топовых корпораций в игре. Это все в игре Hades Star"),
//			},
//			Role: "user",
//		},
//		{
//			Parts: []genai.Part{
//				genai.Text("Расскажи, что ты хочешь обсудить в контексте Hades Star.  У меня есть некоторый базовый уровень знаний об игре, но чем больше деталей ты предоставишь, тем лучше я смогу помочь. Хочешь поговорить о стратегии, о твоих текущих задачах в игре, или о чём-то другом?"),
//			},
//			Role: "model",
//		},
//		{
//			Parts: []genai.Part{
//				genai.Text("я подумаю"),
//			},
//			Role: "user",
//		},
//		{
//			Parts: []genai.Part{
//				genai.Text("жду твои вопросы"),
//			},
//			Role: "model",
//		},
//	}
//	contents, exist := g.history[name]
//	if exist {
//		cs.History = append(cs.History, contents...)
//	}
//	g.history[name] = append(cs.History, []*genai.Content{
//		{
//			Parts: []genai.Part{
//				genai.Text(qustion),
//			},
//			Role: "user",
//		},
//	}...)
//
//	resp, err := cs.SendMessage(ctx, genai.Text(qustion))
//	if err != nil {
//		g.log.ErrorErr(err)
//		return ""
//	}
//	if resp != nil && len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
//		g.history[name] = append(g.history[name], resp.Candidates[0].Content)
//		text = fmt.Sprintf("%+v\n", resp.Candidates[0].Content.Parts[0])
//	}
//	return text
//}
//
//// SplitTextByDot разделяет текст на части с учетом лимита символов, стараясь разделить на ближайшей точке.
//func (g *Gemini) splitTextByDot(text string, limit int) []string {
//	var parts []string
//	for len(text) > limit {
//		// Находим ближайшую точку в пределах лимита
//		end := strings.LastIndex(text[:limit+1], ".")
//		if end == -1 { // Если точка не найдена, просто режем по лимиту
//			end = limit
//		} else {
//			end++ // Включаем точку в конец текущей части
//		}
//		// Добавляем часть текста
//		parts = append(parts, strings.TrimSpace(text[:end]))
//		// Удаляем обработанную часть из текста
//		text = strings.TrimSpace(text[end:])
//	}
//	// Добавляем оставшийся текст
//	if len(text) > 0 {
//		parts = append(parts, text)
//	}
//	return parts
//}

type Request struct {
	Question string
	Name     string
}

func (h *Helpers) GeminiSay(qustion, name string) (answer []string) {
	// Создаем контекст с тайм-аутом
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Освобождаем ресурсы контекста после завершения функции

	body := Request{
		Question: qustion,
		Name:     name,
	}

	data, err := json.Marshal(body)
	if err != nil {
		h.log.ErrorErr(err)
		return nil
	}

	// Выполнение GET-запроса с контекстом
	req, err := http.NewRequestWithContext(ctx, "POST", "http://queue/ai", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Println("Время ожидания запроса истекло")
		} else {
			fmt.Println("Ошибка при выполнении запроса:", err)
		}
		return
	}
	defer response.Body.Close()

	// Проверка кода ответа
	if response.StatusCode != http.StatusOK {
		fmt.Printf("Неправильный статус код: %d\n", response.StatusCode)
		return
	}

	// Чтение тела ответа
	readAll, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return
	}

	// Декодирование
	err = json.Unmarshal(readAll, &answer)
	if err != nil {
		fmt.Println("Ошибка при декодировании JSON:", err)
		return
	}
	return
}
