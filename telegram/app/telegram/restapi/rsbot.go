package restapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"telegram/models"
)

func SendRsBotApp(m models.InMessage) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	_, err = http.Post("http://kz_bot/inbox", "application/json", bytes.NewBuffer(data))
	if err != nil {
		_, err = http.Post("http://192.168.100.155:803/inbox", "application/json", bytes.NewBuffer(data))
		if err != nil {
			return err
		}
	}
	return nil
}
