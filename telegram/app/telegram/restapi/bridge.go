package restapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"telegram/models"
	"time"
)

func SendBridgeApp(m models.ToBridgeMessage) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	url := "http://bridge/bridge/inbox"

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
