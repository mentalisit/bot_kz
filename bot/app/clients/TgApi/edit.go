package TgApi

import (
	"fmt"
	"kz_bot/pkg/utils"
	"strconv"
)

func (t *TgApi) EditText(channel string, messageID int, text string) {
	fmt.Printf("   EditText channel %s, messageID %d, text %s\n", channel, messageID, text)
	ch := utils.WaitForMessage("EditText")
	defer close(ch)
	err := t.EditTextParse(channel, strconv.Itoa(messageID), text, "")
	if err != nil {
		t.log.ErrorErr(err)
	}
}
func (t *TgApi) EditTextParse(channel, messageID string, text, parse string) error {
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
		t.log.ErrorErr(err)
		return err
	}
	fmt.Printf("EditTextParse %+v\n", a)
	return nil
}

func (t *TgApi) EditMessageTextKey(channel string, messageID int, text string, lvlkz string) {
	ch := utils.WaitForMessage("EditMessageTextKey")
	defer close(ch)
	m := apiRs{
		FuncApi:   funcEditMessageTextKey,
		Text:      text,
		Channel:   channel,
		MessageId: strconv.Itoa(messageID),
		LevelRs:   lvlkz,
	}

	_, err := convertAndSend(m)
	if err != nil {
		t.log.ErrorErr(err)
		return
	}
	//fmt.Printf("EditMessageTextKey %+v\n", a)
}
