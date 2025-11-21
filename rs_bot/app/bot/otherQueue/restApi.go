package otherQueue

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"rs/models"
	"time"
)

func GetQueueLevel(level string) (t map[string][]models.QueueStruct, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel() // Освобождаем ресурсы контекста после завершения функции

	// Формирование URL-адреса
	url := fmt.Sprintf("http://queue/queue?level=%s", level)

	// Выполнение GET-запроса с контекстом
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
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
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return
	}

	err = json.Unmarshal(body, &t)
	if err != nil {
		return
	}
	return
}

func GetQueueAll() (t map[string][]models.QueueStruct, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel() // Освобождаем ресурсы контекста после завершения функции

	// Формирование URL-адреса
	url := "http://queue/queue"

	// Выполнение GET-запроса с контекстом
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
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
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return
	}

	err = json.Unmarshal(body, &t)
	if err != nil {
		return
	}
	return
}

func GetUseridTumcha() ([]int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	url := "http://queue/api/readsouzbot"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return nil, err
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Println("Время ожидания запроса истекло")
		} else {
			fmt.Println("Ошибка при выполнении запроса:", err)
		}
		return nil, err
	}
	defer response.Body.Close()

	// Обрабатываем статусы согласно логике ReadQueueTumcha
	if response.StatusCode == http.StatusNoContent {
		// Сервер вернул пустую очередь - это нормальная ситуация
		return []int64{}, nil
	}

	if response.StatusCode != http.StatusOK {
		fmt.Printf("Неправильный статус код: %d\n", response.StatusCode)
		return nil, fmt.Errorf("HTTP error: %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return nil, err
	}

	var result []int64
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// SendSborkzData отправляет POST запрос на https://123bot.ru/rssoyuzbot/Json/sborkz.php с userid
func SendSborkzData(userid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	url := "https://123bot.ru/rssoyuzbot/Json/sborkz.php"

	// Создание структуры с только userid
	data := map[string]string{
		"userid": userid,
	}

	// Сериализация данных в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("ошибка при сериализации данных: %w", err)
	}

	// Создание POST запроса
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	// Установка заголовков
	req.Header.Set("Content-Type", "application/json")

	// Выполнение запроса
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("время ожидания запроса истекло")
		}
		return fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer response.Body.Close()

	// Проверка кода ответа
	if response.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(response.Body)
		return fmt.Errorf("неправильный статус код: %d, ответ: %s", response.StatusCode, string(body))
	}

	return nil
}
