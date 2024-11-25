package ds

import (
	"context"
	"errors"
)

func (c *Client) DeleteMessage(chatId, messageId string) {
	req := &DeleteMessageRequest{
		Chatid: chatId,
		Mesid:  messageId,
	}
	_, err := c.client.DeleteMessage(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
	}
}
func (c *Client) DeleteMessageSecond(chatId, messageId string, second int) {
	req := &DeleteMessageSecondRequest{
		Chatid: chatId,
		Mesid:  messageId,
		Second: int32(second),
	}
	_, err := c.client.DeleteMessageSecond(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
	}
}
func (c *Client) SendChannelDelSecond(ChatId, text string, second int) {
	req := &SendChannelDelSecondRequest{
		Chatid: ChatId,
		Text:   text,
		Second: int32(second),
	}

	_, err := c.client.SendChannelDelSecond(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
	}
}

func (c *Client) SendPollChannel(m map[string]string, options []string) string {
	req := &SendPollRequest{
		Data:    m,
		Options: options,
	}

	pollMid, err := c.client.SendPoll(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
		return ""
	}
	return pollMid.Text
}
func (c *Client) SendHelp(chatid, title, description, oldMidHelps string, ifUser bool) string {
	req := &SendHelpRequest{
		Chatid:      chatid,
		Title:       title,
		Description: description,
		OldMidHelps: oldMidHelps,
		IfUser:      ifUser,
	}
	rR, err := c.client.SendHelp(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
	}
	return rR.Text
}
func (c *Client) CleanOldMessageChannel(chatId, lim string) {
	req := &CleanOldMessageChannelRequest{
		ChatId: chatId,
		Lim:    lim,
	}
	_, err := c.client.CleanOldMessageChannel(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
	}
}
func (c *Client) CleanRsBotOtherMessage() {
	req := &Empty{}
	_, err := c.client.CleanRsBotOtherMessage(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
	}
}

func (c *Client) SendComplexContent(chatid, text string) string {
	req := &SendComplexContentRequest{
		Chatid: chatid,
		Text:   text,
	}
	tr, err := c.client.SendComplexContent(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
		return ""
	}
	return tr.Text
}
func (c *Client) EditComplexButton(dsmesid, dschatid string, mapEmbed map[string]string) error {
	req := &EditComplexButtonRequest{
		Dsmesid:  dsmesid,
		Dschatid: dschatid,
		MapEmbed: mapEmbed,
	}
	er, err := c.client.EditComplexButton(context.Background(), req)
	if err != nil {
		return err
	}
	if er.ErrorMessage != "" {
		return errors.New(er.ErrorMessage)
	}
	return nil
}

func (c *Client) SendWebhook(mtext string, username string, channel string, avatar string) string {
	req := &SendWebhookRequest{
		Text:     mtext,
		Username: username,
		Chatid:   channel,
		Avatar:   avatar,
	}
	tR, err := c.client.SendWebhook(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
		return ""
	}
	return tR.Text
}
func (c *Client) SendDmText(text string, id string) {
	req := &SendDmTextRequest{
		Text:     text,
		AuthorID: id,
	}
	_, err := c.client.SendDmText(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
		return
	}
}

func (c *Client) ReplaceTextMessage(text string, guildid string) string {
	req := &ReplaceTextMessageRequest{
		Text:    text,
		Guildid: guildid,
	}
	tr, err := c.client.ReplaceTextMessage(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
		return ""
	}
	return tr.Text
}

func (c *Client) CheckAdmin(id string, channel string) bool {
	req := &CheckAdminRequest{
		Nameid: id,
		Chatid: channel,
	}
	fr, err := c.client.CheckAdmin(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
		return false
	}
	return fr.Flag
}

func (c *Client) Send(channel string, text string) string {
	req := &SendRequest{
		Chatid: channel,
		Text:   text,
	}
	tr, err := c.client.Send(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
		return ""
	}
	return tr.Text
}

func (c *Client) EditWebhook(text string, username string, channel string, dsmesid string, avatar string) {
	req := &EditWebhookRequest{
		Text:      text,
		Username:  username,
		ChatID:    channel,
		MID:       dsmesid,
		AvatarURL: avatar,
	}
	_, err := c.client.EditWebhook(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
	}
}

func (c *Client) ChannelTyping(channel string) {
	req := &ChannelTypingRequest{
		ChannelID: channel,
	}
	_, err := c.client.ChannelTyping(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
	}
}

func (c *Client) RoleToIdPing(s string, guildid string) (string, error) {
	req := &RoleToIdPingRequest{
		RolePing: s,
		Guildid:  guildid,
	}
	tr, err := c.client.RoleToIdPing(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
		return "", err
	}
	return tr.Text, nil
}

func (c *Client) QueueSend(queue string) {
	req := &QueueSendRequest{
		Text: queue,
	}
	_, _ = c.client.QueueSend(context.Background(), req)
}

func (c *Client) SendEmbedTime(channel string, text string) string {
	req := &SendEmbedTimeRequest{
		Chatid: channel,
		Text:   text,
	}
	tr, err := c.client.SendEmbedTime(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
		return ""
	}
	return tr.Text

}

func (c *Client) SendComplex(channel string, n map[string]string) string {
	req := &SendComplexRequest{
		Chatid:    channel,
		MapEmbeds: n,
	}
	tr, err := c.client.SendComplex(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
		return ""
	}
	return tr.Text
}

func (c *Client) SendEmbedText(chatid, title, text string) string {
	req := &SendEmbedTextRequest{
		Chatid: chatid,
		Title:  title,
		Text:   text,
	}
	tr, err := c.client.SendEmbedText(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
		return ""
	}
	return tr.Text
}

func (c *Client) Subscribe(id string, roles string, guildid string) int {
	req := &SubscrRequest{
		Nameid:   id,
		ArgRoles: roles,
		Guildid:  guildid,
	}
	ir, err := c.client.Subscribe(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
		return -1
	}
	return int(ir.Result)
}

func (c *Client) Unsubscribe(id string, roles string, guildid string) int {
	req := &SubscrRequest{
		Nameid:   id,
		ArgRoles: roles,
		Guildid:  guildid,
	}
	ir, err := c.client.Unsubscribe(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
		return -1
	}
	return int(ir.Result)
}
