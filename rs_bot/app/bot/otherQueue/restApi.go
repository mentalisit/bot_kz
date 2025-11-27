package otherQueue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
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

func SendUserIDRSSOYUZ(userID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	baseURL := "https://123bot.ru/rssoyuzbot/Json/queue.php"

	params := url.Values{}
	params.Add("userid", fmt.Sprintf("%d", userID))

	fullURL := baseURL + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка отправки запроса: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	fmt.Printf("Статус: %d, Ответ: %s\n", resp.StatusCode, string(body))
	return nil
}
