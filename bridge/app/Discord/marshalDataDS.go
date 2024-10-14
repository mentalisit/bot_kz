package ds

import (
	"bridge/models"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

const endpoint = "kz_bot"

func (d *Discord) MarshalDataSendBridgeAsync(message models.BridgeSendToMessenger) (dataReply []models.MessageIds, err error) {
	data, err := json.Marshal(message)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	url := "http://" + endpoint + "/discord/send/bridge"

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&dataReply)
	if err != nil {
		return
	}
	return dataReply, nil
}

func (d *Discord) MarshalDataDiscordDel(message models.DeleteMessageStruct) (err error) {
	data, err := json.Marshal(message)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	url := "http://" + endpoint + "/discord/del"

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	resp.Body.Close()
	return nil
}

func (d *Discord) MarshalDataDiscordSendDel(message models.SendTextDeleteSeconds) (err error) {
	data, err := json.Marshal(message)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	url := "http://" + endpoint + "/discord/SendDel"

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	resp.Body.Close()
	return nil
}
