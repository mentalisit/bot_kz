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

	// Выполняем запрос
	resp, err := client.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Println("Время ожидания запроса истекло")
		} else {
			fmt.Println("Ошибка при выполнении запроса:", err)
		}
		return err
	}
	defer resp.Body.Close()

	// Проверка кода ответа
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неправильный статус код: %d", resp.StatusCode)
	}

	return nil
}
