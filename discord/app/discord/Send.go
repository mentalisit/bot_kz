package DiscordClient

import (
	"bytes"
	"discord/discord/helpers"
	"discord/models"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (d *Discord) SendEmbedText(chatid, title, text string) *discordgo.Message {
	if chatid == "1198012575615561979" || chatid == "1253851338408857641" {
		d.SendOrEditTopEmbedText(chatid, title, text)
		return &discordgo.Message{}
	}
	Emb := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       16711680,
		Description: text,
		Title:       title,
	}
	m, err := d.S.ChannelMessageSendEmbed(chatid, Emb)
	if err != nil {
		d.log.Error("chatid " + chatid + " " + err.Error())
		return &discordgo.Message{}
	}
	return m
}

func (d *Discord) SendOrEditTopEmbedText(chatid, title, text string) {
	messages, err := d.S.ChannelMessages(chatid, 10, "", "", "")
	if err != nil {
		d.log.ErrorErr(err)
		return
	}
	if len(messages) == 0 {
		return
	}
	var mId string
	for _, message := range messages {
		if len(message.Embeds) > 0 {
			if message.Embeds[0].Title == title {
				mId = message.ID
			}
		}
	}

	Emb := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       16711680,
		Description: text,
		Title:       title,
	}

	if mId == "" {
		_, _ = d.S.ChannelMessageSendEmbed(chatid, Emb)
		return
	} else {
		_, _ = d.S.ChannelMessageEditEmbed(chatid, mId, Emb)
	}

}

