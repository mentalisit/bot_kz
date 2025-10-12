package GetDataSborkz

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"queue/models"
	"strconv"
	"time"
)

type Sborkz struct {
	Name      string `json:"name"`
	Userid    string `json:"userid"`
	Chatid    string `json:"chatid"`
	Timestamp string `json:"timestamp"`
	Lvlkz     string `json:"lvlkz"`
	Vid       string `json:"vid"`
	Timedown  string `json:"timedown"`
}

func ReadQueueTumchaNameIds() (namesIds []int64) {
	q, err := fetchDataSborkz()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, tumcha := range q {
		parseInt, _ := strconv.ParseInt(tumcha.Userid, 10, 64)
		if parseInt != 0 {
			namesIds = append(namesIds, parseInt)
		}
	}
	return namesIds
}

func GetData() (q []models.QueueStruct) {
	sborKz, err := fetchDataSborkz()
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(sborKz) == 0 {
		return
	}

	lvl := func(t Sborkz) (s string) {
		if t.Vid == "black" {
			s = "drs"
		} else {
			s = "rs"
		}
		s += t.Lvlkz
		return s
	}

	rs := make(map[string][]Sborkz)
	for _, s := range sborKz {
		rs[s.Chatid] = append(rs[s.Chatid], s)
	}

	for cName, sbor := range rs {
		r := models.QueueStruct{
			CorpName: cName,
		}
		for _, s := range sbor {
			r.Level = lvl(s)
			r.Count = 1
			q = append(q, r)
		}
	}

	return q
}

//func fetchDataSborkz() ([]Sborkz, error) {
//	url := "https://123bot.ru/rssoyuzbot/Json/sborkz.php"
//	resp, err := http.Get(url)
//	if err != nil {
//		return nil, fmt.Errorf("ошибка запроса: %w", err)
//	}
//	defer resp.Body.Close()
//
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return nil, fmt.Errorf("ошибка чтения тела ответа: %w", err)
//	}
//
//	var sborkz []Sborkz
//
//	if err := json.Unmarshal(body, &sborkz); err != nil {
//		return nil, fmt.Errorf("ошибка парсинга JSON: %w", err)
//	}
//
//	return sborkz, nil
//}

func fetchDataSborkz() ([]Sborkz, error) {
	url := "https://123bot.ru/rssoyuzbot/Json/sborkz.php"

	// 1. Создание HTTP-клиента с таймаутом
	client := http.Client{
		Timeout: 15 * time.Second, // 🚨 Устанавливаем общий таймаут на запрос
	}

	// В отличие от http.Get, http.NewRequest позволяет настроить запрос
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	// 2. Выполнение запроса с настроенным клиентом
	resp, err := client.Do(req)
	if err != nil {
		// Ошибка может быть вызвана таймаутом или проблемой сети
		return nil, fmt.Errorf("ошибка выполнения запроса или таймаут: %w", err)
	}
	defer resp.Body.Close()

	// 3. Обработка статуса ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("сервер вернул код ошибки: %d", resp.StatusCode)
	}

	// 4. Чтение тела ответа с использованием io.ReadAll
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения тела ответа: %w", err)
	}
	fmt.Println(string(body))
	// 5. Парсинг JSON
	var sborkz []Sborkz
	if err := json.Unmarshal(body, &sborkz); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	return sborkz, nil
}
