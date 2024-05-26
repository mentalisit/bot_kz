package ds

import (
	"bytes"
	"compendium/models"
	"encoding/json"
	"fmt"
	"net/http"
)

const apinameds = "kz_bot"

func SendChannel(chatId, text string) (string, error) {
	m := models.SendText{
		Text:    text,
		Channel: chatId,
	}
	data, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	resp, errs := http.Post("http://"+apinameds+"/discord/SendText", "application/json", bytes.NewBuffer(data))
	if errs != nil {
		//resp, err = http.Post("http://192.168.100.155:802/data", "application/json", bytes.NewBuffer(data))
		return "", errs
	}
	var mid string
	err = json.NewDecoder(resp.Body).Decode(&mid)
	if err != nil {
		return "", err
	}

	return mid, nil
}

func EditMessage(chatId, mid, text string) error {
	m := models.EditText{
		Text:      text,
		Channel:   chatId,
		MessageId: mid,
	}
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	_, errs := http.Post("http://"+apinameds+"/discord/edit", "application/json", bytes.NewBuffer(data))
	if errs != nil {
		return errs
	}

	return nil
}

func SendChannelPic(chatId, text string, pic []byte) error {
	m := models.SendPic{
		Text:    text,
		Channel: chatId,
		Pic:     pic,
	}
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	resp, err := http.Post("http://"+apinameds+"/discord/SendPic", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
		//resp, err = http.Post("http://192.168.100.155:802/data", "application/json", bytes.NewBuffer(data))
		return err
	}
	fmt.Println("resp.Status", resp.Status)
	return nil
}

func DeleteMessage(ChatId, MesId string) error {
	s := models.DeleteMessageStruct{
		MessageId: MesId,
		Channel:   ChatId,
	}
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	_, err = http.Post("http://"+apinameds+"/discord/del", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	return nil
}
