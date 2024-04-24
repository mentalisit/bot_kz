package tg

import (
	"bytes"
	"compendium/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const apinametg = "kz_bot"

// отправка сообщения в телегу
func Send(chatId string, text string) error {
	m := models.SendText{
		Text:    text,
		Channel: chatId,
	}
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	resp, err := http.Post("http://"+apinametg+"/telegram/SendText", "application/json", bytes.NewBuffer(data))
	if err != nil {
		//_, err = http.Post("http://192.168.100.155:802/data", "application/json", bytes.NewBuffer(data))
		return err
	}
	if resp.StatusCode == http.StatusForbidden {
		return errors.New("forbidden")
	}
	return err
}

func SendPic(channelId string, text string, pic []byte) error {
	m := models.SendPic{
		Text:    text,
		Channel: channelId,
		Pic:     pic,
	}
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	resp, err := http.Post("http://"+apinametg+"/telegram/SendPic", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
		//resp, err = http.Post("http://192.168.100.155:802/data", "application/json", bytes.NewBuffer(data))
		return err
	}
	fmt.Println("resp.Status", resp.Status)
	return nil
}
