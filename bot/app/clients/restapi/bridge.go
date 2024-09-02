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

func SendBridgeApp(m models.ToBridgeMessage) error {
	utils.PrintGoroutine(nil)
	fmt.Printf("Send to bridge: %s %s lenFile:%d Sender: %s Text: %s\n", m.Config.HostRelay, m.ChatId, len(m.Extra), m.Sender, m.Text)
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	// Создаем контекст с тайм-аутом 3 секунды
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Создаем новый запрос с контекстом
	req, err := http.NewRequestWithContext(ctx, "POST", "http://bridge/bridge/inbox", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Создаем клиент
	client := &http.Client{}

	// Канал для отслеживания завершения запроса
	done := make(chan struct{})
	var returnErr error

	// Горутина для выполнения запроса
	go func() {
		resp, err := client.Do(req)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				returnErr = fmt.Errorf("время ожидания запроса истекло")
			} else {
				returnErr = fmt.Errorf("Ошибка при выполнении запроса: %+v\n", err)
			}
		} else {
			defer resp.Body.Close()

			// Проверка кода ответа
			if resp.StatusCode != http.StatusOK {
				returnErr = fmt.Errorf("неправильный статус код: %d", resp.StatusCode)
			}
		}

		// Закрываем канал, чтобы сигнализировать о завершении
		close(done)
	}()

	select {
	case <-done:
		// Запрос завершился до истечения таймаута
	case <-time.After(5 * time.Second):
		// Логируем, если запрос завис
		return fmt.Errorf("SendBridgeApp завис: %+v\n, %+v\n", returnErr, m)
	}

	return nil
}
func ReadBridgeConfig() map[string]models.BridgeConfig {
	var br []models.BridgeConfig
	resp, err := http.Get("http://bridge/bridge/config")
	if err != nil {
		resp, err = http.Get("http://192.168.100.131:808/bridge/config")
		if err != nil {
			return nil
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	err = json.NewDecoder(resp.Body).Decode(&br)
	if err != nil {
		return nil
	}

	var bridgeCounter = 0
	var bridge string

	bridgeConfig := make(map[string]models.BridgeConfig)

	for _, configBridge := range br {
		bridgeConfig[configBridge.NameRelay] = configBridge
		bridgeCounter++
		bridge = bridge + fmt.Sprintf("%s, ", configBridge.HostRelay)
	}
	fmt.Printf("Загружено конфиг мостов %d : %s\n", bridgeCounter, bridge)

	return bridgeConfig
}
