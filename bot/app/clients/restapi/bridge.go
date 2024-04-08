package restapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"kz_bot/models"
	"net/http"
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
