package TgApi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mentalisit/logger"
	"io/ioutil"
	"kz_bot/models"
	"net/http"
	"time"
)

const (
	funcSend                = "send"
	funcDeleteMessage       = "delete_message"
	funcEditMessage         = "edit_message"
	funcEditMessageTextKey  = "edit_message_text_key"
	funcSendDel             = "send_del"
	funcDeleteMessageSecond = "delete_message_second"
	funcCheckAdmin          = "check_admin"
	funcChatTyping          = "chat_typing"
	funcSendHelp            = "send_help"
	funcSendEmbed           = "send_embed"
	funcSendEmbedTime       = "send_embed_time"
	urlTg                   = "http://telegram" //192.168.100.14  telegram
)

type apiRs struct {
	FuncApi   string   `json:"funcApi"`
	Text      string   `json:"text"`
	Channel   string   `json:"channel"`
	MessageId string   `json:"messageId"`
	ParseMode string   `json:"parseMode"`
	Second    int      `json:"second"`
	LevelRs   string   `json:"levelRs"`
	Levels    []string `json:"levels"`
	UserName  string   `json:"userName"`
}

type answer struct {
	ArrString string    `json:"arrString"`
	ArrInt    int       `json:"arrInt"`
	ArrBool   bool      `json:"arrBool"`
	ArrError  error     `json:"arrError"`
	Time      time.Time `json:"time"`
}

type TgApi struct {
	log           *logger.Logger
	ChanRsMessage chan models.InMessage
}

func NewTgApi(log *logger.Logger) *TgApi {
	return &TgApi{log: log, ChanRsMessage: make(chan models.InMessage, 20)}
}

func convertAndSend(m any) (a answer, err error) {
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
		fmt.Printf("convertAndSend %+v\n", m)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("printBody ", string(body))
	}
	return a, err
}
