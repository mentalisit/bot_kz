package server

import (
	"bridge/models"
	"encoding/json"
	"fmt"
	"net/http"
)

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
