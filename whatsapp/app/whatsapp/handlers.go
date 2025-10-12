package wa

import (
	"context"
	"fmt"
	"mime"
	"path/filepath"
	"strings"
	"time"
	"whatsapp/models"

	"github.com/42wim/matterbridge/bridge/config"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

func (b *Whatsapp) eventHandler(evt interface{}) {
	switch e := evt.(type) {
	case *events.Message:
		b.handleMessage(e)
	case *events.GroupInfo:
		b.handleGroupInfo(e)
	case *events.DeleteForMe:
		fmt.Printf("events.DeleteForMe %+v", e)
	case *events.DeleteChat:
		fmt.Printf("events.DeleteChat %+v", e)
	case types.MessageSource:
		fmt.Printf("MessageSource %+v\n", e)
	case *types.MessageSource:
		fmt.Printf("MessageSource* %+v\n", e)

	default:
		fmt.Printf("eventHandler %+v\n", e)

	}
}

func (b *Whatsapp) handleMessage(message *events.Message) {
	msg := message.Message

	// Фильтрация: nil, от меня, старые сообщения
	if msg == nil || message.Info.IsFromMe || message.Info.Timestamp.Before(b.startedAt) {
		return
	}

	b.filter(message)

	//bridge
	wa, bridgeConfig := b.BridgeCheckChannelConfigWA(message.Info.Chat.String())
	if wa {
		b.handleBridgeMessage(message, bridgeConfig)
		return
	}

	switch {
	case msg.Conversation != nil || msg.ExtendedTextMessage != nil:
		b.handleTextMessage(message.Info, msg, nil)

	//case msg.ProtocolMessage != nil && *msg.ProtocolMessage.Type == proto.ProtocolMessage_REVOKE:
	// 	b.handleDelete(msg.ProtocolMessage)
	default:
		fmt.Printf("Receiving message %+v\n", msg)
		return
	}
}

// handleTextMessage принимает *FileInfo для добавления в Extra
func (b *Whatsapp) handleTextMessage(messageInfo types.MessageInfo, msg *waE2E.Message, fileInfo *models.FileInfo) {
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
		GuildId:       "",
		TimestampUnix: messageInfo.Timestamp.Unix(),
		Reply:         nil,
		Config:        nil,
	}

	if strings.HasPrefix(rmsg.Text, ".") {
		g := b.getGroupCommunity(messageInfo)
		rmsg.Config = &models.Bridge2Config{
			HostRelay: g.GuildName,
		}
		b.api.SendBridgeAppRecover(rmsg)
	}

	if fileInfo != nil {
		// Добавляем вложение в срез Extra
		rmsg.Extra = append(rmsg.Extra, *fileInfo)
	}

	if avatarURL, exists := b.userAvatars[senderJID.String()]; exists {
		rmsg.Avatar = avatarURL
	}

	fmt.Printf("Обработанное сообщение: %+v\n", rmsg)
}

// handleMediaMessage получает FileInfo и передает его дальше
func (b *Whatsapp) handleMediaMessage(msg *events.Message) {
	var mediaMsg whatsmeow.DownloadableMessage
	//var mediaType string

	switch {
	case msg.Message.ImageMessage != nil:
		mediaMsg = msg.Message.GetImageMessage()
		//mediaType = "image"
	case msg.Message.VideoMessage != nil:
		mediaMsg = msg.Message.GetVideoMessage()
		//mediaType = "video"
	case msg.Message.AudioMessage != nil:
		mediaMsg = msg.Message.GetAudioMessage()
		//mediaType = "audio"
	case msg.Message.DocumentMessage != nil:
		mediaMsg = msg.Message.GetDocumentMessage()
		//mediaType = "document"
	default:
		b.log.InfoStruct("Unknown media message type in handleMediaMessage: ", msg.Message)
		return
	}

	fileInfo, err := b.handleAttachment(msg, mediaMsg)
	if err != nil {
		b.log.ErrorErr(err)
		return
	}

	// Передаем *FileInfo для добавления в Extra
	b.handleTextMessage(msg.Info, msg.Message, fileInfo)
}

