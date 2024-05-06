package restapi

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func SendCompendiumAppOld(m any) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	_, err = http.Post("http://compendium2/compendium/inbox", "application/json", bytes.NewBuffer(data))
	if err != nil {
		_, err = http.Post("http://192.168.100.131:809/compendium/inbox", "application/json", bytes.NewBuffer(data))
		if err != nil {
			return err
		}
	}
	return nil
}
func SendCompendiumApp(m any) error {
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
