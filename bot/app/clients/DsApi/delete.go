package DsApi

import (
	"fmt"
	"kz_bot/pkg/utils"
)

func (d *DsApi) DelMessage(ChatId string, messageID string) {
	ch := utils.WaitForMessage("DelMessage")
	defer close(ch)
	fmt.Printf("DelMessage %d\n", messageID)
	m := apiRs{
		FuncApi:   funcDeleteMessage,
		Channel:   ChatId,
		MessageId: messageID,
	}

	a, err := convertAndSend(m)
	if err != nil {
		d.log.ErrorErr(err)
		return
	}
	fmt.Printf("DelMessage %+v\n", a)
}

func (d *DsApi) DelMessageSecond(chatId, messageId string, second int) {
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
		d.log.ErrorErr(err)
		return
	}
	fmt.Printf("DelMessageSecond %+v\n", a)
}
