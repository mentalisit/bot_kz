package DiscordClient

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"kz_bot/clients/helper"
	"kz_bot/clients/restapi"
	"kz_bot/models"
	"path/filepath"
)

func (d *Discord) filterNewBridge(m *discordgo.MessageCreate, mes models.ToBridgeMessage) {
	mes.Text = d.replaceTextMessage(m.Content, m.GuildID)
	mes.Sender = d.getAuthorName(m)
	mes.Tip = "ds"
	mes.MesId = m.ID
	mes.GuildId = m.GuildID
	mes.TimestampUnix = m.Timestamp.Unix()
	mes.Avatar = m.Author.AvatarURL("")

	d.handleDownloadBridge(&mes, m)

	if m.ReferencedMessage != nil {
		usernameR := m.ReferencedMessage.Author.String() //.Username
		if m.ReferencedMessage.Member != nil && m.ReferencedMessage.Member.Nick != "" {
			usernameR = m.ReferencedMessage.Member.Nick
		}
		mes.Reply = &models.BridgeMessageReply{
			TimeMessage: m.ReferencedMessage.Timestamp.Unix(),
			Text:        d.replaceTextMessage(m.ReferencedMessage.Content, m.GuildID),
			Avatar:      m.ReferencedMessage.Author.AvatarURL(""),
			UserName:    usernameR,
		}
	}
	if mes.Text != "" || len(mes.Extra) > 0 {
		err := restapi.SendBridgeApp(mes)
		if err != nil {
			d.log.ErrorErr(err)
		}
	}
}

func (d *Discord) handleDownloadBridge(mes *models.ToBridgeMessage, m *discordgo.MessageCreate) {
	if len(m.StickerItems) > 0 {
		mes.Text = fmt.Sprintf("https://cdn.discordapp.com/stickers/%s.png", m.Message.StickerItems[0].ID)
	}
	if len(m.Attachments) > 0 {
		for _, a := range m.Attachments {
			f := models.FileInfo{
				Name: a.Filename,
				Data: nil,
				URL:  a.URL,
				Size: int64(a.Size),
			}
			if filepath.Ext(a.Filename) == ".apk" {
				f.URL = ""
				data, err := helper.DownloadFile(a.URL)
				if err != nil {
					d.log.ErrorErr(err)
				}
				f.Data = data
			}

			mes.Extra = append(mes.Extra, f)
		}
	}
}
