package GetDataSborkz

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"queue/models"
	"strconv"
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

func fetchDataSborkz() ([]Sborkz, error) {
	url := "https://123bot.ru/rssoyuzbot/Json/sborkz.php"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения тела ответа: %w", err)
	}

	var sborkz []Sborkz

	//fmt.Println(string(body))

	if err := json.Unmarshal(body, &sborkz); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	return sborkz, nil
}
