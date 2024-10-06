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
		_, err = http.Post("http://192.168.100.131:802/inbox", "application/json", bytes.NewBuffer(data))
		if err != nil {
			return err
		}
	}
	return nil
}

//func GetRsConfig() ([]models.CorporationConfig, error) {
//	var br []models.CorporationConfig
//	resp, err := http.Get("http://storage/storage/rsbot/read")
//	if err != nil {
//		resp, err = http.Get("http://192.168.100.131:805/storage/rsbot/read")
//		if err != nil {
//			return nil, err
//		}
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusOK {
//		return nil, fmt.Errorf("error calling API: %d", resp.StatusCode)
//	}
//
//	err = json.NewDecoder(resp.Body).Decode(&br)
//	if err != nil {
//		return nil, err
//	}
//	return br, nil
//}
