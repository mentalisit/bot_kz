package restapi

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func SendCompendiumApp(m any) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	_, err = http.Post("http://compendium/compendium/inbox", "application/json", bytes.NewBuffer(data))
	if err != nil {
		_, err = http.Post("http://192.168.100.155:803/compendium/inbox", "application/json", bytes.NewBuffer(data))
		if err != nil {
			return err
		}
	}
	return nil
}
