package TgApi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kz_bot/pkg/utils"
	"net/http"
	"time"
)

func (t *TgApi) SendChannelDelSecond(chatId, text string, second int) {
	ch := utils.WaitForMessage("SendChannelDelSecond")
	defer close(ch)
	m := apiRs{
		//FuncApi: funcSendDel,
		Text:    text,
		Channel: chatId,
		Second:  second,
	}

	data, err := json.Marshal(m)
	if err != nil {
		t.log.ErrorErr(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", "http://telegram/send_del", bytes.NewBuffer(data))
	if err != nil {
		t.log.ErrorErr(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		t.log.ErrorErr(err)
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
	fmt.Printf("SendChannelDelSecond %+v\n", a)
}
func (t *TgApi) SendChannel(chatId string, text string) int {
	ch := utils.WaitForMessage("SendChannel")
	defer close(ch)
	m := apiRs{
		//FuncApi: funcSend,
		Text:    text,
		Channel: chatId,
	}
	data, err := json.Marshal(m)
	if err != nil {
		t.log.ErrorErr(err)
		return 0
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", "http://telegram/send", bytes.NewBuffer(data))
	if err != nil {
		t.log.ErrorErr(err)
		return 0
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		t.log.ErrorErr(err)
		return 0
	}
	defer resp.Body.Close()

	var a answer
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		t.log.ErrorErr(err)
		t.log.Info(string(body))
		return 0
	}
	fmt.Printf("SendChannel %+v\n", a)
	return a.ArrInt
}
func (t *TgApi) ChatTyping(chatId string) {
	ch := utils.WaitForMessage("ChatTyping")
	defer close(ch)
	m := apiRs{
		//FuncApi: funcChatTyping,
		Channel: chatId,
	}
	data, err := json.Marshal(m)
	if err != nil {
		t.log.ErrorErr(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", "http://telegram/chat_typing", bytes.NewBuffer(data))
	if err != nil {
		t.log.ErrorErr(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		t.log.ErrorErr(err)
		return
	}
	defer resp.Body.Close()

	var a answer
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(err)
		fmt.Println(string(body))
	}
	//fmt.Printf("ChatTyping %+v\n", a)
}
func (t *TgApi) SendHelp(chatId, text string, mIdOld string) string {
	ch := utils.WaitForMessage("SendHelp")
	defer close(ch)
	m := apiRs{
		//FuncApi:   funcSendHelp,
		Text:      text,
		Channel:   chatId,
		MessageId: mIdOld,
	}

	data, err := json.Marshal(m)
	if err != nil {
		t.log.ErrorErr(err)
		return ""
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", "http://telegram/send_help", bytes.NewBuffer(data))
	if err != nil {
		t.log.ErrorErr(err)
		return ""
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		t.log.ErrorErr(err)
		return ""
	}
	defer resp.Body.Close()

	var a answer
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		t.log.ErrorErr(err)
		t.log.Info(string(body))
		return ""
	}

	fmt.Printf("SendHelp %+v\n", a)
	return a.ArrString
}
func (t *TgApi) SendEmbed(lvlkz string, chatid string, text string) int {
	ch := utils.WaitForMessage("SendEmbed")
	defer close(ch)
	m := apiRs{
		//FuncApi: funcSendEmbed,
		Text:    text,
		Channel: chatid,
		LevelRs: lvlkz,
	}

	data, err := json.Marshal(m)
	if err != nil {
		t.log.ErrorErr(err)
		return 0
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", "http://telegram/send_embed", bytes.NewBuffer(data))
	if err != nil {
		t.log.ErrorErr(err)
		return 0
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		t.log.ErrorErr(err)
		return 0
	}
	defer resp.Body.Close()

	var a answer
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		t.log.ErrorErr(err)
		t.log.Info(string(body))
		return 0
	}

	fmt.Printf("SendEmbed %+v\n", a)
	return a.ArrInt
}
func (t *TgApi) SendEmbedTime(chatid string, text string) int {
	ch := utils.WaitForMessage("SendEmbedTime")
	defer close(ch)
	m := apiRs{
		//FuncApi: funcSendEmbedTime,
		Text:    text,
		Channel: chatid,
	}
	data, err := json.Marshal(m)
	if err != nil {
		t.log.ErrorErr(err)
		return 0
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", "http://telegram/send_embed_time", bytes.NewBuffer(data))
	if err != nil {
		t.log.ErrorErr(err)
		return 0
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		t.log.ErrorErr(err)
		return 0
	}
	defer resp.Body.Close()

	var a answer
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		t.log.ErrorErr(err)
		t.log.Info(string(body))
		return 0
	}

	fmt.Printf("SendEmbedTime %+v\n", a)
	return a.ArrInt
}
