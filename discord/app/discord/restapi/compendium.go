package restapi

import (
	"bytes"
	"context"
	"discord/models"
	"encoding/json"
	"net/http"
	"time"
)

func SendCompendiumApp(m models.IncomingMessage) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	url := "http://compendiumnew/compendium/inbox"

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