// handleAttachment скачивает медиа, извлекает метаданные и возвращает *FileInfo
//func (b *Whatsapp) handleAttachment(msg *events.Message, media whatsmeow.DownloadableMessage, mediaType string) (*models.FileInfo, error) {
//	data, err := b.wc.Download(context.Background(), media)
//	if err != nil {
//		return nil, fmt.Errorf("download %s failed: %w", mediaType, err)
//	}
//
//	var filename string
//	var mimeType string
//	var size int64
//
//	// URL устанавливаем пустым, так как прямого URL нет.
//	// Если ваш мост генерирует публичный URL для скачанных данных, вставьте его здесь.
//	url := ""
//
//	// Используем утверждение типа для получения конкретных полей
//	switch v := media.(type) {
//	case *waE2E.DocumentMessage:
//		filename = v.GetFileName()
//		mimeType = v.GetMimetype()
//		size = int64(v.GetFileLength())
//	case *waE2E.ImageMessage:
//		mimeType = v.GetMimetype()
//		size = int64(v.GetFileLength())
//	case *waE2E.VideoMessage:
//		mimeType = v.GetMimetype()
//		size = int64(v.GetFileLength())
//	case *waE2E.AudioMessage:
//		mimeType = v.GetMimetype()
//		size = int64(v.GetFileLength())
//	default:
//		return nil, fmt.Errorf("unsupported media type passed to handleAttachment: %T", media)
//	}
//
//	// Логика определения имени файла (если имя не получено из DocumentMessage)
//	if filename == "" {
//		fileExt, _ := mime.ExtensionsByType(mimeType)
//		ext := ".bin"
//		if len(fileExt) > 0 {
//			ext = fileExt[0]
//		}
//		if ext == ".jfif" || ext == ".jpe" {
//			ext = ".jpg"
//		} else if mediaType == "video" && ext == "" {
//			ext = ".mp4"
//		} else if mediaType == "audio" && ext == "" {
//			ext = ".ogg"
//		}
//		filename = fmt.Sprintf("%v%v", msg.Info.ID, ext)
//	}
//
//	return &models.FileInfo{
//		Name: filename,
//		Data: data,
//		URL:  url, // Установлено выше как ""
//		Size: size,
//	}, nil
//}

// mediaType в сигнатуре функции не нужен, если вы используете утверждение типа.
func (b *Whatsapp) handleAttachment(msg *events.Message, media whatsmeow.DownloadableMessage) (*models.FileInfo, error) {
	// 1. Скачивание файла
	data, err := b.wc.Download(context.Background(), media)
	if err != nil {
		// В логировании mediaType не нужен, можно просто указать, что это медиа
		return nil, fmt.Errorf("download media failed: %w", err)
	}

	var filename string
	var mimeType string
	var size int64
	url := "" // URL остается пустым, если не генерируется мостом

	// 2. Извлечение метаданных
	switch v := media.(type) {
	case *waE2E.DocumentMessage:
		filename = v.GetFileName()
		mimeType = v.GetMimetype()
		size = int64(v.GetFileLength())
		fmt.Println("mimeTypeDocumentMessage ", mimeType)
	case *waE2E.ImageMessage:
		// ImageMessage и VideoMessage не имеют поля FileName
		mimeType = v.GetMimetype()
		size = int64(v.GetFileLength())
		fmt.Println("mimeTypeImageMessage ", mimeType)
	case *waE2E.VideoMessage:
		mimeType = v.GetMimetype()
		size = int64(v.GetFileLength())
		fmt.Println("mimeTypeVideoMessage ", mimeType)
	case *waE2E.AudioMessage:
		mimeType = v.GetMimetype()
		size = int64(v.GetFileLength())
		fmt.Println("mimeTypeAudioMessage ", mimeType)
	default:
		fmt.Println("mimeTypeDefaultMessage ", mimeType)
		return nil, fmt.Errorf("unsupported media type passed to handleAttachment: %T", media)
	}

	// 3. Генерация имени файла, если оно отсутствует
	if filename == "" {
		// Получаем все возможные расширения для данного MIME-типа
		fileExts, err := mime.ExtensionsByType(mimeType)
		ext := ""

		if err == nil && len(fileExts) > 0 {
			// Берем первое расширение и убираем начальную точку
			ext = strings.TrimPrefix(fileExts[0], ".")
		}

		// --- УСИЛЕНИЕ ЛОГИКИ ДЛЯ АУДИО (ИСПРАВЛЕНИЕ *.bin) ---
		// Если стандартное определение не дало результата или MIME-тип указывает на Ogg
		if ext == "" || mimeType == "audio/ogg" {

			// Проверяем, является ли сообщение голосовым (аудио)
			// Нам нужно определить, был ли media типом AudioMessage,
			// но поскольку мы не знаем, как именно это сделать без mediaType,
			// мы используем mimeType как основной индикатор.

			if strings.HasPrefix(mimeType, "audio/") {
				// Для WhatsApp-аудио часто используется Ogg Opus.
				if mimeType == "audio/ogg" || mimeType == "audio/opus" {
					ext = "ogg" // Или "opus"
				} else if mimeType == "audio/mp4" || mimeType == "audio/aac" {
					ext = "mp4"
				} else {
					ext = "mp3" // Запасной вариант для аудио
				}
			}
		}

		// Если расширение все еще не определено, используем 'bin'
		if ext == "" {
			ext = "bin"
		}

		// Формируем уникальное имя файла
		filename = fmt.Sprintf("%s-%d.%s",
			msg.Info.ID,
			time.Now().UnixNano()/int64(time.Millisecond),
			ext)

	} else {
		// Если имя файла получено (из DocumentMessage), убедимся, что для аудио есть .ogg
		if strings.HasPrefix(mimeType, "audio/") && filepath.Ext(filename) == "" {
			filename = filename + ".ogg"
		}
	}

	// 4. Возвращаем результат
	return &models.FileInfo{
		Name: filename,
		Data: data,
		URL:  url,
		Size: size,
	}, nil
}

