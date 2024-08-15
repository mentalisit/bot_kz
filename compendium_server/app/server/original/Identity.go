package original

import (
	"compendium_s/models"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func CheckIdentity(code string) (*models.Identity, error) {
	apiURL := "https://bot.hs-compendium.com/compendium/applink/identities?ver=2&code=1"

	// Подготовка запроса
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return &models.Identity{}, err
	}

	// Установка параметров запроса
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", code)

	// Отправка запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &models.Identity{}, err
	}
	defer resp.Body.Close()

	// Проверка успешного ответа
	if resp.StatusCode < 200 || resp.StatusCode >= 500 {
		return &models.Identity{}, errors.New("Server Error")
	}

	// Чтение тела ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &models.Identity{}, err
	}

	if len(body) < 50 {
		return &models.Identity{}, errors.New("Invalid User Id")
	}

	// Декодирование JSON-строки в структуру Identity
	var identity models.IdentityGET
	err = json.Unmarshal(body, &identity)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return &models.Identity{}, err
	}
	fmt.Printf("JSON-строкa %+v\n", string(body))

	return &models.Identity{
		User:  identity.User,
		Guild: identity.Guild[0],
		Token: identity.Token,
	}, nil
}
