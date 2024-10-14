package restapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"telegram/models"
	"time"
)

func SendRsBotApp(m models.InMessage) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	url := "http://kz_bot/inbox"

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

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