func (b *Whatsapp) handleGroupInfo(event *events.GroupInfo) {
	fmt.Printf("Receiving event %#v", event)

	switch {
	case event.Join != nil:
		b.handleUserJoin(event)
	case event.Leave != nil:
		b.handleUserLeave(event)
	case event.Topic != nil:
		b.handleTopicChange(event)
	case event.Delete != nil:
		fmt.Printf("Delete event: %+v\n", event.Delete)

	default:
		fmt.Printf("Unknown event: %+v\n", event)
	}
}

func (b *Whatsapp) handleUserJoin(event *events.GroupInfo) {
	for _, joinedJid := range event.Join {
		senderName := b.getSenderNameFromJID(joinedJid)

		rmsg := config.Message{
			UserID:   joinedJid.String(),
			Username: senderName,
			Channel:  event.JID.String(),
			//Account:  b.Account,
			//Protocol: b.Protocol,
			Event: config.EventJoinLeave,
			Text:  "joined chat",
		}

		fmt.Printf("%+v\n", rmsg)
	}
}

func (b *Whatsapp) handleUserLeave(event *events.GroupInfo) {
	for _, leftJid := range event.Leave {
		senderName := b.getSenderNameFromJID(leftJid)

		rmsg := config.Message{
			UserID:   leftJid.String(),
			Username: senderName,
			Channel:  event.JID.String(),
			//Account:  b.Account,
			//Protocol: b.Protocol,
			Event: config.EventJoinLeave,
			Text:  "left chat",
		}

		fmt.Printf("%+v\n", rmsg)
	}
}

func (b *Whatsapp) handleTopicChange(event *events.GroupInfo) {
	msg := event.Topic
	senderJid := msg.TopicSetBy
	senderName := b.getSenderNameFromJID(senderJid)

	text := msg.Topic
	if text == "" {
		text = "removed topic"
	}

	rmsg := config.Message{
		UserID:   senderJid.String(),
		Username: senderName,
		Channel:  event.JID.String(),
		//Account:  b.Account,
		//Protocol: b.Protocol,
		Event: config.EventTopicChange,
		Text:  "Topic changed: " + text,
	}

	fmt.Printf("handleTopicChange %+v\n", rmsg)
}

