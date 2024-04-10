package corpPercent

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func getKey(key string) *CorporationsData {
	// URL для загрузки JSON-данных
	url := "https://storage.googleapis.com/hades-star-public-xq8f-d4rg-v0d9/" + key
	// Получаем JSON-данные с помощью HTTP GET-запроса
	resp, err := http.Get(url)
	if err != nil {
		log.Error(fmt.Sprintln("Ошибка при загрузке данных:", err))
		return nil
	}
	defer resp.Body.Close()

	// Читаем данные из тела ответа
	jsonData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(fmt.Sprintln("Ошибка при чтении данных:", err))
		return nil
	}

	// Создаем переменную для хранения данных о корпорациях
	var data CorporationsData

	// Распаковываем JSON в структуру данных
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Error(fmt.Sprintln("Ошибка при распаковке JSON:", err))
		return nil
	}

	return data.sort()
}
func (data *CorporationsData) sort() *CorporationsData {
	if data.Corporation1Score > data.Corporation2Score {
		//fmt.Printf("1Score>2Score %d-%d %s-%s\n", data.Corporation1Score, data.Corporation2Score, data.Corporation1Name, data.Corporation2Name)
		return data
	} else if data.Corporation2Score > data.Corporation1Score {
		corp := CorporationsData{
			Corporation1Name:  data.Corporation2Name,
			Corporation1Id:    data.Corporation2Id,
			Corporation2Name:  data.Corporation1Name,
			Corporation2Id:    data.Corporation1Id,
			Corporation1Score: data.Corporation2Score,
			Corporation2Score: data.Corporation1Score,
			DateEnded:         data.DateEnded,
		}
		//fmt.Printf("2Score>1Score %d-%d %s-%s\n", data.Corporation2Score, data.Corporation1Score, data.Corporation2Name, data.Corporation1Name)
		return &corp
	}
	return &CorporationsData{}
}
