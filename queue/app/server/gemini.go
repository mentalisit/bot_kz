package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/mentalisit/logger"
	"google.golang.org/api/option"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	Question string
	Name     string
}

func (s *Server) GeminiAI(c *gin.Context) {
	s.PrintGoroutine()
	var req Request

	// Читаем тело запроса как необработанные JSON-данные
	if err := c.ShouldBindJSON(&req); err != nil {
		s.log.ErrorErr(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Question == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Question is empty"})
	}
	if req.Question != "" && req.Name != "" {
		answer := s.ai.GetAnswerGemini(req.Question, req.Name)
		if len(answer) != 0 {
			c.JSON(http.StatusOK, gin.H{"strings": answer})
		}
	}
}

type Gemini struct {
	log     *logger.Logger
	history map[string][]*genai.Content
}

func NewGemini(log *logger.Logger) *Gemini {
	return &Gemini{log: log, history: make(map[string][]*genai.Content)}
}

func (g *Gemini) GetAnswerGemini(question, name string) []string {
	answerBig := g.getAnswerGemini(question, name)
	if len(answerBig) == 0 {
		return nil
	}
	return g.splitTextByDot(answerBig, 1900)
}
func (g *Gemini) getAnswerGemini(question, name string) (text string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := genai.NewClient(ctx, option.WithAPIKey("AIzaSyAhZ6vmmalmtKWuTOk62Bld9MR87aZYKy0"))
	if err != nil {
		g.log.ErrorErr(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	cs := model.StartChat()

	cs.History = []*genai.Content{
		{
			Parts: []genai.Part{
				genai.Text("Привет меня зовут " + name + " я пишу тебе через бота "),
				genai.Text("КЗ или RS это красная звезда, тКз или DRS это мрачная красная звезда,  Bso или бсо  это одна из топовых корпораций в игре. Это все в игре Hades Star"),
			},
			Role: "user",
		},
		{
			Parts: []genai.Part{
				genai.Text(fmt.Sprintf("привет %s!", name)),
			},
			Role: "model",
		},
		{
			Parts: []genai.Part{
				genai.Text("КЗ или RS это красная звезда, тКз или DRS это мрачная красная звезда,  Bso или бсо  это одна из топовых корпораций в игре. Это все в игре Hades Star"),
			},
			Role: "user",
		},
		{
			Parts: []genai.Part{
				genai.Text("Расскажи, что ты хочешь обсудить в контексте Hades Star.  У меня есть некоторый базовый уровень знаний об игре, но чем больше деталей ты предоставишь, тем лучше я смогу помочь. Хочешь поговорить о стратегии, о твоих текущих задачах в игре, или о чём-то другом?"),
			},
			Role: "model",
		},
		{
			Parts: []genai.Part{
				genai.Text("я подумаю"),
			},
			Role: "user",
		},
		{
			Parts: []genai.Part{
				genai.Text("жду твои вопросы"),
			},
			Role: "model",
		},
	}
	contents, exist := g.history[name]
	if exist {
		cs.History = append(cs.History, contents...)
	}
	g.history[name] = append(cs.History, []*genai.Content{
		{
			Parts: []genai.Part{
				genai.Text(question),
			},
			Role: "user",
		},
	}...)

	resp, err := cs.SendMessage(ctx, genai.Text(question))
	if err != nil {
		g.log.ErrorErr(err)
		return ""
	}
	if resp != nil && len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		g.history[name] = append(g.history[name], resp.Candidates[0].Content)
		text = fmt.Sprintf("%+v\n", resp.Candidates[0].Content.Parts[0])
	}
	return text
}

// SplitTextByDot разделяет текст на части с учетом лимита символов, стараясь разделить на ближайшей точке.
func (g *Gemini) splitTextByDot(text string, limit int) []string {
	var parts []string
	for len(text) > limit {
		// Находим ближайшую точку в пределах лимита
		end := strings.LastIndex(text[:limit+1], ".")
		if end == -1 { // Если точка не найдена, просто режем по лимиту
			end = limit
		} else {
			end++ // Включаем точку в конец текущей части
		}
		// Добавляем часть текста
		parts = append(parts, strings.TrimSpace(text[:end]))
		// Удаляем обработанную часть из текста
		text = strings.TrimSpace(text[end:])
	}
	// Добавляем оставшийся текст
	if len(text) > 0 {
		parts = append(parts, text)
	}
	return parts
}
