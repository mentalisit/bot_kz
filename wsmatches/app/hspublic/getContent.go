package hspublic

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"time"
	"ws/models"
)

func (h *HS) GetContentSevenDays() []models.Content {
	var contentNew []models.Content
	var marker string
	for {
		// Получаем XML-данные с помощью HTTP GET-запроса
		url := "https://storage.googleapis.com/hades-star-public-xq8f-d4rg-v0d9"
		if marker != "" {
			url = url + "/?&marker=" + marker
		}
		//fmt.Println("get " + url)
		resp, err := http.Get(url)
		if err != nil {
			h.log.ErrorErr(err)
		}
		defer resp.Body.Close()

		// Читаем данные из тела ответа
		xmlData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			h.log.ErrorErr(err)
		}

		// Создаем переменную для хранения данных о содержимом корзины Amazon S3
		var listBucketResult models.ListBucketResult

		// Распарсиваем XML-данные
		err = xml.Unmarshal(xmlData, &listBucketResult)
		if err != nil {
			h.log.ErrorErr(err)
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
				h.log.ErrorErr(err)
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
func (h *HS) GetContentAll() []models.Content {
	var contentNew []models.Content
	var marker string
	for {
		// Получаем XML-данные с помощью HTTP GET-запроса
		url := "https://storage.googleapis.com/hades-star-public-xq8f-d4rg-v0d9"
		if marker != "" {
			url = url + "/?&marker=" + marker
		}
		//fmt.Println("get " + url)
		resp, err := http.Get(url)
		if err != nil {
			h.log.Error(fmt.Sprintln("Ошибка при загрузке XML:", err))
		}
		defer resp.Body.Close()

		// Читаем данные из тела ответа
		xmlData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			h.log.Error(fmt.Sprintln("Ошибка при чтении XML:", err))
		}

		// Создаем переменную для хранения данных о содержимом корзины Amazon S3
		var listBucketResult models.ListBucketResult

		// Распарсиваем XML-данные
		err = xml.Unmarshal(xmlData, &listBucketResult)
		if err != nil {
			h.log.Error(fmt.Sprintln("Ошибка при парсинге XML:", err))
			h.log.InfoStruct("xmlData", string(xmlData))
		}

		if listBucketResult.NextMarker != "" {
			marker = listBucketResult.NextMarker
		}

		for _, content := range listBucketResult.Contents {
			contentNew = append(contentNew, content)
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