func (d *Discord) SendOrEditEmbedImage(channelID, title, imageUrl string) error {
	messages, err := d.S.ChannelMessages(channelID, 10, "", "", "")
	if err != nil {
		return fmt.Errorf("ошибка получения сообщений: %v", err)
	}

	var mId string

	// Ищем сообщение с таким же заголовком
	for _, message := range messages {
		for _, embed := range message.Embeds {
			if embed.Title == title {
				mId = message.ID
				break // Нашли нужное сообщение, выходим из цикла
			}
		}
		if mId != "" {
			break // Выход из основного цикла тоже
		}
	}

	// Создаём новый эмбед
	embed := &discordgo.MessageEmbed{
		Title:     title,
		Color:     0x00ff00,
		Image:     &discordgo.MessageEmbedImage{URL: imageUrl},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	if mId == "" {
		// Если сообщение не найдено, создаём новое
		_, err = d.S.ChannelMessageSendEmbed(channelID, embed)
	} else {
		// Если заголовок совпадает, просто обновляем сообщение
		_, err = d.S.ChannelMessageEditEmbed(channelID, mId, embed)
	}

	return err
}
func (d *Discord) SendOrEditEmbedImageFileName(channelID, title, fileNameScoreboard string) error {
	messages, err := d.S.ChannelMessages(channelID, 5, "", "", "")
	if err != nil {
		return fmt.Errorf("ошибка получения сообщений: %v", err)
	}
	// Ищем сообщение с таким же заголовком
	for _, message := range messages {
		for _, embed := range message.Embeds {
			if embed.Title == title {
				d.DeleteMessage(channelID, message.ID)
				break // Нашли нужное сообщение, выходим из цикла
			}
		}
	}

	open, err := os.Open("docker/scoreboard/" + fileNameScoreboard)
	if err != nil {
		return err
	}
	file := &discordgo.File{
		Name:   fileNameScoreboard,
		Reader: open,
	}
	dFile := []*discordgo.File{file}

	// Создаём новый эмбед
	embed := &discordgo.MessageEmbed{
		Title:     title,
		Color:     0x00ff00,
		Image:     &discordgo.MessageEmbedImage{URL: "attachment://" + fileNameScoreboard},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	wparams := &discordgo.WebhookParams{
		AvatarURL: "https://cdn.discordapp.com/icons/700238199070523412/9ddd77ab4f137a64a0cd74019318cd34.png",
		Username:  "Scoreboard",
		Embeds:    []*discordgo.MessageEmbed{embed},
		Files:     dFile,
	}

	_, err = d.webhook.Send(channelID, wparams)
	if err != nil {
		return err
	}

	return nil
}

func (d *Discord) SendChannelDelSecond(chatid, text string, second int) {
	if text != "" {
		text = d.re.replaceNameForDiscord(text)
		message, err := d.S.ChannelMessageSend(chatid, text)
		if err != nil {
			d.log.ErrorErr(err)
			d.log.Info(chatid + " " + text)
			return
		}
		tu := int(time.Now().UTC().Unix())
		d.storage.Db.TimerInsert(models.Timer{
			Tip:    "ds",
			ChatId: chatid,
			MesId:  message.ID,
			Timed:  tu + second,
		})
	}
}
func (d *Discord) SendComplexContent(chatid, text string) (mesId string) { //отправка текста комплексного сообщения
	mesCompl, err := d.S.ChannelMessageSendComplex(chatid, &discordgo.MessageSend{
		Content: text})
	if err != nil {
		channel, _ := d.S.Channel(chatid)
		d.log.Info("Ошибка отправки комплексного сообщения text " + channel.Name)
		d.log.ErrorErr(err)
		mesCompl, err = d.S.ChannelMessageSendComplex(chatid, &discordgo.MessageSend{
			Content: text})
		if err == nil {
			return mesCompl.ID
		}
		return ""
	}
	return mesCompl.ID
}
func (d *Discord) SendComplex(chatid string, mapEmbeds map[string]string) (mesId string) { //отправка текста комплексного сообщения
	mesCompl, err := d.S.ChannelMessageSendComplex(chatid, &discordgo.MessageSend{
		Content:    mesContentNil,
		Embed:      d.embedDS(mapEmbeds),
		Components: d.addButtonsQueue(mapEmbeds["buttonLevel"]),
	})
	if err != nil {
		channel, _ := d.S.Channel(chatid)
		d.log.Info("Ошибка отправки комплексного сообщения embed " + channel.Name)
		d.log.ErrorErr(err)
		mesCompl, err = d.S.ChannelMessageSendComplex(chatid, &discordgo.MessageSend{
			Content: mesContentNil,
			Embed:   d.embedDS(mapEmbeds),
		})
		if err == nil {
			return mesCompl.ID
		}
		return ""
	}
	return mesCompl.ID
}
func (d *Discord) Send(chatid, text string) (mesId string) { //отправка текста
	message, err := d.S.ChannelMessageSend(chatid, text)
	if err != nil {
		d.log.ErrorErr(err)
	}
	return message.ID
}

//func (d *Discord) SendEmbedTime1(chatid, text string) (mesId string) { //отправка текста с двумя реакциями
//	message, err := d.S.ChannelMessageSend(chatid, text)
//	if err != nil {
//		d.log.ErrorErr(err)
//	}
//	err = d.S.MessageReactionAdd(chatid, message.ID, emPlus)
//	if err != nil {
//		d.log.ErrorErr(err)
//	}
//	err = d.S.MessageReactionAdd(chatid, message.ID, emMinus)
//	if err != nil {
//		d.log.ErrorErr(err)
//	}
//	return message.ID
//}

func (d *Discord) SendEmbedTime(chatid, text string) (mesId string) { //отправка текста с двумя реакциями
	message, err := d.S.ChannelMessageSendComplex(chatid, &discordgo.MessageSend{
		Content: text,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Style:    discordgo.SecondaryButton,
						Label:    "+",
						CustomID: "+",
						Emoji: &discordgo.ComponentEmoji{
							Name: emPlus},
					},

					&discordgo.Button{
						Style:    discordgo.SecondaryButton,
						Label:    "-",
						CustomID: "-",
						Emoji: &discordgo.ComponentEmoji{
							Name: emMinus},
					}}},
		},
	})
	if err != nil {
		d.log.ErrorErr(err)
		return ""
	}
	return message.ID
}

func (d *Discord) SendWebhook(text, username, chatid, Avatar string) (mesId string) {
	if text == "" {
		return ""
	}
	pp := discordgo.WebhookParams{
		Content:   text,
		Username:  username,
		AvatarURL: Avatar,
	}
	mes, err := d.webhook.Send(chatid, &pp)
	if err != nil {
		fmt.Println(err)
		m := d.Send(chatid, text)
		return m
	}
	return mes.ID
}

