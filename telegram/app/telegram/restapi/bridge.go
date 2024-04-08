package restapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"telegram/models"
)

func SendBridgeApp(m models.ToBridgeMessage) error {
	fmt.Printf("Send to bridge: %s %s relay %s channel %s lenFile:%d\n", m.Sender, m.Text, m.Config.HostRelay, m.ChatId, len(m.Extra))
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	_, err = http.Post("http://bridge/bridge/inbox", "application/json", bytes.NewBuffer(data))
	if err != nil {
		_, err = http.Post("http://192.168.100.155:808/bridge/inbox", "application/json", bytes.NewBuffer(data))
		if err != nil {
			return err
		}
	}
	return nil
}
func GetBridgeConfig() ([]models.BridgeConfig, error) {
	var br []models.BridgeConfig
	resp, err := http.Get("http://storage/storage/bridge/read")
	if err != nil {
		resp, err = http.Get("http://192.168.100.155:804/storage/bridge/read")
		if err != nil {
			return nil, err
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error calling API: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&br)
	if err != nil {
		return nil, err
	}
	return br, nil
}
