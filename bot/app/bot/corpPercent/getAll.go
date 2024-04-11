package corpPercent

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func (b *Percent) getKeyAll() []string {
	var keys []string
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
			b.log.Error(fmt.Sprintln("Ошибка при загрузке XML:", err))
		}
		defer resp.Body.Close()

		// Читаем данные из тела ответа
		xmlData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			b.log.Error(fmt.Sprintln("Ошибка при чтении XML:", err))
		}

		// Создаем переменную для хранения данных о содержимом корзины Amazon S3
		var listBucketResult ListBucketResult

		// Распарсиваем XML-данные
		err = xml.Unmarshal(xmlData, &listBucketResult)
		if err != nil {
			b.log.Error(fmt.Sprintln("Ошибка при парсинге XML:", err))
		}

		if listBucketResult.NextMarker != "" {
			marker = listBucketResult.NextMarker
		}

		// Получаем текущее время минус 7 дней
		tenDaysAgo := time.Now().AddDate(0, 0, -7)

		for _, content := range listBucketResult.Contents {
			// Преобразуем строку даты в тип time.Time
			lastModifiedTime, errp := time.Parse(time.RFC3339, content.LastModified)
			if errp != nil {
				b.log.Error(fmt.Sprintln("Ошибка при преобразовании даты:", errp))
				continue
			}

			if !lastModifiedTime.Before(tenDaysAgo) {
				keys = append(keys, content.Key)
			}
		}
		if !listBucketResult.IsTruncated {
			break
		}
	}
	if len(keys) == 0 {
		b.log.Error("len(keys) == 0")
	}
	return keys
}
