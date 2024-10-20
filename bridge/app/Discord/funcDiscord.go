package ds

import (
	"bridge/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mentalisit/logger"
	"io/ioutil"
	"net/http"
	"time"
)

type Discord struct {
	log *logger.Logger
}

func NewDiscord(log *logger.Logger) *Discord {
	return &Discord{log: log}
}

func (d *Discord) DeleteMessageDs(ChatId, MesId string) {
	s := models.DeleteMessageStruct{
		MessageId: MesId,
		Channel:   ChatId,
	}
	err := d.MarshalDataDiscordDel(s)
	if err != nil {
		d.log.ErrorErr(err)
		d.log.InfoStruct("DeleteMessageDs s: ", s)
	}
}
func (d *Discord) SendChannelDelSecondDs(chatId, text string, second int) {
	s := models.SendTextDeleteSeconds{
		Text:    text,
		Channel: chatId,
		Seconds: second,
	}
	err := d.MarshalDataDiscordSendDel(s)
	if err != nil {
		d.log.ErrorErr(err)
		d.log.InfoStruct("SendChannelDelSecondDs s:", s)
	}
}

func (d *Discord) SendPollChannel(m map[string]string, options []string) string {
	s := models.Request{
		Data:    m,
		Options: options,
	}

	data, err := json.Marshal(s)
	if err != nil {
		d.log.ErrorErr(err)
		return ""
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//todo need move to nes service
	req, err := http.NewRequestWithContext(ctx, "POST", "http://kz_bot/discord/send_poll", bytes.NewBuffer(data))
	if err != nil {
		d.log.ErrorErr(err)
		return ""
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		d.log.ErrorErr(err)
		return ""
	}
	defer resp.Body.Close()

	var a string
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		d.log.ErrorErr(err)
		d.log.Info(string(body))
		return ""
	}

	fmt.Printf("SendPoll %+v\n", a)
	return a
}
