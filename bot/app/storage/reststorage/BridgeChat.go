package reststorage

import (
	"encoding/json"
	"kz_bot/models"
	"net/http"
)

func (d Db) DBReadBridgeConfig() []models.BridgeConfig {
	var br []models.BridgeConfig
	resp, err := http.Get("http://storage/storage/bridge/read")
	if err != nil {
		resp, err = http.Get("http://192.168.100.155:804/storage/bridge/read")
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
	return br
}

func (d Db) UpdateBridgeChat(br models.BridgeConfig) {
	//TODO implement me
	panic("implement me")
}

func (d Db) InsertBridgeChat(br models.BridgeConfig) {
	//TODO implement me
	panic("implement me")
}
