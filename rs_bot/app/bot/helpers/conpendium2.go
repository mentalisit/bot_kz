package helpers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// type techLevelArray map[int][2]int
type TechLevels map[int]TechLevel
type TechLevel struct {
	Ts    int64 `json:"ts"`
	Level int   `json:"level"`
}

type corpMember struct {
	Name         string     `json:"name"`
	UserId       string     `json:"userId"`
	GuildId      string     `json:"guildId"`
	Avatar       string     `json:"avatar"`
	Tech         TechLevels `json:"tech"`
	AvatarUrl    string     `json:"avatarUrl"`
	LocalTime    string     `json:"localTime"`
	LocalTime24  string     `json:"localTime24"`
	TimeZone     string     `json:"timeZone"`
	ZoneOffset   int        `json:"zoneOffset"`
	AfkFor       string     `json:"afkFor"`
	AfkWhen      int        `json:"afkWhen"`
	MultiAccount *MultiAccount
	MAcc         *MultiAccount
}
type MultiAccount struct {
	UUID             uuid.UUID
	Nickname         string
	TelegramID       string
	TelegramUsername string
	DiscordID        string
	DiscordUsername  string
	WhatsappID       string
	WhatsappUsername string
	CreatedAt        time.Time
	AvatarURL        string
	Alts             []string
}

func Get2TechDataUserId(name, userID, guildid string) (genesis, enrich, rsextender int) {
	// Создаем контекст с тайм-аутом 2 секунды
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // Освобождаем ресурсы контекста после завершения функции

	// Формирование URL-адреса
	url := fmt.Sprintf("http://compendiumnew/compendium/api?userid=%s&guildid=%s", userID, guildid)

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

	if technicalData[0].MultiAccount != nil {
		ma := technicalData[0].MultiAccount
		if (ma.DiscordID == userID && ma.DiscordUsername == name) ||
			(ma.TelegramID == userID && ma.TelegramUsername == name) ||
			(ma.WhatsappID == userID && ma.WhatsappUsername == name) {
			name = ma.Nickname
		}

	}

	for _, datum := range technicalData {
		if strings.ToLower(datum.Name) == strings.ToLower(name) {
			rsextender = datum.Tech[603].Level
			enrich = datum.Tech[503].Level
			genesis = datum.Tech[508].Level
		}
	}
	if rsextender == 0 && enrich == 0 && genesis == 0 && len(technicalData) != 0 {
		rsextender = technicalData[0].Tech[603].Level
		enrich = technicalData[0].Tech[503].Level
		genesis = technicalData[0].Tech[508].Level
	}

	return
}

func Get3TechDataUserId(name, userIdTg string) (genesis, enrich, rse int) {
	url := fmt.Sprintf("https://123bot.ru/rssoyuzbot/Json/module.php?userid=%s", userIdTg)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var tech []map[string]string

	if err = json.Unmarshal(body, &tech); err != nil {
		return
	}
	if len(tech) == 1 {
		genesis, _ = strconv.Atoi(tech[0]["g"])
		enrich, _ = strconv.Atoi(tech[0]["o"])
		rse, _ = strconv.Atoi(tech[0]["ikz"])
		return genesis, enrich, rse
	} else {
		for _, m := range tech {
			if strings.ToLower(m["nameacc"]) == strings.ToLower(name) {
				genesis, _ = strconv.Atoi(m["g"])
				enrich, _ = strconv.Atoi(m["o"])
				rse, _ = strconv.Atoi(m["ikz"])
				return genesis, enrich, rse
			}
		}
	}
	return
}

func Get2AltsUserId(userID string) (alts []string) {
	// Создаем контекст с тайм-аутом 2 секунды
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // Освобождаем ресурсы контекста после завершения функции

	// Формирование URL-адреса
	url := fmt.Sprintf("http://compendiumnew/compendium/api/user?userid=%s", userID)

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

	// Декодирование JSON-данных в структуру TechnicalData
	err = json.Unmarshal(body, &alts)
	if err != nil {
		fmt.Println(string(body))
		fmt.Println("Ошибка при декодировании JSON:", err)
		return
	}
	return alts
}
