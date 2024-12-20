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

type Request struct {
	Question string
	Name     string
}
type RequestData struct {
	Strings []string `json:"strings"`
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

	var reqData RequestData
	// Декодирование
	err = json.Unmarshal(readAll, &reqData)
	if err != nil {
		fmt.Println("Ошибка при декодировании JSON:", err)
		return
	}

	return reqData.Strings
}