//func (d *Discord) SendWebhookAsync(text, username, chatID, guildID, avatarURL string, resultChannel chan<- models.MessageIds, wg *sync.WaitGroup) {
//	defer wg.Done()
//
//	if text == "" {
//		return
//	}
//
//	web := transmitter.New(d.S, guildID, "KzBot", true, d.log)
//	params := &discordgo.WebhookParams{
//		Content:   text,
//		Username:  username,
//		AvatarURL: avatarURL,
//	}
//	mes, err := web.Send(chatID, params)
//	if err != nil {
//		fmt.Println(err)
//		d.Send(chatID, text) // Если вебхук не отправился, отправляем через обычное сообщение
//		return
//	}
//
//	messageData := models.MessageIds{
//		MessageId: mes.ID,
//		ChatId:    chatID,
//	}
//
//	resultChannel <- messageData
//}
//
//func (d *Discord) SendWebhookReplyAsync(text, username, chatid, guildId, Avatar string, reply *models.BridgeMessageReply, resultChannel chan<- models.MessageIds, wg *sync.WaitGroup) {
//	defer wg.Done()
//
//	if text == "" {
//		return
//	}
//	web := transmitter.New(d.S, guildId, "KzBot", true, d.log)
//	var embeds []*discordgo.MessageEmbed
//	e := discordgo.MessageEmbed{
//		Description: reply.Text,
//		Timestamp:   time.Unix(reply.TimeMessage, 0).Format(time.RFC3339),
//		Color:       14232643,
//		Author: &discordgo.MessageEmbedAuthor{
//			Name:    reply.UserName,
//			IconURL: reply.Avatar,
//		},
//	}
//
//	embeds = append(embeds, &e)
//
//	pp := &discordgo.WebhookParams{
//		Content:   text,
//		Username:  username,
//		AvatarURL: Avatar,
//		Embeds:    embeds,
//	}
//	mes, err := web.Send(chatid, pp)
//	if err != nil {
//		d.log.ErrorErr(err)
//		d.Send(chatid, text)
//		return
//	}
//	messageData := models.MessageIds{
//		MessageId: mes.ID,
//		ChatId:    chatid,
//	}
//
//	resultChannel <- messageData
//}
//
//func (d *Discord) SendFileAsync(text, username, channelID, guildId, fileURL, Avatar string, resultChannel chan<- models.MessageIds, wg *sync.WaitGroup) {
//	defer wg.Done()
//	fileName, i := utils.Convert(fileURL)
//	// convert byte slice to io.Reader
//	reader := bytes.NewReader(i)
//
//	web := transmitter.New(d.S, guildId, "KzBot", true, d.log)
//
//	// Подготавливаем параметры вебхука
//	webhook := &discordgo.WebhookParams{
//		Content:   text,
//		Username:  username,
//		AvatarURL: Avatar,
//		Files: []*discordgo.File{{
//			Name:   fileName, // Имя файла, которое будет видно в Discord
//			Reader: reader,
//		},
//		},
//	}
//
//	// Отправляем файл в Discord
//	m, err := web.Send(channelID, webhook)
//	if err != nil {
//		return
//	}
//	messageData := models.MessageIds{
//		MessageId: m.ID,
//		ChatId:    channelID,
//	}
//
//	resultChannel <- messageData
//}
//
//func (d *Discord) SendFilePic(channelID string, f *bytes.Reader) {
//	_, err := d.S.ChannelFileSend(channelID, "image.png", f)
//	if err != nil {
//		d.log.ErrorErr(err)
//		return
//	}
//}

func (d *Discord) SendPic(channelID, text string, imageBytes []byte) (string, error) {
	// Отправляем сообщение с вложенным файлом (изображением)
	m, err := d.S.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content: text,
		Files: []*discordgo.File{
			{
				Name:   "image.jpg",
				Reader: bytes.NewReader(imageBytes),
			},
		},
	})
	if err != nil {
		return "", err
	}

	return m.ID, nil
}

