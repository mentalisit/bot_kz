package TgApi

import (
	"context"
	"errors"
	"strconv"
)

func (c *Client) DelMessage(ChatId string, messageID int) {
	er, err := c.client.DeleteMessage(context.Background(), &DeleteMessageRequest{
		Chatid: ChatId,
		Mesid:  strconv.Itoa(messageID),
	})
	if err != nil {
		//c.log.ErrorErr(err)
		return
	}
	if er.ErrorMessage != "" {
		c.log.Error(er.ErrorMessage)
		return
	}
}
func (c *Client) DelMessageSecond(chatId, messageId string, second int) {
	er, err := c.client.DeleteMessageSecond(context.Background(), &DeleteMessageSecondRequest{
		Chatid: chatId,
		Mesid:  messageId,
		Second: int32(second),
	})
	if err != nil {
		c.log.ErrorErr(err)
		return
	}
	if er.ErrorMessage != "" {
		c.log.Error(er.ErrorMessage)
	}
}
func (c *Client) EditTextParse(channel, messageID string, text, parse string) error {
	errorResponse, err := c.client.EditMessage(context.Background(), &EditMessageRequest{
		TextEdit:  text,
		ChatID:    channel,
		MID:       messageID,
		ParseMode: parse,
	})
	if err != nil {
		return err
	}
	if errorResponse.ErrorMessage != "" {
		return errors.New(errorResponse.ErrorMessage)
	}
	return nil
}
func (c *Client) EditMessageTextKey(channel string, messageID int, text string, lvlkz string) error {
	errorResponse, err := c.client.EditMessageTextKey(context.Background(), &EditMessageTextKeyRequest{
		ChatId:    channel,
		EditMesId: int32(messageID),
		TextEdit:  text,
		Lvlkz:     lvlkz,
	})
	if err != nil {
		if err.Error() != "rpc error: code = Unknown desc = Bad Request: message is not modified: specified new message content and reply markup are exactly the same as a current content and reply markup of the message" {
			return err
		}
	}
	if errorResponse != nil && errorResponse.ErrorMessage != "" {
		return errors.New(errorResponse.ErrorMessage)
	}
	return nil
}

func (c *Client) CheckAdminTg(ChatId string, userName string) (bool, error) {
	flagResponse, err := c.client.CheckAdmin(context.Background(), &CheckAdminRequest{
		Name:   userName,
		Chatid: ChatId,
	})
	if err != nil {
		return false, err
	}
	return flagResponse.Flag, nil
}
func (c *Client) SendChannelDelSecond(chatId, text string, second int) {
	_, err := c.client.SendChannelDelSecond(context.Background(), &SendMessageRequest{
		Text:   text,
		ChatID: chatId,
		Second: int32(second),
	})
	if err != nil {
		return
	}
}
func (c *Client) SendChannel(chatId string, text string) int {
	response, err := c.client.Send(context.Background(), &SendMessageRequest{
		Text:   text,
		ChatID: chatId,
	})
	if err != nil {
		return 0
	}
	mid, err := strconv.Atoi(response.Text)
	if err != nil {
		return 0
	}
	return mid

}
func (c *Client) ChatTyping(chatId string) {
	_, err := c.client.SendChannelTyping(context.Background(), &SendChannelTypingRequest{
		ChannelID: chatId,
	})
	if err != nil {
		return
	}
}
func (c *Client) SendHelp(chatId, text string, mIdOld string, ifUser bool) string {
	sendHelp, err := c.client.SendHelp(context.Background(), &SendHelpRequest{
		ChatId:      chatId,
		Text:        text,
		OldMidHelps: mIdOld,
		IfUser:      ifUser,
	})
	if err != nil {
		return ""
	}
	return sendHelp.Text
}
func (c *Client) SendEmbed(lvlkz string, chatid string, text string) int {
	response, err := c.client.SendEmbedText(context.Background(), &SendEmbedRequest{
		ChatId: chatid,
		Text:   text,
		Level:  lvlkz,
	})
	if err != nil {
		return 0
	}
	return int(response.Result)
}
func (c *Client) SendEmbedTime(chatid string, text string) int {
	response, err := c.client.SendEmbedTime(context.Background(), &SendMessageRequest{
		Text:   text,
		ChatID: chatid,
	})
	if err != nil {
		return 0
	}
	return int(response.Result)
}
