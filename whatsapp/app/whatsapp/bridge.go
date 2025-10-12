package wa

import (
	"fmt"
	"strings"
	"whatsapp/models"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

func (b *Whatsapp) handleBridgeMessage(message *events.Message, conf models.Bridge2Config) {
	msg := message.Message

	switch {
	case msg.Conversation != nil || msg.ExtendedTextMessage != nil:
		b.handleBridgeTextMessage(message.Info, msg, nil, &conf)
	case msg.ImageMessage != nil, msg.VideoMessage != nil, msg.PtvMessage != nil, msg.AudioMessage != nil, msg.DocumentMessage != nil:
		b.handleBridgeMediaMessage(message, conf)
	case msg.PollCreationMessageV3 != nil:
		b.handleBridgePollMessage(message)
	case msg.ProtocolMessage != nil:
		b.handleBridgeProtocolMessage(message.Info, msg, &conf)
	case msg.LocationMessage != nil, msg.ContactMessage != nil, msg.MessageContextInfo != nil:
		b.handleTestMessage(message, conf)

	//case msg.ProtocolMessage != nil && *msg.ProtocolMessage.Type == proto.ProtocolMessage_REVOKE:
	// 	b.handleDelete(msg.ProtocolMessage)
	default:
		fmt.Printf("Receiving message %+v\n", msg)
		return
	}
}

func (b *Whatsapp) handleBridgeProtocolMessage(messageInfo types.MessageInfo, msg *waE2E.Message, conf *models.Bridge2Config) {
	pr := msg.GetProtocolMessage()
	senderJID := messageInfo.Sender
	if pr.GetType() == waE2E.ProtocolMessage_REVOKE {

		rmsg := models.ToBridgeMessage{
			Tip:    "delWa",
			ChatId: pr.GetKey().GetRemoteJID(),
			MesId:  getMessageIdFormat(senderJID, messageInfo.ID),
			Config: conf,
		}
		b.api.SendBridgeAppRecover(rmsg)
		return
	}
	if pr.GetType() == waE2E.ProtocolMessage_MESSAGE_EDIT {
		if pr.GetKey() != nil && pr.GetEditedMessage() != nil {
			rmsg := models.ToBridgeMessage{
				Text:   pr.GetEditedMessage().GetConversation(),
				Sender: b.getSenderName(messageInfo),
				Tip:    "wae",
				ChatId: pr.GetKey().GetRemoteJID(),
				MesId:  getMessageIdFormat(senderJID, pr.GetKey().GetID()),
				Config: conf,
			}
			b.api.SendBridgeAppRecover(rmsg)
			return
		}
	}
}

// handleBridgeMediaMessage получает FileInfo и передает его дальше
func (b *Whatsapp) handleTestMessage(msg *events.Message, conf models.Bridge2Config) {
	var mediaMsg whatsmeow.DownloadableMessage
	//var mediaType string
	var captionText string

	switch {
	case msg.Message.LocationMessage != nil:
		fmt.Printf("GetLocationMessage %+v\n", msg.Message.GetLocationMessage())
		return
	case msg.Message.ContactMessage != nil:
		fmt.Printf("GetContactMessage %+v\n", msg.Message.GetContactMessage())
		return
	case msg.Message.MessageContextInfo != nil:
		fmt.Printf("GetMessageContextInfo %+v\n", msg.Message.GetMessageContextInfo())
		return
	case msg.Message.ImageMessage != nil:
		mediaMsg = msg.Message.GetImageMessage()
		//mediaType = "image"
		captionText = msg.Message.GetImageMessage().GetCaption()
	case msg.Message.VideoMessage != nil:
		mediaMsg = msg.Message.GetVideoMessage()
		//mediaType = "video"
		captionText = msg.Message.GetVideoMessage().GetCaption()
	case msg.Message.PtvMessage != nil:
		mediaMsg = msg.Message.GetPtvMessage()
		captionText = msg.Message.GetPtvMessage().GetCaption()
	case msg.Message.AudioMessage != nil:
		mediaMsg = msg.Message.GetAudioMessage()
		//mediaType = "audio"
	case msg.Message.DocumentMessage != nil:
		mediaMsg = msg.Message.GetDocumentMessage()
		//mediaType = "document"
		captionText = msg.Message.GetDocumentMessage().GetCaption()
	default:
		b.log.InfoStruct("Unknown media message type in handleMediaMessage: ", msg.Message)
		return
	}

	fileInfo, err := b.handleAttachment(msg, mediaMsg)
	if err != nil {
		b.log.ErrorErr(err)
		return
	}
	fileInfo.Comment = captionText

	// Передаем *FileInfo для добавления в Extra
	b.handleBridgeTextMessage(msg.Info, msg.Message, fileInfo, &conf)
}

// handleBridgeTextMessage принимает *FileInfo для добавления в Extra
func (b *Whatsapp) handleBridgeTextMessage(messageInfo types.MessageInfo, msg *waE2E.Message, fileInfo *models.FileInfo, conf *models.Bridge2Config) {
	senderJID := messageInfo.Sender
	channel := messageInfo.Chat

	senderName := b.getSenderName(messageInfo)

	isTextMessage := msg.GetExtendedTextMessage() != nil || msg.GetConversation() != ""
	if !isTextMessage && fileInfo == nil {
		b.log.InfoStruct("message without text content or attachment? ", msg)
		return
	}

	var text string

	if msg.GetExtendedTextMessage() == nil {
		text = msg.GetConversation()
	} else {
		extendedText := msg.GetExtendedTextMessage()
		text = extendedText.GetText()
		ci := extendedText.GetContextInfo()

		if senderJID == (types.JID{}) && ci.Participant != nil {
			senderJID = types.NewJID(ci.GetParticipant(), types.DefaultUserServer)
		}

		if ci.MentionedJID != nil {
			for _, mentionedJID := range ci.MentionedJID {
				numberAndSuffix := strings.SplitN(mentionedJID, "@", 2)
				mention := b.getSenderNotify(types.NewJID(numberAndSuffix[0], types.DefaultUserServer))
				text = strings.Replace(text, "@"+numberAndSuffix[0], "@"+mention, 1)
			}
		}
	}

	rmsg := models.ToBridgeMessage{
		Text:          text,
		Sender:        senderName,
		Tip:           "wa",
		ChatId:        channel.String(),
		MesId:         getMessageIdFormat(senderJID, messageInfo.ID),
		GuildId:       channel.String(),
		TimestampUnix: messageInfo.Timestamp.Unix(),
		Reply:         nil,
		Config:        conf,
	}

	if fileInfo != nil {
		// Добавляем вложение в срез Extra
		rmsg.Extra = append(rmsg.Extra, *fileInfo)
		rmsg.Text = fmt.Sprintf("%s\n%s", rmsg.Text, fileInfo.Comment)
	}

	if avatarURL, exists := b.userAvatars[senderJID.String()]; exists {
		rmsg.Avatar = avatarURL
	}

	b.api.SendBridgeAppRecover(rmsg)
}

// handleBridgeMediaMessage получает FileInfo и передает его дальше
func (b *Whatsapp) handleBridgeMediaMessage(msg *events.Message, conf models.Bridge2Config) {
	var mediaMsg whatsmeow.DownloadableMessage
	//var mediaType string
	var captionText string

	switch {
	case msg.Message.ImageMessage != nil:
		mediaMsg = msg.Message.GetImageMessage()
		//mediaType = "image"
		captionText = msg.Message.GetImageMessage().GetCaption()
	case msg.Message.VideoMessage != nil:
		mediaMsg = msg.Message.GetVideoMessage()
		//mediaType = "video"
		captionText = msg.Message.GetVideoMessage().GetCaption()
	case msg.Message.PtvMessage != nil:
		mediaMsg = msg.Message.GetPtvMessage()
		captionText = msg.Message.GetPtvMessage().GetCaption()
	case msg.Message.AudioMessage != nil:
		mediaMsg = msg.Message.GetAudioMessage()
		//mediaType = "audio"
	case msg.Message.DocumentMessage != nil:
		mediaMsg = msg.Message.GetDocumentMessage()
		//mediaType = "document"
		captionText = msg.Message.GetDocumentMessage().GetCaption()
	default:
		b.log.InfoStruct("Unknown media message type in handleMediaMessage: ", msg.Message)
		return
	}

	fileInfo, err := b.handleAttachment(msg, mediaMsg)
	if err != nil {
		b.log.ErrorErr(err)
		return
	}
	fileInfo.Comment = captionText

	// Передаем *FileInfo для добавления в Extra
	b.handleBridgeTextMessage(msg.Info, msg.Message, fileInfo, &conf)
}

func (b *Whatsapp) handleBridgePollMessage(message *events.Message) {
	pollMsg := message.Message.GetPollCreationMessageV3()
	if pollMsg == nil {
		return // Дополнительная проверка
	}

	// 1. Извлечение данных опроса
	question := pollMsg.GetName()
	options := make([]string, len(pollMsg.GetOptions()))
	for i, opt := range pollMsg.GetOptions() {
		options[i] = opt.GetOptionName()
	}

	// 0 = Multiple Select, 1 = Single Select.
	// В whatsmeow 0 означает, что можно выбрать любое количество.
	multiSelect := pollMsg.GetSelectableOptionsCount() == 0

	// 2. Создание базового сообщения
	senderJID := message.Info.Sender

	rmsg := models.ToBridgeMessage{
		Text:          question, // Вопрос часто используется как основной текст
		Sender:        b.getSenderName(message.Info),
		Tip:           "wa_poll",
		ChatId:        message.Info.Chat.String(),
		MesId:         getMessageIdFormat(senderJID, message.Info.ID),
		TimestampUnix: message.Info.Timestamp.Unix(),
		Reply:         nil,
		Config:        nil,
	}

	// 3. Добавление информации об опросе в Extra
	// Поскольку Extra - это []FileInfo, мы вынуждены использовать его.
	// Если ваша модель ToBridgeMessage может принимать map[string]interface{} в Extra,
	// используйте лучше его.

	// Если Extra должно быть []FileInfo, вы должны решить, как преобразовать PollInfo в FileInfo.
	// Здесь мы просто добавим строку с метаданными как "файл"

	pollDataJSON := fmt.Sprintf(`{"type": "poll", "question": "%s", "options": ["%s"], "multiSelect": %t}`,
		question, strings.Join(options, `", "`), multiSelect)

	pollFileInfo := models.FileInfo{
		Name:   "poll_metadata.json", // Имя файла для метаданных
		Data:   []byte(pollDataJSON),
		Size:   int64(len(pollDataJSON)),
		FileID: getMessageIdFormat(senderJID, message.Info.ID),
		// URL остается пустым
	}

	rmsg.Extra = append(rmsg.Extra, pollFileInfo)

	fmt.Printf("Обработанное сообщение-опрос: %+v\n", rmsg)
}

func (b *Whatsapp) SendForBridge(in models.BridgeSendToMessenger) []models.MessageIds {
	var res []models.MessageIds
	for _, s := range in.ChannelId {
		m := Message{
			Text:     in.Text,
			Channel:  s,
			Username: in.Sender,
			Avatar:   in.Avatar,
			Extra:    in.Extra,
		}
		if in.ReplyMap != nil && in.ReplyMap[s] != "" {
			m.ParentID = in.ReplyMap[s]
		}
		send, err := b.Send(m)
		if err != nil {
			b.log.ErrorErr(err)
		}
		res = append(res, models.MessageIds{
			MessageId: send,
			ChatId:    s,
		})
	}
	return res
}

func (b *Whatsapp) DeleteMessage(ChatId, mId string) error {
	groupJID, _ := b.ParseJID(ChatId)

	extendedMsgID, _ := b.parseMessageID(mId)
	ID := extendedMsgID.MessageID

	_, err := b.wc.RevokeMessage(groupJID, ID)
	if err != nil {
		return err
	}
	return nil
}