//func (d *Discord) SendBridgeAsync(text, username string, channelID, fileURL []string, Avatar string, reply *models.BridgeMessageReply, resultChannel chan<- models.MessageIds, wg *sync.WaitGroup) {
//	web := transmitter.New(d.S, "", "KzBot", true, d.log)
//	params := &discordgo.WebhookParams{
//		Content:   text,
//		Username:  username,
//		AvatarURL: Avatar,
//	}
//	if reply != nil {
//		params.Embeds = append(params.Embeds, &discordgo.MessageEmbed{
//			Description: reply.Text,
//			Timestamp:   time.Unix(reply.TimeMessage, 0).Format(time.RFC3339),
//			Color:       14232643,
//			Author: &discordgo.MessageEmbedAuthor{
//				Name:    reply.UserName,
//				IconURL: reply.Avatar,
//			},
//		})
//	} else if len(fileURL) > 0 && fileURL[0] != "" {
//		for _, s := range fileURL {
//			resp, err := http.Get(s)
//			if err != nil {
//				return
//			}
//			defer resp.Body.Close()
//			if resp.StatusCode != http.StatusOK {
//				return
//			}
//			fileName := filepath.Base(s)
//			params.Files = append(params.Files, &discordgo.File{
//				Name:   fileName,
//				Reader: resp.Body})
//		}
//
//	}
//
//	for _, channelId := range channelID {
//		d.sendWebhookBridge(channelId, params, web, resultChannel, wg)
//	}
//}
//
//func (d *Discord) sendWebhookBridge(channelId string, webhook *discordgo.WebhookParams, web *transmitter.Transmitter, resultChannel chan<- models.MessageIds, wg *sync.WaitGroup) {
//	// Отправляем файл в Discord
//	m, err := web.Send(channelId, webhook)
//	if err != nil {
//		d.log.ErrorErr(err)
//		d.log.InfoStruct(channelId, webhook)
//		return
//	}
//	messageData := models.MessageIds{
//		MessageId: m.ID,
//		ChatId:    channelId,
//	}
//
//	resultChannel <- messageData
//	wg.Done()
//}
//
//func (d *Discord) SendBridgeFunc(text, username string, channelID []string, Extra map[string][]interface{}, Avatar string, reply *models.BridgeMessageReply, resultChannel chan<- models.MessageIds, wg *sync.WaitGroup) {
//	params := &discordgo.WebhookParams{
//		Content:   text,
//		Username:  username,
//		AvatarURL: Avatar,
//		Files:     []*discordgo.File{},
//	}
//	if reply != nil {
//		params.Embeds = append(params.Embeds, &discordgo.MessageEmbed{
//			Description: reply.Text,
//			Timestamp:   time.Unix(reply.TimeMessage, 0).Format(time.RFC3339),
//			Color:       14232643,
//			Author: &discordgo.MessageEmbedAuthor{
//				Name:    reply.UserName,
//				IconURL: reply.Avatar,
//			},
//		})
//	} else if len(Extra) > 0 {
//		for _, f := range Extra["file"] {
//			fi := f.(models.FileInfo)
//			file := discordgo.File{
//				Name:        fi.Name,
//				ContentType: "",
//				Reader:      bytes.NewReader(fi.Data),
//			}
//			params.Files = append(params.Files, &file)
//		}
//	}
//
//	for _, channelId := range channelID {
//		// Отправляем файл в Discord
//		m, err := d.webhook.Send(channelId, params)
//		if err != nil {
//			d.log.ErrorErr(err)
//			return
//		}
//		messageData := models.MessageIds{
//			MessageId: m.ID,
//			ChatId:    channelId,
//		}
//
//		resultChannel <- messageData
//		wg.Done()
//	}
//}

