package server

import (
	"bridge/models"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func (b *Bridge) ifTipDelSend(text string) {
	if b.in.Tip == "ds" {
		go b.discord.SendChannelDelSecondDs(b.in.ChatId, "```"+text+"```", 30)
		go b.discord.DeleteMessageDs(b.in.ChatId, b.in.MesId)
	} else if b.in.Tip == "tg" {
		go b.telegram.SendChannelDelSecond(b.in.ChatId, text, 30)
		mid, err := strconv.Atoi(b.in.MesId)
		if err != nil {
			return
		}
		go b.telegram.DeleteMessage(b.in.ChatId, mid)
	}
}

func (b *Bridge) ifChannelTip() {
	if b.in.Tip == "ds" {
		b.in.Config.ChannelDs = append(b.in.Config.ChannelDs, models.BridgeConfigDs{
			ChannelId:       b.in.ChatId,
			GuildId:         b.in.GuildId,
			CorpChannelName: b.in.Config.HostRelay,
			AliasName:       "",
			MappingRoles:    map[string]string{},
		})
	}
	if b.in.Tip == "tg" {
		b.in.Config.ChannelTg = append(b.in.Config.ChannelTg, models.BridgeConfigTg{
			ChannelId:       b.in.ChatId,
			CorpChannelName: b.in.Config.HostRelay,
			AliasName:       "",
			MappingRoles:    map[string]string{},
		})
	}
}
func GetRandomColor() string {
	// Генерируем случайные значения для красного, зеленого и синего цветов
	red := rand.Intn(256)
	green := rand.Intn(256)
	blue := rand.Intn(256)

	// Форматируем цвет в HEX
	colorHex := fmt.Sprintf("%02X%02X%02X", red, green, blue)

	return colorHex
}
func ExtractUppercase(input string) string {
	var result string

	for _, char := range input {
		if unicode.IsUpper(char) {
			result += string(char)
		}
	}
	if result == "" {
		return input
	}

	return result
}

var messageText string
var messageAuthor string

// проверка на повторное сообщение
func (b *Bridge) checkingForIdenticalMessage() bool {
	if messageAuthor == b.in.Sender {
		if len(b.in.Extra) > 0 {
			if b.in.Extra[0].Name == messageText {
				b.delIncomingMessage()
				return true
			} else {
				messageText = b.in.Extra[0].Name
			}
		} else if messageText == b.in.Text {
			b.delIncomingMessage()
			return true
		}
	}
	messageText = b.in.Text
	messageAuthor = b.in.Sender
	return false
}

// удаление входящего сообщения
func (b *Bridge) delIncomingMessage() {
	if b.in.Tip == "ds" {
		go b.discord.DeleteMessageDs(b.in.ChatId, b.in.MesId)
	} else if b.in.Tip == "tg" {
		mid, err := strconv.Atoi(b.in.MesId)
		if err != nil {
			return
		}
		go b.telegram.DeleteMessage(b.in.ChatId, mid)
	}
}

// TODO
func (b *Bridge) replaceTextMentionRsRole(input, guildId string) string {
	//re := regexp.MustCompile(`@&rs([4-9]|1[0-2])`)
	//output := re.ReplaceAllStringFunc(input, func(s string) string {
	//	return b.client.Ds.TextToRoleRsPing(s[2:], guildId)
	//})
	return input
}

func replaceTextMap(text string, m map[string]string) string {
	mentionPattern := `@(\w+)|<@(\w+)>`
	mentionRegex := regexp.MustCompile(mentionPattern)
	text = mentionRegex.ReplaceAllStringFunc(text, func(match string) string {
		if value, ok := m[match]; ok {
			// Если значение найдено, заменяем упоминание на значение из map
			return value
		}

		// Если значение не найдено, оставляем упоминание без изменений
		return match
	})
	//fmt.Println(modifiedText)
	return text
}

// GetSenderName конконтенация имени
func (b *Bridge) GetSenderName() string {
	AliasName := ""
	if b.in.Tip == "ds" {
		for _, d := range b.in.Config.ChannelDs {
			if d.ChannelId == b.in.ChatId {
				AliasName = d.AliasName
			}
		}
	} else if b.in.Tip == "tg" {
		for _, d := range b.in.Config.ChannelTg {
			if d.ChannelId == b.in.ChatId {
				AliasName = d.AliasName
			}
		}
	}
	return fmt.Sprintf("%s ([%s]%s)", b.in.Sender, strings.ToUpper(b.in.Tip), AliasName)
}
