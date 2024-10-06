package TgApi

import (
	"fmt"
	"kz_bot/pkg/utils"
	"strconv"
)

func (t *TgApi) DelMessage(ChatId string, messageID int) {
	ch := utils.WaitForMessage("DelMessage")
	defer close(ch)
	m := apiRs{
		FuncApi:   funcDeleteMessage,
		Channel:   ChatId,
		MessageId: strconv.Itoa(messageID),
	}

	_, err := convertAndSend(m)
	if err != nil {
		t.log.ErrorErr(err)
		return
	}
	//fmt.Printf("DelMessage %+v\n", a)
}

func (t *TgApi) DelMessageSecond(chatId, messageId string, second int) {
	ch := utils.WaitForMessage("DelMessageSecond")
	defer close(ch)
	m := apiRs{
		FuncApi:   funcDeleteMessageSecond,
		Channel:   chatId,
		MessageId: messageId,
		Second:    second,
	}
	a, err := convertAndSend(m)
	if err != nil {
		t.log.ErrorErr(err)
		return
	}
	fmt.Printf("DelMessageSecond %+v\n", a)
}
