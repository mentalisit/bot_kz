package DiscordClient

import (
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
)

func (d *Discord) SendHelp(chatid, title, description, oldMidHelps string, ifUser bool) string {
	if oldMidHelps != "" {
		if !ifUser {
			messages, _ := d.S.ChannelMessages(chatid, 10, "", oldMidHelps, "")
			if len(messages) < 3 {
				return oldMidHelps
			}
		}
		go d.DeleteMessage(chatid, oldMidHelps)
	}
	Emb := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       16711680,
		Description: description,
		Title:       title,
	}

	m, err := d.S.ChannelMessageSendComplex(chatid, &discordgo.MessageSend{
		Components: d.AddButtonsStartQueue(chatid),
		Embed:      Emb,
	})
	if err != nil {
		d.log.ErrorErr(err)
		return ""
	}

	return m.ID
}
func (d *Discord) AddButtonsStartQueue(chatid string) []discordgo.MessageComponent {
	var mc []discordgo.MessageComponent
	var components []discordgo.MessageComponent
	_, config := d.CheckChannelConfigDS(chatid)
	levels := d.storage.Db.ReadTop5Level(config.CorpName)
	getButton := func(level string) string {
		after, found := strings.CutPrefix(level, "rs")
		if found {
			return after + "+"
		}
		after, found = strings.CutPrefix(level, "drs")
		if found {
			return after + "+"
		}
		after, found = strings.CutPrefix(level, "solo")
		if found {
			return after + "+"
		}
		return level
	}

	if len(levels) > 2 {
		for _, level := range levels {
			button := discordgo.Button{
				Style:    discordgo.DangerButton,
				Label:    getButton(level),
				CustomID: getButton(level),
			}
			components = append(components, button)
		}
	}

	if len(components) == 0 {
		for i := 7; i < 12; i++ {
			l := strconv.Itoa(i)

			button := discordgo.Button{
				Label:    l + "+",
				Style:    discordgo.SecondaryButton,
				CustomID: l + "+",
			}
			components = append(components, button)

		}
	}
	mc = append(mc, discordgo.ActionsRow{Components: RemoveDuplicates(components)})
	return mc
}
func RemoveDuplicates[T comparable](a []T) []T {
	result := make([]T, 0, len(a))
	temp := map[T]struct{}{}
	for _, item := range a {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