func (d *Discord) SendBridgeFuncRest(in models.BridgeSendToMessenger) []models.MessageIds {

	var message []models.MessageIds

	for _, channelId := range in.ChannelId {
		params := &discordgo.WebhookParams{
			Content:   d.re.replaceNameForDiscord(in.Text),
			Username:  in.Sender,
			AvatarURL: in.Avatar,
			Files:     []*discordgo.File{},
		}
		if in.Reply != nil {
			params.Embeds = append(params.Embeds, &discordgo.MessageEmbed{
				Description: d.re.replaceNameForDiscord(in.Reply.Text),
				Timestamp:   time.Unix(in.Reply.TimeMessage, 0).Format(time.RFC3339),
				Color:       14232643,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    in.Reply.UserName,
					IconURL: in.Reply.Avatar,
				},
			})
		}
		if len(in.Extra) > 0 {
			for _, f := range in.Extra {
				contentType := mime.TypeByExtension(filepath.Ext(f.Name))
				file := discordgo.File{
					Name:        f.Name,
					ContentType: contentType,
				}
				if f.URL == "" && len(f.Data) == 0 {
					continue
				}
				if len(f.Data) > 0 {
					file.Reader = bytes.NewReader(f.Data)
				} else if f.URL != "" {
					downloadFile, err := helpers.DownloadFile(f.URL)
					if err != nil {
						d.log.ErrorErr(err)
						return nil
					}
					file.Reader = bytes.NewReader(downloadFile)
				}
				params.Files = append(params.Files, &file)
			}
		}
		if in.ReplyMap[channelId] != "" && len(in.Extra) == 0 {
			replyId := in.ReplyMap[channelId]
			text := fmt.Sprintf("%s:\n%s", in.Sender, in.Text)
			sendReply, err := d.S.ChannelMessageSendReply(channelId, text, &discordgo.MessageReference{
				MessageID: replyId,
				ChannelID: channelId,
			})
			if err != nil {
				d.log.ErrorErr(err)
			} else {
				messageData := models.MessageIds{
					MessageId: sendReply.ID,
					ChatId:    channelId,
				}
				message = append(message, messageData)
				continue
			}
		}
		m, err := d.webhook.Send(channelId, params)
		if err != nil {
			d.log.ErrorErr(err)
			restErr, _ := err.(*discordgo.RESTError)
			if restErr != nil {
				d.log.ErrorErr(restErr)
			}
			if restErr != nil && restErr.Message != nil && restErr.Message.Code == discordgo.ErrCodeUnknownChannel {
				d.log.Info("нужно сделать удаление этого канала : " + channelId)
			} else {
				d.log.ErrorErr(err)
			}
		} else {
			messageData := models.MessageIds{
				MessageId: m.ID,
				ChatId:    channelId,
			}
			message = append(message, messageData)
		}
	}
	return message
}

func (d *Discord) SendPoll(data map[string]string, options []string) string {
	chatid := data["chatid"]
	question := data["question"]
	url := data["url"]
	createTime := data["createTime"]
	description := ""
	for i, option := range options {
		description += fmt.Sprintf("\n%d. %s", i+1, option)
	}
	title := fmt.Sprintf("Опрос от %s: \n  %s", data["author"], question)
	Emb := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  16711680,
		Title:  title,
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:  description,
				Value: fmt.Sprintf("[результат](%s)", url),
			},
		},
	}
	fmt.Println("createTime ", createTime)
	mes, err := d.S.ChannelMessageSendComplex(chatid, &discordgo.MessageSend{
		Components: d.AddButtonPoll(createTime, options),
		Embed:      Emb,
	})
	if err != nil {
		d.log.ErrorErr(err)
		return ""
	}

	return mes.ID
}

const (
	emOne  = "1️⃣"
	emTwo  = "2️⃣"
	emTree = "3️⃣"
	emFour = "4️⃣"
	emFive = "5️⃣"
)

func (d *Discord) AddButtonPoll(createTime string, option []string) []discordgo.MessageComponent {
	var components []discordgo.MessageComponent
	if len(option) > 0 {
		if len(option) > 0 && option[0] != "" {
			button := discordgo.Button{
				Style:    discordgo.SecondaryButton,
				Label:    "",
				CustomID: createTime + ".1",
				Emoji: &discordgo.ComponentEmoji{
					Name: emOne,
				},
			}
			components = append(components, button)
		}
		if len(option) > 1 && option[1] != "" {
			button := discordgo.Button{
				Style:    discordgo.SecondaryButton,
				Label:    "",
				CustomID: createTime + ".2",
				Emoji: &discordgo.ComponentEmoji{
					Name: emTwo,
				},
			}
			components = append(components, button)
		}
		if len(option) > 2 && option[2] != "" {
			button := discordgo.Button{
				Style:    discordgo.SecondaryButton,
				Label:    "",
				CustomID: createTime + ".3",
				Emoji: &discordgo.ComponentEmoji{
					Name: emTree,
				},
			}
			components = append(components, button)
		}

		if len(option) > 3 && option[3] != "" {
			button := discordgo.Button{
				Style:    discordgo.SecondaryButton,
				Label:    "",
				CustomID: createTime + ".4",
				Emoji: &discordgo.ComponentEmoji{
					Name: emFour,
				},
			}
			components = append(components, button)
		}
		if len(option) > 4 && option[4] != "" {
			button := discordgo.Button{
				Style:    discordgo.SecondaryButton,
				Label:    "",
				CustomID: createTime + ".5",
				Emoji: &discordgo.ComponentEmoji{
					Name: emFive,
				},
			}
			components = append(components, button)
		}
	}
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: components,
		},
	}
}
