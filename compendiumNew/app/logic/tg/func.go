package tg

import (
	"bytes"
	"compendium/models"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const apinametg = "kz_bot"

// отправка сообщения в телегу
func Send(chatId string, text string) (string, error) {
	m := models.SendText{
		Text:    text,
		Channel: chatId,
	}
	data, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	resp, err := http.Post("http://"+apinametg+"/telegram/SendText", "application/json", bytes.NewBuffer(data))
	if err != nil {
		//_, err = http.Post("http://192.168.100.155:802/data", "application/json", bytes.NewBuffer(data))
		return "", err
	}
	if resp.StatusCode == http.StatusForbidden {
		return "", errors.New("forbidden")
	}
	var mid string
	err = json.NewDecoder(resp.Body).Decode(&mid)
	if err != nil {
		return "", err
	}
	return mid, err
}

func SendMarkdownV2(chatId string, text string) (string, error) {
	m := models.SendText{
		Text:    escapeMarkdownV2(text),
		Channel: chatId,
	}
	data, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	resp, err := http.Post("http://"+apinametg+"/telegram/SendText?parse=MarkdownV2", "application/json", bytes.NewBuffer(data))
	if err != nil {
		//_, err = http.Post("http://192.168.100.155:802/data", "application/json", bytes.NewBuffer(data))
		return "", err
	}
	if resp.StatusCode == http.StatusForbidden {
		return "", errors.New("forbidden")
	}
	var mid string
	err = json.NewDecoder(resp.Body).Decode(&mid)
	if err != nil {
		return "", err
	}
	return mid, err
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

func DeleteMessage(ChatId string, MesId string) error {
	s := models.DeleteMessageStruct{
		MessageId: MesId,
		Channel:   ChatId,
	}
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	_, err = http.Post("http://"+apinametg+"/telegram/del", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err

	}
	return nil
}

func EditMessage(chatId, mid string, text, ParseMode string) error {
	m := models.EditText{
		Text:      text,
		Channel:   chatId,
		MessageId: mid,
	}

	url := "http://" + apinametg + "/telegram/edit"
	if ParseMode != "" {
		url = url + "?parse=" + ParseMode
		m.Text = escapeMarkdownV2(text)
	}
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	resp, errr := http.Post(url, "application/json", bytes.NewBuffer(data))
	defer resp.Body.Close()
	if errr != nil {
		return errr
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(fmt.Sprintf("Code %d body:%s", resp.StatusCode, string(body)))
	}
	return nil
}
func escapeMarkdownV2(text string) string {
	// Специальные символы, которые нужно экранировать в MarkdownV2
	specialChars := "_*[]()~`>#+-=|{}.!"

	// Буфер для результата
	var builder strings.Builder

	// Переменные для отслеживания состояния
	var inLinkText bool
	var inLinkURL bool
	var linkTextBuffer strings.Builder
	var linkURLBuffer strings.Builder

	for i := 0; i < len(text); i++ {
		char := text[i]

		if char == '[' && !inLinkText && !inLinkURL {
			inLinkText = true
			builder.WriteByte(char)
			continue
		}

		if char == ']' && inLinkText && !inLinkURL {
			inLinkText = false
			builder.WriteString(linkTextBuffer.String())
			linkTextBuffer.Reset()
			builder.WriteByte(char)
			continue
		}

		if char == '(' && !inLinkText && !inLinkURL && i > 0 && text[i-1] == ']' {
			inLinkURL = true
			builder.WriteByte(char)
			continue
		}

		if char == ')' && !inLinkText && inLinkURL {
			inLinkURL = false
			builder.WriteString(linkURLBuffer.String())
			linkURLBuffer.Reset()
			builder.WriteByte(char)
			continue
		}

		if inLinkText {
			linkTextBuffer.WriteByte(char)
			continue
		}

		if inLinkURL {
			linkURLBuffer.WriteByte(char)
			continue
		}

		if strings.ContainsRune(specialChars, rune(char)) {
			builder.WriteByte('\\')
		}
		builder.WriteByte(char)
	}

	return builder.String()
}
