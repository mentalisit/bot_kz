package restapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"kz_bot/models"
	"kz_bot/pkg/utils"
	"net/http"
	"time"
)

func SendCompendiumApp(m models.IncomingMessage) error {
	utils.PrintGoroutine(nil)
	fmt.Printf("Sending compendium app %s %s %s\n", m.GuildName, m.Name, m.Text)
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	// Создаем контекст с тайм-аутом 3 секунды
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel() // Освобождаем ресурсы контекста после завершения функции

	// Создаем новый запрос с контекстом
	req, err := http.NewRequestWithContext(ctx, "POST", "http://compendiumnew/compendium/inbox", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Создаем клиент с тайм-аутом
	client := &http.Client{}

	// Канал для отслеживания завершения запроса
	done := make(chan struct{})
	var returnErr error

	go func() {
		// Выполняем запрос
		resp, err := client.Do(req)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				returnErr = fmt.Errorf("время ожидания запроса истекло")
			} else {
				returnErr = fmt.Errorf("Ошибка при выполнении запроса: %+v\n", err)
			}
			return
		}
		defer resp.Body.Close()

		// Проверка кода ответа
		if resp.StatusCode != http.StatusOK {
			returnErr = fmt.Errorf("неправильный статус код: %d", resp.StatusCode)
		}

		close(done)
	}()

	if returnErr != nil {
		return returnErr
	}

	return nil
}
