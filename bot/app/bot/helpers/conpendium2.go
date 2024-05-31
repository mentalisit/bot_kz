package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type techLevelArray map[int][2]int

type corpMember struct {
	Name        string         `json:"name"`
	UserId      string         `json:"userId"`
	GuildId     string         `json:"guildId"`
	Avatar      string         `json:"avatar"`
	Tech        techLevelArray `json:"tech"`
	AvatarUrl   string         `json:"avatarUrl"`
	LocalTime   string         `json:"localTime"`
	LocalTime24 string         `json:"localTime24"`
	TimeZone    string         `json:"timeZone"`
	ZoneOffset  int            `json:"zoneOffset"`
	AfkFor      string         `json:"afkFor"`
	AfkWhen     int            `json:"afkWhen"`
}

func Get2TechDataUserId(name, userID, guildid string) (genesis, enrich, rsextender int) {
	apiKey := "gGUBIlUAU1uTKWd8HssP27ojG0DugoAaPslwFGTDSAbEM6UM"

	// Формирование URL-адреса
	url := fmt.Sprintf("https://compendiumnew.mentalisit.myds.me/compendium/api/tech?token=%s&userid=%s&guildid=%s", apiKey, userID, guildid)

	// Выполнение GET-запроса
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
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

	// Декодирование JSON-данных в структуру TechnicalData
	var technicalData []corpMember
	err = json.Unmarshal(body, &technicalData)
	if err != nil {
		fmt.Println("Ошибка при декодировании JSON:", err)
		return
	}
	for _, datum := range technicalData {
		if datum.Name == name {
			rsextender = datum.Tech[603][0]
			enrich = datum.Tech[503][0]
			genesis = datum.Tech[508][0]
		}
	}
	return
}
