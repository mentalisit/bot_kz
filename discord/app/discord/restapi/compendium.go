package restapi

import (
	"bytes"
	"discord/models"
	"encoding/json"
	"net/http"
)

func SendCompendiumApp(m models.IncomingMessage) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	_, err = http.Post("http://compendiumnew/compendium/inbox", "application/json", bytes.NewBuffer(data))
	if err != nil {
		_, err = http.Post("http://192.168.100.131:880/compendium/inbox", "application/json", bytes.NewBuffer(data))
		if err != nil {
			return err
		}
	}
	return nil
}
