package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"rs/models"
)

func (b *Bot) GameWebhook(in models.InMessage) {
	if in.Username != "fakeData" {
		err := b.FetchJSON(in.Mtext, in.Username)
		if err != nil {
			b.log.ErrorErr(err)
			return
		}
	}
}

func (b *Bot) FetchJSON(url string, corpName string) error {
	// Отправка GET-запроса
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("ошибка при отправке запроса: %v", err)
	}
	defer resp.Body.Close()

	// Проверка статуса ответа
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неожиданный статус ответа: %s", resp.Status)
	}

	// Чтение тела ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка при чтении ответа: %v", err)
	}

	// Декодирование JSON в map для удобного отображения
	var jsonData map[string]interface{}
	if err = json.Unmarshal(body, &jsonData); err != nil {
		return fmt.Errorf("ошибка при парсинге JSON: %v", err)
	}

	// Форматированный вывод JSON
	formattedJSON, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка при форматировании JSON: %v", err)
	}
	switch jsonData["EventType"] {
	case "RedStarStarted":
		var rss models.RedStarEvent
		err = rss.UnmarshalJSON(formattedJSON)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Sprintf("%+v\n", rss)
	case "RedStarEnded":
		var rse models.RedStarEvent
		if err = rse.UnmarshalJSON(formattedJSON); err != nil {
			fmt.Println(err)
		}

		if rse.RSEventPoints != 0 {
			points := rse.RSEventPoints / len(rse.Players)
			for _, player := range rse.Players {
				err = b.storage.Battles.BattlesInsert(models.Battles{
					EventId:  32,
					CorpName: corpName,
					Name:     player.PlayerName,
					Level:    rse.StarLevel,
					Points:   points,
				})
				if err != nil {
					fmt.Println(err)
					return err
				}
			}
		}

		//fmt.Printf("%+v\n", rse)
	default:

		fmt.Println(string(formattedJSON))

	}

	//fmt.Println(len(string(formattedJSON)))
	return nil
}
