package server

import (
	"bridge/models"
	"fmt"
	"strings"
	"sync"
)

func (b *Bridge) logicMessage() {
	if b.checkingForIdenticalMessage() {
		return
	}
	if b.in.Tip == "delDs" {
		b.RemoveMessage()
		return
	}
	if b.in.Tip == "dse" {
		//b.EditMessageDS()
		return
	}
	if b.in.Tip == "tge" {
		//b.EditMessageTG()// нужно исправить
		return
	}

	var memory models.BridgeTempMemory

	memory.Wg.Add(1)

	// Создаем WaitGroup для ожидания завершения всех горутин
	var wg sync.WaitGroup
	chatIdsTG, chatIdsDS := b.Channels()

	lenChTG := len(chatIdsTG)
	resultChannelTg := make(chan models.MessageIds, lenChTG)
	if lenChTG > 0 {
		wg.Add(1)
		textTg := fmt.Sprintf("%s\n%s", b.GetSenderName(), b.in.Text)
		if b.in.Reply != nil && b.in.Reply.Text != "" {
			textTg = fmt.Sprintf("%s\n%s\nReply: %s", b.GetSenderName(), b.in.Text, b.in.Reply.Text)
		}
		go b.telegram.SendBridgeArrayMessage(chatIdsTG, textTg, b.in.Extra, resultChannelTg, &wg)
	}

	// DS
	lenChannels := len(chatIdsDS)
	resultChannelDs := make(chan models.MessageIds, lenChannels)
	if b.in.Reply != nil && b.in.Reply.Text != "" && b.in.Reply.UserName == "gote1st_bot" {
		at := strings.SplitN(b.in.Reply.Text, "\n", 2)
		b.in.Reply.UserName = at[0]
		if len(at) > 1 {
			b.in.Reply.Text = at[1]
		}
	}
	if lenChannels > 0 {
		wg.Add(1)
		if b.in.Avatar == "" {
			b.in.Avatar = fmt.Sprintf("https://via.placeholder.com/128x128.png/%s/FFFFFF/?text=%s",
				GetRandomColor(), ExtractUppercase(b.in.Sender))
		}
		b.discord.SendBridgeFunc(b.in.Text, b.GetSenderName(), chatIdsDS, b.in.Extra, b.in.Avatar, b.in.Reply, resultChannelDs, &wg)
	}

	go func() {
		wg.Wait()
		close(resultChannelTg)
		close(resultChannelDs)
		for value := range resultChannelTg {
			memory.MessageTg = append(memory.MessageTg, value)
		}
		for value := range resultChannelDs {
			memory.MessageDs = append(memory.MessageDs, value)
		}
		memory.Wg.Done()
	}()
	memory.Wg.Wait()
	b.in = models.ToBridgeMessage{}
	b.messages = append(b.messages, memory)
}
