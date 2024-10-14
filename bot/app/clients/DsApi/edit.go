package DsApi

import (
	"fmt"
	"kz_bot/pkg/utils"
	"strconv"
)

func (d *DsApi) EditText(channel string, messageID string, text string) {
	ch := utils.WaitForMessage("EditText")
	defer close(ch)
	err := d.EditTextParse(channel, messageID, text, "")
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *DsApi) EditTextParse(channel, messageID string, text, parse string) error {
	ch := utils.WaitForMessage("EditTextParse")
	defer close(ch)
	m := apiRs{
		FuncApi:   funcEditMessage,
		Text:      text,
		Channel:   channel,
		MessageId: messageID,
		ParseMode: parse,
	}

	a, err := convertAndSend(m)
	if err != nil {
		d.log.ErrorErr(err)
		return err
	}
	fmt.Printf("EditTextParse %+v\n", a)
	return nil
}

func (d *DsApi) EditMessageTextKey(channel string, messageID int, text string, lvlkz string) {
	ch := utils.WaitForMessage("EditMessageTextKey")
	defer close(ch)
	m := apiRs{
		FuncApi:   funcEditMessageTextKey,
		Text:      text,
		Channel:   channel,
		MessageId: strconv.Itoa(messageID),
		LevelRs:   lvlkz,
	}

	a, err := convertAndSend(m)
	if err != nil {
		d.log.ErrorErr(err)
		return
	}
	fmt.Printf("EditMessageTextKey %+v\n", a)
}
