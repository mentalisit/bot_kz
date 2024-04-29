package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/mentalisit/logger"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"time"
	"ws/config"
	"ws/dbpostgres"
	"ws/dbredis"
	"ws/models"
)

var log *logger.Logger

func main() {
	cfg := config.InitConfig()
	log = logger.LoggerZap(cfg.Logger.Token, cfg.Logger.ChatId, cfg.Logger.Webhook)

	for {
		now := time.Now()

		if now.Second() == 0 && now.Minute() == 0 {
			newContent := getKeyAllSevenDays()
			downloadFromCloud(newContent)
		}

		//if now.Hour() == 5 && now.Minute() == 30 && now.Second() == 0 {
		//	newContent := getKeyAll()
		//	downloadFromCloud(newContent)
		//}
		time.Sleep(1 * time.Second)
	}
}

func downloadFromCloud(newContent []models.Content) {
	r := dbredis.NewDb(log)
	p := dbpostgres.NewDb(log)
	var count int

	listHCorp := make(map[string]models.LevelCorp)

	all, err := p.ReadCorpLevelAll()
	if err != nil {
		log.ErrorErr(err)
		return
	}

	for _, corp := range all {
		listHCorp[corp.HCorp] = corp
	}

	var corpsdata []models.CorpsData

	for _, cont := range newContent {
		count++
		corpData := r.ReadCorpData(cont.Key)
		if corpData == nil {
			corpData = getKey(cont.Key)
			r.SaveCorpDate(cont.Key, *corpData)
		}
		corpsdata = append(corpsdata, *corpData)

		if listHCorp[corpData.Corp1Name].HCorp != "" {
			c := listHCorp[corpData.Corp1Name]
			if c.EndDate.Before(corpData.DateEnded) {
				c.EndDate = corpData.DateEnded
				c.Percent = c.Level - 1

				fmt.Println(c)
				p.InsertUpdateCorpLevel(c)
			}

		}
	}
	fmt.Println("count", count)

	marshal, err := json.Marshal(corpsdata)
	if err != nil {
		log.ErrorErr(err)
		return
	}
	err = os.WriteFile("ws/ws.json", marshal, 0644)
	if err != nil {
		log.ErrorErr(err)
		return
	}
}

func getKey(key string) *models.CorpsData {
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
	var data models.CorporationsData

	// Распаковываем JSON в структуру данных
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Error(fmt.Sprintln("Ошибка при распаковке JSON:", err))
		return nil
	}
	data.SortWin()

	corpData := &models.CorpsData{
		Corp1Name:  data.Corporation1Name,
		Corp2Name:  data.Corporation2Name,
		Corp1Score: data.Corporation1Score,
		Corp2Score: data.Corporation2Score,
		DateEnded:  data.DateEnded,
	}

	return corpData
}

func getKeyAll() []models.Content {
	var contentNew []models.Content
	var marker string
	for {
		// Получаем XML-данные с помощью HTTP GET-запроса
		url := "https://storage.googleapis.com/hades-star-public-xq8f-d4rg-v0d9"
		if marker != "" {
			url = url + "/?&marker=" + marker
		}
		fmt.Println("get " + url)
		resp, err := http.Get(url)
		if err != nil {
			log.Error(fmt.Sprintln("Ошибка при загрузке XML:", err))
		}
		defer resp.Body.Close()

		// Читаем данные из тела ответа
		xmlData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error(fmt.Sprintln("Ошибка при чтении XML:", err))
		}

		// Создаем переменную для хранения данных о содержимом корзины Amazon S3
		var listBucketResult models.ListBucketResult

		// Распарсиваем XML-данные
		err = xml.Unmarshal(xmlData, &listBucketResult)
		if err != nil {
			log.Error(fmt.Sprintln("Ошибка при парсинге XML:", err))
		}

		if listBucketResult.NextMarker != "" {
			marker = listBucketResult.NextMarker
		}

		// Получаем текущее время минус 7 дней
		//tenDaysAgo := time.Now().AddDate(0, 0, -7)

		for _, content := range listBucketResult.Contents {
			// Преобразуем строку даты в тип time.Time
			//lastModifiedTime, errp := time.Parse(time.RFC3339, content.LastModified)
			//if errp != nil {
			//	log.Error(fmt.Sprintln("Ошибка при преобразовании даты:", errp))
			//	continue
			//}

			//if !lastModifiedTime.Before(tenDaysAgo) {
			contentNew = append(contentNew, content)
			//}
		}
		if !listBucketResult.IsTruncated {
			break
		}
	}
	sort.Slice(contentNew, func(i, j int) bool {
		return contentNew[i].LastModified > contentNew[j].LastModified
	})

	return contentNew
}
func getKeyAllSevenDays() []models.Content {
	var contentNew []models.Content
	var marker string
	for {
		// Получаем XML-данные с помощью HTTP GET-запроса
		url := "https://storage.googleapis.com/hades-star-public-xq8f-d4rg-v0d9"
		if marker != "" {
			url = url + "/?&marker=" + marker
		}
		fmt.Println("get " + url)
		resp, err := http.Get(url)
		if err != nil {
			log.Error(fmt.Sprintln("Ошибка при загрузке XML:", err))
		}
		defer resp.Body.Close()

		// Читаем данные из тела ответа
		xmlData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error(fmt.Sprintln("Ошибка при чтении XML:", err))
		}

		// Создаем переменную для хранения данных о содержимом корзины Amazon S3
		var listBucketResult models.ListBucketResult

		// Распарсиваем XML-данные
		err = xml.Unmarshal(xmlData, &listBucketResult)
		if err != nil {
			log.Error(fmt.Sprintln("Ошибка при парсинге XML:", err))
		}

		if listBucketResult.NextMarker != "" {
			marker = listBucketResult.NextMarker
		}

		//Получаем текущее время минус 7 дней
		tenDaysAgo := time.Now().AddDate(0, 0, -7)

		for _, content := range listBucketResult.Contents {
			// Преобразуем строку даты в тип time.Time
			lastModifiedTime, errp := time.Parse(time.RFC3339, content.LastModified)
			if errp != nil {
				log.Error(fmt.Sprintln("Ошибка при преобразовании даты:", errp))
				continue
			}

			if !lastModifiedTime.Before(tenDaysAgo) {
				contentNew = append(contentNew, content)
			}
		}
		if !listBucketResult.IsTruncated {
			break
		}
	}
	sort.Slice(contentNew, func(i, j int) bool {
		return contentNew[i].LastModified > contentNew[j].LastModified
	})

	return contentNew
}
