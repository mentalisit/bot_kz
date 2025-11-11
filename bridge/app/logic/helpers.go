package logic

import (
	"bridge/models"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"unicode"
)

func (b *Bridge) ifTipDelSend(text string) {
	if b.in.Tip == "ds" {
		go b.discord.SendChannelDelSecondDs(b.in.ChatId, "```"+text+"```", 30)
		go b.discord.DeleteMessageDs(b.in.ChatId, b.in.MesId)
	} else if b.in.Tip == "tg" {
		go b.telegram.SendChannelDelSecond(b.in.ChatId, text, 30)
		go b.telegram.DeleteMessage(b.in.ChatId, b.in.MesId)
	} else if b.in.Tip == "wa" {
		go b.whatsapp.SendChannelDelSecond(b.in.ChatId, text, 30)
		go b.whatsapp.DeleteMessage(b.in.ChatId, b.in.MesId)
	}
}
func (b *Bridge) ifTipSendDel(text string) {
	if b.in.Tip == "ds" {
		go b.discord.SendChannelDelSecondDs(b.in.ChatId, text, 10)
	} else if b.in.Tip == "tg" {
		go b.telegram.SendChannelDelSecond(b.in.ChatId, text, 10)
	} else if b.in.Tip == "wa" {
		go b.whatsapp.SendChannelDelSecond(b.in.ChatId, text, 10)
	}
}

func (b *Bridge) ifChannelTip() {
	if b.in.Config == nil {
		b.in.Config = &models.Bridge2Config{}
	}
	b2c := models.Bridge2Configs{
		ChannelId:       b.in.ChatId,
		GuildId:         b.in.GuildId,
		CorpChannelName: b.in.Config.HostRelay,
		AliasName:       "",
		MappingRoles:    map[string]string{},
	}

	if b.in.Tip == "ds" {
		if b.in.Config.Channel["ds"] == nil {
			b.in.Config.Channel["ds"] = []models.Bridge2Configs{}
		}
		b.in.Config.Channel["ds"] = append(b.in.Config.Channel["ds"], b2c)
	}

	if b.in.Tip == "tg" {
		if b.in.Config.Channel["tg"] == nil {
			b.in.Config.Channel["tg"] = []models.Bridge2Configs{}
		}
		b.in.Config.Channel["tg"] = append(b.in.Config.Channel["tg"], b2c)
	}
	if b.in.Tip == "wa" {
		if b.in.Config.Channel["wa"] == nil {
			b.in.Config.Channel["wa"] = []models.Bridge2Configs{}
		}
		b.in.Config.Channel["wa"] = append(b.in.Config.Channel["wa"], b2c)
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
	if strings.HasPrefix(result, " ") {
		split := strings.Split(result, " ")
		return split[0]
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
		go b.telegram.DeleteMessage(b.in.ChatId, b.in.MesId)
	} else if b.in.Tip == "wa" {
		go b.whatsapp.DeleteMessage(b.in.ChatId, b.in.MesId)
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
	for _, d := range b.in.Config.Channel[b.in.Tip] {
		if d.ChannelId == b.in.ChatId {
			AliasName = d.AliasName
		}
	}
	return fmt.Sprintf("%s ([%s]%s)", b.in.Sender, strings.ToUpper(b.in.Tip), AliasName)
}

func ReplaceParticipantJIDForMap(data map[string]string) map[string]string {
	separator := "/"

	updatedData := make(map[string]string)
	for chatJID, valueString := range data {
		parts := strings.Split(valueString, separator)
		if len(parts) < 2 {
			updatedData[chatJID] = valueString // Если формат не соответствует ожидаемому, оставляем строку как есть.
			continue
		}
		currentJIDPart := parts[0] // Это JID, который нужно заменить (например, 85178361896964@lid)
		if strings.Contains(currentJIDPart, "85178361896964@lid") {
			newVal := "79991399754@s.whatsapp.net" + separator + parts[1]
			updatedData[chatJID] = newVal
		} else {
			updatedData[chatJID] = valueString
		}
	}

	return updatedData
}
