package otherQueue

import (
	"encoding/json"
	"fmt"
	"kz_bot/models"
	"net/http"
	"time"
)

func GetQueueLevel(level string) (t map[string][]models.QueueStruct, err error) {
	done := make(chan struct{})

	go func() {
		resp, err := http.Get("http://queue/queue?level=" + level)
		if err != nil {
			resp, err = http.Get("http://192.168.100.131:888/queue?level=" + level)
			if err != nil {
				return
			}
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("error calling API: %d", resp.StatusCode)
		}

		err = json.NewDecoder(resp.Body).Decode(&t)
		if err != nil {
			return
		}
		close(done)
	}()

	select {
	case <-done:
		// Запрос завершился до истечения таймаута
	case <-time.After(10 * time.Second):
		// Логируем, если запрос завис
		err = fmt.Errorf("GetQueueLevel завис: %+v\n, %+v\n", err, t)
	}
	return
}

func GetQueueAll() (t map[string][]models.QueueStruct, err error) {
	done := make(chan struct{})

	go func() {
		resp, err := http.Get("http://queue/queue")
		if err != nil {
			resp, err = http.Get("http://192.168.100.131:888/queue")
			if err != nil {
				return
			}
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("error calling API: %d", resp.StatusCode)
		}

		err = json.NewDecoder(resp.Body).Decode(&t)
		if err != nil {
			return
		}
		close(done)
		return

	}()

	select {
	case <-done:
		// Запрос завершился до истечения таймаута
	case <-time.After(10 * time.Second):
		// Логируем, если запрос завис
		err = fmt.Errorf("GetQueueAll завис: %+v\n, %+v\n", err, t)
	}
	return
}
