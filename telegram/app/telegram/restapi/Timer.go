package restapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"telegram/models"
)

func SendInsertTimer(m models.Timer) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	_, err = http.Post("http://storage/storage/timer/insert", "application/json", bytes.NewBuffer(data))
	if err != nil {
		_, err = http.Post("http://192.168.100.155:803", "application/json", bytes.NewBuffer(data))
		if err != nil {
			return err
		}
	}
	return nil
}
