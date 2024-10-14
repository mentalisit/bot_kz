package hspublic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"ws/models"
)

func (h *HS) GetCorporationsData(key string) *models.CorporationsData {
	// URL для загрузки JSON-данных
	url := "https://storage.googleapis.com/hades-star-public-xq8f-d4rg-v0d9/" + key
	// Получаем JSON-данные с помощью HTTP GET-запроса
	resp, err := http.Get(url)
	if err != nil {
		h.log.Error(fmt.Sprintln("Ошибка при загрузке данных:", err))
		return nil
	}
	defer resp.Body.Close()

	// Читаем данные из тела ответа
	jsonData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		h.log.Error(fmt.Sprintln("Ошибка при чтении данных:", err))
		return nil
	}

	fmt.Println(key + " " + string(jsonData))

	// Создаем переменную для хранения данных о корпорациях
	var data models.CorporationsData

	// Распаковываем JSON в структуру данных
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		h.log.Error(fmt.Sprintln("Ошибка при распаковке JSON:", err))
		return nil
	}
	data.SortWin()

	return &data
}
