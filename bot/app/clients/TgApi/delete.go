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

func (t *TgApi) DelMessage(ChatId string, messageID int) {
	ch := utils.WaitForMessage("DelMessage")
	defer close(ch)

	m := apiRs{
		//FuncApi:   funcDeleteMessage,
		Channel:   ChatId,
		MessageId: strconv.Itoa(messageID),
	}

	data, err := json.Marshal(m)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", "http://telegram/delete_message", bytes.NewBuffer(data))
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
	}
}

func (t *TgApi) DelMessageSecond(chatId, messageId string, second int) {
	ch := utils.WaitForMessage("DelMessageSecond")
	defer close(ch)
	m := apiRs{
		//FuncApi:   funcDeleteMessageSecond,
		Channel:   chatId,
		MessageId: messageId,
		Second:    second,
	}

	data, err := json.Marshal(m)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", "http://telegram/delete_message_second", bytes.NewBuffer(data))
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
	}
	fmt.Printf("DelMessageSecond %+v\n", a)
}
