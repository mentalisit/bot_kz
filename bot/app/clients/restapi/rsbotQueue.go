package restapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func RsbotQueue() (string, error) {
	var text string
	resp, err := http.Get("http://storage/storage/rsbot/readqueue")
	if err != nil {
		resp, err = http.Get("http://192.168.100.155:804/storage/rsbot/readqueue")
		if err != nil {
			return "", err
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error calling API: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&text)
	if err != nil {
		return "", err
	}
	return text, nil
}
