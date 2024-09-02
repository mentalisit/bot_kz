package restapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func RsbotQueue() (text string, err error) {
	done := make(chan struct{})

	go func() {
		resp, err := http.Get("http://storage/storage/rsbot/readqueue")
		if err != nil {
			resp, err = http.Get("http://192.168.100.131:804/storage/rsbot/readqueue")
			if err != nil {
				return
			}
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("error calling API: %d", resp.StatusCode)
		}

		err = json.NewDecoder(resp.Body).Decode(&text)
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
		err = fmt.Errorf("RsbotQueue завис: %+v\n, %+v\n", err, text)
	}
	return
}
