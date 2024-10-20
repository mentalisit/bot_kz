package tg

import (
	"bridge/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	send          = "send"
	deleteMessage = "delete_message"
	editMessage   = "edit_message"
	sendDel       = "send_del"
	urlTg         = "http://telegram" //192.168.100.14  telegram
)

//// отправка сообщения в телегу
//func (t *Telegram) Send(chatId string, text string) (string, error) {
//	m := apiRs{
//		FuncApi: send,
//		Text:    text,
//		Channel: chatId,
//	}
//	resp, err := convertAndSend(m)
//	if err != nil {
//		return "", err
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode == http.StatusForbidden {
//		return "", errors.New("forbidden")
//	}
//
//	var mid string
//	err = json.NewDecoder(resp.Body).Decode(&mid)
//	if err != nil {
//		return "", err
//	}
//	return mid, err
//}

func (t *Telegram) DeleteMessage(ChatId string, messageID string) error {
	m := apiRs{
		FuncApi:   deleteMessage,
		Channel:   ChatId,
		MessageId: messageID,
	}
	a, err := convertAndSend(m)
	if err != nil {
		t.log.ErrorErr(err)
		return err
	}
	fmt.Printf("DelMessageSecond %+v\n", a)
	return nil
}
func (t *Telegram) SendChannelDelSecond(chatId, text string, second int) {
	m := apiRs{
		FuncApi: sendDel,
		Text:    text,
		Channel: chatId,
		Second:  second,
	}

	a, err := convertAndSend(m)
	if err != nil {
		t.log.ErrorErr(err)
		return
	}
	fmt.Printf("SendChannelDelSecond %+v\n", a)
}

func (t *Telegram) SendBridgeArrayMessage(chatid []string, text string, extra []models.FileInfo, resultChannel chan<- models.MessageIds, wg *sync.WaitGroup) {

	defer wg.Done()

	// Создаем URL с параметрами запроса
	u, err := url.Parse(urlTg + "/bridge")
	if err != nil {
		return
	}

	m := models.BridgeSendToMessenger{
		Text:      text,
		ChannelId: chatid,
		Extra:     extra,
	}

	data, err := json.Marshal(m)
	if err != nil {
		t.log.ErrorErr(err)
		return
	}
	resp, err := http.Post(u.String(), "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.log.ErrorErr(err)
		return
	}
	var MessageIds []models.MessageIds
	err = json.NewDecoder(resp.Body).Decode(&MessageIds)
	if err != nil {
		t.log.Info(fmt.Sprintf("err resp.Body %+v\n", resp.Body))
		t.log.ErrorErr(err)
		return
	}

	for _, id := range MessageIds {
		resultChannel <- id
	}
}

func convertAndSend(m any) (a answer, err error) {
	fmt.Printf("convertAndSend %+v\n", m)
	data, err := json.Marshal(m)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", "http://telegram/func", bytes.NewBuffer(data))
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

	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		printBody(resp.Body)
	}
	return a, err
}
func printBody(r io.Reader) {
	body, _ := ioutil.ReadAll(r)
	fmt.Println("printBody ", string(body))
}

type apiRs struct {
	FuncApi   string `json:"funcApi"`
	Text      string `json:"text"`
	Channel   string `json:"channel"`
	MessageId string `json:"messageId"`
	//ParseMode string `json:"parseMode"`
	Second int `json:"second"`
	//LevelRs   string `json:"levelRs"`
	//Levels    []string `json:"levels"`
	//UserName  string `json:"userName"`
}
type answer struct {
	ArrString string    `json:"arrString"`
	ArrInt    int       `json:"arrInt"`
	ArrBool   bool      `json:"arrBool"`
	ArrError  error     `json:"arrError"`
	Time      time.Time `json:"time"`
}

func (t *Telegram) SendPollChannel(m map[string]string, options []string) string {
	s := models.Request{
		Data:    m,
		Options: options,
	}

	data, err := json.Marshal(s)
	if err != nil {
		t.log.ErrorErr(err)
		return ""
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", "http://telegram/send_poll", bytes.NewBuffer(data))
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

	var a string
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		t.log.ErrorErr(err)
		t.log.Info(string(body))
		return ""
	}

	fmt.Printf("SendPoll %+v\n", a)
	return a
}
