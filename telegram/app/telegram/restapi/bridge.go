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
		_, err = http.Post("http://192.168.100.131:808/bridge/inbox", "application/json", bytes.NewBuffer(data))
		if err != nil {
			return err
		}
	}
	return nil
}

func GetBridgeConfig() map[string]models.BridgeConfig {
	var br []models.BridgeConfig
	resp, err := http.Get("http://bridge/bridge/config")
	if err != nil {
		resp, err = http.Get("http://192.168.100.131:808/bridge/config")
		if err != nil {
			return nil
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	err = json.NewDecoder(resp.Body).Decode(&br)
	if err != nil {
		return nil
	}

	var bridgeCounter = 0
	var bridge string

	bridgeConfig := make(map[string]models.BridgeConfig)

	for _, configBridge := range br {
		bridgeConfig[configBridge.NameRelay] = configBridge
		bridgeCounter++
		bridge = bridge + fmt.Sprintf("%s, ", configBridge.HostRelay)
	}
	fmt.Printf("Загружено конфиг мостов %d : %s\n", bridgeCounter, bridge)

	return bridgeConfig
}
