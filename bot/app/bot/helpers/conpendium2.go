package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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
	// Формирование URL-адреса
	url := fmt.Sprintf("http://compendiumnew/compendium/api?userid=%s&guildid=%s", userID, guildid)

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
	if len(technicalData) < 1 {
		return
	}
	if len(technicalData) == 1 {
		rsextender = technicalData[0].Tech[603][0]
		enrich = technicalData[0].Tech[503][0]
		genesis = technicalData[0].Tech[508][0]
	}
	for _, datum := range technicalData {
		if strings.ToLower(datum.Name) == strings.ToLower(name) {
			rsextender = datum.Tech[603][0]
			enrich = datum.Tech[503][0]
			genesis = datum.Tech[508][0]
		}
	}
	return
}

func Get2AltsUserId(userID string) (alts []string) {
	// Формирование URL-адреса
	url := fmt.Sprintf("http://compendiumnew/compendium/api/user?userid=%s", userID)

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
	err = json.Unmarshal(body, &alts)
	if err != nil {
		fmt.Println("Ошибка при декодировании JSON:", err)
		return
	}
	return alts
}