/*
MessageIDs:[3EB0E1F0D6E0BA4B4B3A54 3EB050B9FF42E2648649B0 3EB066CDC17D9FC4D7C983 3EB0D9FC83BD73049BDCCD 3EB0B39B339BB0B7E27AD5 3EB0D9B6B9599432DB986D 3EB056E8C3B52E89C62CC6 3EB0093BE366152D769BD0 3EB09A93C2BE698D9E1D9E 3EB04D9E255C988311F6FE 3EB04D6742DC07EE2B4175 3EB0A0A5FBF418F46EF504 3EB0D115407D6FF2301EBD] Timestamp:2025-10-11 00:51:32 -0500 CDT Type:read MessageSender:}
MessageIDs:[3EB029A533BE2CFB7B3D22] Type:read MessageSender:}
MessageIDs:[3EB04D6742DC07EE2B4175] Type: MessageSender:}
MessageIDs:[3EB052937087C95DD5A5DC] Type:read MessageSender:}
MessageIDs:[3EB0E1F0D6E0BA4B4B3A54] Type: MessageSender:}
MessageIDs:[3EB0585328A52746579DCE] Type:read MessageSender:}
MessageIDs:[3EB0330DEE90EF6E9DC1AC] Type:read MessageSender:}
MessageIDs:[3EB09A93C2BE698D9E1D9E] Type: MessageSender:}
MessageIDs:[3EB04D9E255C988311F6FE] Type: MessageSender:}
MessageIDs:[3EB0B48427B18BC96DF788] Type:read MessageSender:}
MessageIDs:[3EB0EE456001916444172E] Type:read MessageSender:}
MessageIDs:[3EB064C67397BF9AF7DA8E] Type:read MessageSender:}
MessageIDs:[3EB00BF33D33E735AD4243] Type:read MessageSender:}
MessageIDs:[3EB0093BE366152D769BD0] Type: MessageSender:}
MessageIDs:[3EB0A5D97D1167501E6C1C] Type:read MessageSender:}
MessageIDs:[3EB0E7277CF942EE57945B] Type:read MessageSender:}
MessageIDs:[3EB0CB55DDB8C6CE44140C] Type:read MessageSender:}
MessageIDs:[3EB0CA37F94297CA86081F] Type:read MessageSender:}
MessageIDs:[3EB0D9B6B9599432DB986D] Type: MessageSender:}
MessageIDs:[3EB050B9FF42E2648649B0] Type: MessageSender:}
MessageIDs:[3EB0BD999155F61D2AB3A8] Type:read MessageSender:}
MessageIDs:[3EB0139C8ED48689590081] Type:read MessageSender:}
MessageIDs:[3EB0E887BD482240AF4F02] Type:read MessageSender:}
MessageIDs:[3EB0F35C7765647416E617] Type:read MessageSender:}
MessageIDs:[3EB051EFA5B9A9AC0415F0] Type:read MessageSender:}
MessageIDs:[3EB066CDC17D9FC4D7C983] Type: MessageSender:}
MessageIDs:[3EB081E2614AB6B927FE7A] Type:read MessageSender:}
MessageIDs:[3EB0FAD8F47668E221DB2D] Type:read MessageSender:}
MessageIDs:[3EB0C59BFE7E8059C1B0E9] Type:read MessageSender:}
MessageIDs:[3EB056E8C3B52E89C62CC6] Type: MessageSender:}
MessageIDs:[3EB0B39B339BB0B7E27AD5] Type: MessageSender:}
MessageIDs:[3EB03047737235E454E709] Type:read MessageSender:}
MessageIDs:[3EB0BC95AD2550BD248CB4] Type:read MessageSender:}
MessageIDs:[3EB0B23CF1CDF920FBF436] Type:read MessageSender:}
MessageIDs:[3EB08D8BFE0A158DF3F54F] Type:read MessageSender:}
MessageIDs:[3EB04645ADABF2441D1C5E] Type:read MessageSender:}
MessageIDs:[3EB03C47AB45522D80B918] Type:read MessageSender:}
MessageIDs:[3EB0D420752AC590250FBA] Type:read MessageSender:}
MessageIDs:[3EB0D9FC83BD73049BDCCD] Type: MessageSender:}
MessageIDs:[3EB02B7346CAB008513FBE] Type:read MessageSender:}
MessageIDs:[3EB03F6A5A21F64C80AD95] Type:read MessageSender:}
MessageIDs:[3EB0EDEBCDFE7EC87A2F44] Type:read MessageSender:}
MessageIDs:[3EB0F4E8789A30166FB1CB] Type:read MessageSender:}
MessageIDs:[3EB0908976A7797229E808] Type:read MessageSender:}
MessageIDs:[3EB0873B170130B9479AF4] Type:read MessageSender:}
MessageIDs:[3EB044F43F3D0A4AE31804] Type:read MessageSender:}
MessageIDs:[3EB0E7D634C9FBE867144D] Type:read MessageSender:}
MessageIDs:[3EB018A320A298F2EECF2C] Type:read MessageSender:}
Sender:278258197205177@lid MessageIDs:[3EB075F17C4002842397AD] Type:read MessageSender:}
Sender:278258197205177@lid MessageIDs:[3EB0712651A65E4231E9B3] Type:read MessageSender:}
Sender:278258197205177@lid MessageIDs:[3EB017A769FC1FAF83B66D] Type:read MessageSender:}
Sender:278258197205177@lid MessageIDs:[3EB067251CD7B3DD1386A0] Type:read MessageSender:}
Sender:278258197205177@lid MessageIDs:[3EB00EC3B36CB2EE80C162] Type:read MessageSender:}
Sender:278258197205177@lid MessageIDs:[3EB031D01239AA8E112EBE] Type:read MessageSender:}
Sender:278258197205177@lid MessageIDs:[3EB0D9B6B9599432DB986D] Type:played MessageSender:}
Sender:93574603485437@lid MessageIDs:[3EB0F66516F8BABF2768CF] Timestamp:2025-10-11 00:54:08 -0500 CDT Type: MessageSender:}
Sender:278258197205177@lid MessageIDs:[3EB0F66516F8BABF2768CF]T Type: MessageSender:}
Sender:1507617460237@lid MessageIDs:[3EB0F66516F8BABF2768CF] TType: MessageSender:}
Sender:101056335872247@lid MessageIDs:[3EB0F66516F8BABF2768CF]T Type: MessageSender:}
Sender:174032141983900@lid MessageIDs:[3EB0F66516F8BABF2768CF]T Type: MessageSender:}
Sender:21926445588488@lid MessageIDs:[3EB0F66516F8BABF2768CF]  Type: MessageSender:}
Sender:217278486474905@lid MessageIDs:[3EB0F66516F8BABF2768CF]T Type: MessageSender:}
Sender:140939452526807@lid MessageIDs:[3EB0F66516F8BABF2768CF]T Type: MessageSender:}
Sender:10213633605725@lid MessageIDs:[3EB0F66516F8BABF2768CF] Timestamp:2025-10-11 01:25:50 -0500 CDT Type: MessageSender:}
Sender:174032141983900@lid MessageIDs:[3EB0848F045D872E35DCFC 3EB0AB49A748D6CA738546 3EB07FCAC400AE2036CAFF 3EB04E4EEC1A51782E1380] Timestamp:2025-10-11 02:33:21 -0500 CDT Type:read MessageSender:}
Sender:230004155760702@lid MessageIDs:[3EB0F66516F8BABF2768CF] Timestamp:2025-10-11 02:54:55 -0500 CDT Type: MessageSender:}
Sender:217278486474905:2@lid MessageIDs:[3EB0F66516F8BABF2768CF] Timestamp:2025-10-11 03:20:58 -0500 CDT Type: MessageSender:}
Sender:93574603485437@lid MessageIDs:[3EB0712651A65E4231E9B3 3EB018A320A298F2EECF2C 3EB075F17C4002842397AD 3EB067251CD7B3DD1386A0 3EB017A769FC1FAF83B66D 3EB031D01239AA8E112EBE 3EB00EC3B36CB2EE80C162 3EB0BC95AD2550BD248CB4 3EB0CB55DDB8C6CE44140C 3EB03047737235E454E709 3EB0BD999155F61D2AB3A8 3EB051EFA5B9A9AC0415F0 3EB044F43F3D0A4AE31804 3EB0D420752AC590250FBA 3EB0E887BD482240AF4F02 3EB0908976A7797229E808 3EB0E7277CF942EE57945B 3EB081E2614AB6B927FE7A 3EB0E7D634C9FBE867144D 3EB08D8BFE0A158DF3F54F 3EB0F35C7765647416E617 3EB052937087C95DD5A5DC 3EB0F4E8789A30166FB1CB 3EB0139C8ED48689590081 3EB0FAD8F47668E221DB2D 3EB029A533BE2CFB7B3D22 3EB064C67397BF9AF7DA8E 3EB0A5D97D1167501E6C1C 3EB0EE456001916444172E 3EB0330DEE90EF6E9DC1AC 3EB0C59BFE7E8059C1B0E9 3EB00BF33D33E735AD4243 3EB04645ADABF2441D1C5E 3EB0CA37F94297CA86081F 3EB03C47AB45522D80B918 3EB0B48427B18BC96DF788 3EB0585328A52746579DCE 3EB03F6A5A21F64C80AD95 3EB0EDEBCDFE7EC87A2F44 3EB0873B170130B9479AF4 3EB02B7346CAB008513FBE 3EB0B23CF1CDF920FBF436 3EB0E1F0D6E0BA4B4B3A54 3EB050B9FF42E2648649B0 3EB066CDC17D9FC4D7C983] Timestamp:2025-10-11 06:15:57 -0500 CDT Type:read MessageSender:}
eventHandler &{Reason: Type:new CreateKey:79103982920-3f7553bde24b41bd9cff8832bfc27198@temp Sender:93574603485437@lid SenderPN:79103982920@s.whatsapp.net Notify:CTALKEP1001 GroupInfo:{JID:120363422464172527@g.us OwnerJID:93574603485437@lid OwnerPN:79103982920@s.whatsapp.net GroupName:{Name:Кз ивент счет NameSetAt:2025-10-11 06:16:35 -0500 CDT NameSetBy:93574603485437@lid NameSetByPN:79103982920@s.whatsapp.net} GroupTopic:{Topic: TopicID: TopicSetAt:0001-01-01 00:00:00 +0000 UTC TopicSetBy: TopicSetByPN: TopicDeleted:false} GroupLocked:{IsLocked:false} GroupAnnounce:{IsAnnounce:false AnnounceVersionID:1760181395769729} GroupEphemeral:{IsEphemeral:false DisappearingTimer:0} GroupIncognito:{IsIncognito:false} GroupParent:{IsParent:false DefaultMembershipApprovalMode:} GroupLinkedParent:{LinkedParentJID:120363200871949928@g.us} GroupIsDefaultSub:{IsDefaultSubGroup:false} GroupMembershipApprovalMode:{IsJoinApprovalRequired:false} AddressingMode: GroupCreated:2025-10-11 06:16:35 -0500 CDT CreatorCountryCode:RU ParticipantVersionID:1760181395939180 Participants:[{JID:85178361896964@lid PhoneNumber:79991399754@s.whatsapp.net LID:85178361896964@lid IsAdmin:false IsSuperAdmin:false DisplayName: Error:0 AddRequest:<nil>} {JID:93574603485437@lid PhoneNumber:79103982920@s.whatsapp.net LID:93574603485437@lid IsAdmin:true IsSuperAdmin:true DisplayName: Error:0 AddRequest:<nil>}] MemberAddMode:all_member_add}}
Receiving event &events.GroupInfo{JID:types.JID{User:"120363200871949928", Notify:"CTALKEP1001", Sender:(*types.JID)(0xc000645c80), SenderPN:(*types.JID)(0xc000645cb0), Timestamp:time.Date(2025, time.October, 11, 6, 16, 36, 0, time.Local), Name:(*types.GroupName)(nil), Topic:(*types.GroupTopic)(nil), Locked:(*types.GroupLocked)(nil), Announce:(*types.GroupAnnounce)(nil), Ephemeral:(*types.GroupEphemeral)(nil), MembershipApprovalMode:(*types.GroupMembershipApprovalMode)(nil), Delete:(*types.GroupDelete)(nil), Link:(*types.GroupLinkChange)(0xc00027dc70), Unlink:(*types.GroupLinkChange)(nil), NewInviteLink:(*string)(nil), PrevParticipantVersionID:"", ParticipantVersionID:"", JoinReason:"", Join:[]types.JID(nil), Leave:[]types.JID(nil), Promote:[]types.JID(nil), Demote:[]types.JID(nil), UnknownChanges:[]*binary.Node(nil)}Unknown event: &{JID:120363200871949928@g.us Notify:CTALKEP1001 Sender:93574603485437@lid SenderPN:79103982920@s.whatsapp.net Timestamp:2025-10-11 06:16:36 -0500 CDT Name:<nil> Topic:<nil> Locked:<nil> Announce:<nil> Ephemeral:<nil> MembershipApprovalMode:<nil> Delete:<nil> Link:0xc00027dc70 Unlink:<nil> NewInviteLink:<nil> PrevParticipantVersionID: ParticipantVersionID: JoinReason: Join:[] Leave:[] Promote:[] Demote:[] UnknownChanges:[]}
Receiving event &events.GroupInfo{JID:types.JID{User:"120363202366842960", Notify:"CTALKEP1001", Sender:(*types.JID)(0xc000483530), SenderPN:(*types.JID)(0xc000483560), Timestamp:time.Date(2025, time.October, 11, 6, 16, 36, 0, time.Local), Name:(*types.GroupName)(nil), Topic:(*types.GroupTopic)(nil), Locked:(*types.GroupLocked)(nil), Announce:(*types.GroupAnnounce)(nil), Ephemeral:(*types.GroupEphemeral)(nil), MembershipApprovalMode:(*types.GroupMembershipApprovalMode)(nil), Delete:(*types.GroupDelete)(nil), Link:(*types.GroupLinkChange)(0xc000504dd0), Unlink:(*types.GroupLinkChange)(nil), NewInviteLink:(*string)(nil), PrevParticipantVersionID:"", ParticipantVersionID:"", JoinReason:"", Join:[]types.JID(nil), Leave:[]types.JID(nil), Promote:[]types.JID(nil), Demote:[]types.JID(nil), UnknownChanges:[]*binary.Node(nil)}Unknown event: &{JID:120363202366842960@g.us Notify:CTALKEP1001 Sender:93574603485437@lid SenderPN:79103982920@s.whatsapp.net Timestamp:2025-10-11 06:16:36 -0500 CDT Name:<nil> Topic:<nil> Locked:<nil> Announce:<nil> Ephemeral:<nil> MembershipApprovalMode:<nil> Delete:<nil> Link:0xc000504dd0 Unlink:<nil> NewInviteLink:<nil> PrevParticipantVersionID: ParticipantVersionID: JoinReason: Join:[] Leave:[] Promote:[] Demote:[] UnknownChanges:[]}
Receiving event &events.GroupInfo{JID:types.JID{User:"120363422464172527", Notify:"WhatsApp", Sender:(*types.JID)(nil), SenderPN:(*types.JID)(nil), Timestamp:time.Date(2025, time.October, 11, 8, 0, 10, 0, time.Local), Name:(*types.GroupName)(nil), Topic:(*types.GroupTopic)(nil), Locked:(*types.GroupLocked)(nil), Announce:(*types.GroupAnnounce)(nil), Ephemeral:(*types.GroupEphemeral)(nil), MembershipApprovalMode:(*types.GroupMembershipApprovalMode)(nil), Delete:(*types.GroupDelete)(nil), Link:(*types.GroupLinkChange)(nil), Unlink:(*types.GroupLinkChange)(nil), NewInviteLink:(*string)(nil), PrevParticipantVersionID:"1760181395939180", ParticipantVersionID:"1760187609812467", JoinReason:"linked_group_join", Join:[]types.JID{types.JID{User:"278258197205177", RawAgent:0x0, Device:0x0, Integrator:0x0, Server:"lid"}}, Leave:[]types.JID(nil), Promote:[]types.JID(nil), Demote:[]types.JID(nil), UnknownChanges:[]*binary.Node(nil)}{Text:joined chat Channel:120363422464172527@g.us Username:Юлия_Рефлекс UserID:278258197205177@lid Avatar: Account: Event:join_leave Protocol: Gateway: ParentID: Timestamp:0001-01-01 00:00:00 +0000 UTC ID: Extra:map[]}
Receiving event &events.GroupInfo{JID:types.JID{User:"120363422464172527", Notify:"WhatsApp", Sender:(*types.JID)(0xc00037cf30), SenderPN:(*types.JID)(0xc00037cf60), Timestamp:time.Date(2025, time.October, 11, 8, 20, 1, 0, time.Local), Name:(*types.GroupName)(nil), Topic:(*types.GroupTopic)(nil), Locked:(*types.GroupLocked)(nil), Announce:(*types.GroupAnnounce)(nil), Ephemeral:(*types.GroupEphemeral)(nil), MembershipApprovalMode:(*types.GroupMembershipApprovalMode)(nil), Delete:(*types.GroupDelete)(nil), Link:(*types.GroupLinkChange)(nil), Unlink:(*types.GroupLinkChange)(nil), NewInviteLink:(*string)(nil), PrevParticipantVersionID:"1760187609812467", ParticipantVersionID:"1760188801051036", JoinReason:"", Join:[]types.JID{types.JID{User:"21926445588488", RawAgent:0x0, Device:0x0, Integrator:0x0, Server:"lid"}}, Leave:[]types.JID(nil), Promote:[]types.JID(nil), Demote:[]types.JID(nil), UnknownChanges:[]*binary.Node(nil)}{Text:joined chat Channel:120363422464172527@g.us Username:mentalisit UserID:21926445588488@lid Avatar: Account: Event:join_leave Protocol: Gateway: ParentID: Timestamp:0001-01-01 00:00:00 +0000 UTC ID: Extra:map[]}

*/
