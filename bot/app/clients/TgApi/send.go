package TgApi

import (
	"fmt"
	"kz_bot/pkg/utils"
)

func (t *TgApi) SendChannelDelSecond(chatId, text string, second int) {
	ch := utils.WaitForMessage("SendChannelDelSecond")
	defer close(ch)
	m := apiRs{
		FuncApi: funcSendDel,
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
func (t *TgApi) SendChannel(chatId string, text string) int {
	ch := utils.WaitForMessage("SendChannel")
	defer close(ch)
	m := apiRs{
		FuncApi: funcSend,
		Text:    text,
		Channel: chatId,
	}

	a, err := convertAndSend(m)
	if err != nil {
		t.log.ErrorErr(err)
		return 0
	}
	fmt.Printf("SendChannel %+v\n", a)
	return a.ArrInt
}
func (t *TgApi) ChatTyping(chatId string) {
	ch := utils.WaitForMessage("ChatTyping")
	defer close(ch)
	m := apiRs{
		FuncApi: funcChatTyping,
		Channel: chatId,
	}
	_, err := convertAndSend(m)
	if err != nil {
		t.log.ErrorErr(err)
		return
	}
	//fmt.Printf("ChatTyping %+v\n", a)
}
func (t *TgApi) SendHelp(chatId, text string, mIdOld string) string {
	ch := utils.WaitForMessage("SendHelp")
	defer close(ch)
	m := apiRs{
		FuncApi:   funcSendHelp,
		Text:      text,
		Channel:   chatId,
		MessageId: mIdOld,
	}

	a, err := convertAndSend(m)
	if err != nil {
		t.log.ErrorErr(err)
		return ""
	}
	fmt.Printf("SendHelp %+v\n", a)
	return a.ArrString
}
func (t *TgApi) SendEmbed(lvlkz string, chatid string, text string) int {
	ch := utils.WaitForMessage("SendEmbed")
	defer close(ch)
	m := apiRs{
		FuncApi: funcSendEmbed,
		Text:    text,
		Channel: chatid,
		LevelRs: lvlkz,
	}

	a, err := convertAndSend(m)
	if err != nil {
		t.log.ErrorErr(err)
		return 0
	}
	fmt.Printf("SendEmbed %+v\n", a)
	return a.ArrInt
}
func (t *TgApi) SendEmbedTime(chatid string, text string) int {
	ch := utils.WaitForMessage("SendEmbedTime")
	defer close(ch)
	m := apiRs{
		FuncApi: funcSendEmbedTime,
		Text:    text,
		Channel: chatid,
	}
	a, err := convertAndSend(m)
	if err != nil {
		t.log.ErrorErr(err)
		return 0
	}
	fmt.Printf("SendEmbedTime %+v\n", a)
	return a.ArrInt
}
