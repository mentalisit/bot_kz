package logic

import (
	"bridge/models"
	"fmt"
	"strings"
	"sync"
)

func (b *Bridge) logicSendMessage() {
	var memory models.BridgeTempMemory
	memory.Message = make(map[string]string)
	memory.Message[b.in.ChatId] = b.in.MesId
	var flag bool
	var wg sync.WaitGroup // Общая WaitGroup для ВСЕХ сетевых операций

	chatIdsTG, chatIdsDS, chatIdsWA := b.GetChannels()

	toMessenger := models.BridgeSendToMessenger{
		Text:   b.in.Text,
		Sender: b.GetSenderName(),
		Avatar: b.in.Avatar,
		Extra:  b.in.Extra,
		Reply:  b.in.Reply,
	}

	if len(b.in.ReplyMap) > 0 {
		replyMap, err := b.storage.GetMapByLinkedID(b.in.ReplyMap)
		if len(replyMap) == 0 || err != nil {
			replyMap, err = b.storage.GetMapByLinkedID(ReplaceParticipantJIDForMap(b.in.ReplyMap))
		}
		if err == nil {
			toMessenger.ReplyMap = replyMap
			flag = true
		}
	}

	// DS
	lenChannels := len(chatIdsDS)
	resultChannelDs := make(chan models.MessageIds, lenChannels)

	if lenChannels > 0 {
		wg.Add(1)
		toDs := toMessenger
		if b.in.Reply != nil && b.in.Reply.Text != "" && b.in.Reply.UserName == "gote1st_bot" {
			at := strings.SplitN(b.in.Reply.Text, "\n", 2)
			toDs.Reply.UserName = at[0]
			if len(at) > 1 {
				toDs.Reply.Text = at[1]
			}
		}
		if toDs.Avatar == "" {
			toDs.Avatar = fmt.Sprintf("https://via.placeholder.com/128x128.png/%s/FFFFFF/?text=%s",
				GetRandomColor(), ExtractUppercase(b.in.Sender))
		}
		toDs.ChannelId = chatIdsDS
		if toDs.ReplyMap == nil {
			toDs.ReplyMap = make(map[string]string)
		}
		for _, configs := range b.in.Config.Channel["ds"] {
			toDs.ReplyMap["guild-"+configs.ChannelId] = configs.GuildId
		}
		go b.discord.SendBridgeArrayMessage(resultChannelDs, &wg, toDs)
	}

	//for TG and WA
	toMessenger.Text = fmt.Sprintf("%s\n%s", b.GetSenderName(), b.in.Text)
	if b.in.Reply != nil && b.in.Reply.Text != "" && !flag {
		toMessenger.Text = fmt.Sprintf("%s\n%s\nReply: %s", b.GetSenderName(), b.in.Text, b.in.Reply.Text)
	}

	//TG
	lenChTG := len(chatIdsTG)
	resultChannelTg := make(chan models.MessageIds, lenChTG)
	if lenChTG > 0 {
		wg.Add(1)
		toTg := toMessenger
		toTg.ChannelId = chatIdsTG

		go b.telegram.SendBridgeArrayMessage(resultChannelTg, &wg, toTg)
	}

	//WA
	lenChannelsWA := len(chatIdsWA)
	resultChannelWA := make(chan models.MessageIds, lenChannelsWA)
	if lenChannelsWA > 0 {
		wg.Add(1)
		toWa := toMessenger
		toWa.ChannelId = chatIdsWA
		go b.whatsapp.SendBridgeArrayMessage(resultChannelWA, &wg, toWa)
	}

	wg.Wait()
	close(resultChannelTg)
	close(resultChannelDs)
	close(resultChannelWA)
	for value := range resultChannelTg {
		memory.Message[value.ChatId] = value.MessageId
	}
	for value := range resultChannelDs {
		memory.Message[value.ChatId] = value.MessageId
	}
	for value := range resultChannelWA {
		memory.Message[value.ChatId] = value.MessageId
	}
	err := b.storage.SaveBridgeMap(memory.Message)
	if err != nil {
		b.log.ErrorErr(err)
	}
	b.in = models.ToBridgeMessage{}
}

/*
original := models.MessageIds{
			MessageId: b.in.MesId,
			ChatId:    b.in.ChatId,
		}
		if b.in.Tip == "ds" {
			memory.MessageDs = append(memory.MessageDs, original)
		} else if b.in.Tip == "tg" {
			memory.MessageTg = append(memory.MessageTg, original)
		} else if b.in.Tip == "wa" {
			memory.MessageWa = append(memory.MessageWa, original)
		}
*/
