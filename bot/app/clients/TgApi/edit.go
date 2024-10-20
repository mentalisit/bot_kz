package TgApi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kz_bot/pkg/utils"
	"net/http"
	"strconv"
	"time"
)

func (t *TgApi) EditText(channel string, messageID int, text string) {
	fmt.Printf("   EditText channel %s, messageID %d, text %s\n", channel, messageID, text)
	ch := utils.WaitForMessage("EditText")
	defer close(ch)
	err := t.EditTextParse(channel, strconv.Itoa(messageID), text, "")
	if err != nil {
		t.log.ErrorErr(err)
	}
}
func (t *TgApi) EditTextParse(channel, messageID string, text, parse string) error {
	ch := utils.WaitForMessage("EditTextParse")
	defer close(ch)
	m := apiRs{
		//FuncApi:   funcEditMessage,
		Text:      text,
		Channel:   channel,
		MessageId: messageID,
		ParseMode: parse,
	}

	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", "http://telegram/edit_message", bytes.NewBuffer(data))
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

	var a answer
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		t.log.ErrorErr(err)
		t.log.Info(string(body))
	}
	if !a.ArrBool && a.ArrError != nil {
		return a.ArrError
	}
	fmt.Printf("EditTextParse %+v\n", a)
	return nil
}

func (t *TgApi) EditMessageTextKey(channel string, messageID int, text string, lvlkz string) {
	ch := utils.WaitForMessage("EditMessageTextKey")
	defer close(ch)
	m := apiRs{
		//FuncApi:   funcEditMessageTextKey,
		Text:      text,
		Channel:   channel,
		MessageId: strconv.Itoa(messageID),
		LevelRs:   lvlkz,
	}

	data, err := json.Marshal(m)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", "http://telegram/edit_message_text_key", bytes.NewBuffer(data))
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

	var a answer
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		t.log.ErrorErr(err)
		t.log.Info(string(body))
	}
	if !a.ArrBool && a.ArrError != nil {
		t.log.ErrorErr(a.ArrError)
		return
	}
	//fmt.Printf("EditMessageTextKey %+v\n", a)
}
