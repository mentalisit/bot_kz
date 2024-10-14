package DsApi

import (
	"fmt"
	"kz_bot/pkg/utils"
)

func (d *DsApi) SendChannelDelSecond(chatId, text string, second int) {
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
		d.log.ErrorErr(err)
		return
	}
	fmt.Printf("SendChannelDelSecond %+v\n", a)
}
func (d *DsApi) SendChannel(chatId string, text string) int {
	ch := utils.WaitForMessage("SendChannel")
	defer close(ch)
	m := apiRs{
		FuncApi: funcSend,
		Text:    text,
		Channel: chatId,
	}

	a, err := convertAndSend(m)
	if err != nil {
		d.log.ErrorErr(err)
		return 0
	}
	fmt.Printf("SendChannel %+v\n", a)
	return a.ArrInt
}
func (d *DsApi) ChatTyping(chatId string) {
	ch := utils.WaitForMessage("ChatTyping")
	defer close(ch)
	m := apiRs{
		FuncApi: funcChatTyping,
		Channel: chatId,
	}
	a, err := convertAndSend(m)
	if err != nil {
		d.log.ErrorErr(err)
		return
	}
	fmt.Printf("ChatTyping %+v\n", a)
}
func (d *DsApi) SendHelp(chatId, text, title string, levels []string) string {
	ch := utils.WaitForMessage("SendHelp")
	defer close(ch)
	m := apiRs{
		FuncApi:  funcSendHelp,
		Text:     text,
		UserName: title, //todo need sync
		Channel:  chatId,
		Levels:   levels,
	}

	a, err := convertAndSend(m)
	if err != nil {
		d.log.ErrorErr(err)
		return ""
	}
	fmt.Printf("SendHelp %+v\n", a)
	return a.ArrString
}
func (d *DsApi) SendEmbed(lvlkz string, chatid string, text string) int {
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
		d.log.ErrorErr(err)
		return 0
	}
	fmt.Printf("SendEmbed %+v\n", a)
	return a.ArrInt
}
func (d *DsApi) SendEmbedTime(chatid string, text string) int {
	ch := utils.WaitForMessage("SendEmbedTime")
	defer close(ch)
	m := apiRs{
		FuncApi: funcSendEmbedTime,
		Text:    text,
		Channel: chatid,
	}
	a, err := convertAndSend(m)
	if err != nil {
		d.log.ErrorErr(err)
		return 0
	}
	fmt.Printf("SendEmbedTime %+v\n", a)
	return a.ArrInt
}

func (d *DsApi) SendWebhook(mtext string, username string, channel string, guildid string, avatar string) {
	//todo
	panic("implement me")
}
