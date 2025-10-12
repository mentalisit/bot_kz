package wa

import (
	"context"
	"time"
	"whatsapp/models"

	"go.mau.fi/util/ptr"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waCommon"
	"go.mau.fi/whatsmeow/proto/waE2E"
	goproto "google.golang.org/protobuf/proto"
)

func (b *Whatsapp) EditText(ChatID, mId, TextEdit string) error {
	groupJID, _ := b.ParseJID(ChatID)
	// 1. Создание ключа сообщения
	messageKey := waCommon.MessageKey{
		RemoteJID: ptr.Ptr(groupJID.String()),
		FromMe:    ptr.Ptr(true),
		ID:        ptr.Ptr(mId),
	}
	mes := waE2E.Message{Conversation: ptr.Ptr(TextEdit)}

	// 2. Создание ProtocolMessage
	protocolMessage := &waE2E.ProtocolMessage{
		Type:          ptr.Ptr(waE2E.ProtocolMessage_MESSAGE_EDIT),
		Key:           &messageKey, // В новых версиях - EditMessageKey
		EditedMessage: &mes,
	}

	// 3. Отправка сообщения
	content := &waE2E.Message{Conversation: &TextEdit,
		ProtocolMessage: protocolMessage,
	}

	ID := b.wc.GenerateMessageID()

	_, err := b.wc.SendMessage(context.Background(), groupJID, content, whatsmeow.SendRequestExtra{ID: ID})
	if err != nil {
		return err
	}

	return nil
}

func (b *Whatsapp) SendChannel(ChatID, text, parseMode string) (id string, err error) {
	var message waE2E.Message
	if parseMode == "" {
		message.Conversation = &text
	} else {
		message.ExtendedTextMessage = &waE2E.ExtendedTextMessage{Text: &text}
	}

	groupJID, _ := b.ParseJID(ChatID)
	ID := b.wc.GenerateMessageID()

	_, err = b.wc.SendMessage(context.Background(), groupJID, &message, whatsmeow.SendRequestExtra{ID: ID})

	return getMessageIdFormat(*b.wc.Store.ID, ID), err

}
func (b *Whatsapp) SendChannelDelSecond(chatid string, text string, second int) (bool, error) {
	sendMessage, err := b.SendChannel(chatid, text, "")
	if err != nil {
		return false, err
	}
	tu := int(time.Now().UTC().Unix())
	b.Storage.Db.TimerInsert(models.Timer{
		Tip:    "wa",
		ChatId: chatid,
		MesId:  sendMessage,
		Timed:  tu + second,
	})

	if sendMessage != "" {
		return true, nil
	}
	return false, nil
}

func (b *Whatsapp) SendPic(Chatid, Text string, ImageBytes []byte) error {
	resp, err := b.wc.Upload(context.Background(), ImageBytes, whatsmeow.MediaImage)
	if err != nil {
		return err
	}

	var message waE2E.Message
	var ctx *waE2E.ContextInfo
	mimetype := "image/jpeg"

	message.ImageMessage = &waE2E.ImageMessage{
		Mimetype:      &mimetype,
		Caption:       &Text,
		MediaKey:      resp.MediaKey,
		FileEncSHA256: resp.FileEncSHA256,
		FileSHA256:    resp.FileSHA256,
		FileLength:    goproto.Uint64(resp.FileLength),
		URL:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		ContextInfo:   ctx,
	}

	groupJID, _ := b.ParseJID(Chatid)
	ID := b.wc.GenerateMessageID()

	_, err = b.wc.SendMessage(context.Background(), groupJID, &message, whatsmeow.SendRequestExtra{ID: ID})
	if err != nil {
		return err
	}
	return nil

}
func (b *Whatsapp) GetAvatarUrl(Userid string) string {
	if avatarURL, exists := b.userAvatars[Userid]; exists {
		_, url := b.SaveAvatarLocalCache(Userid, avatarURL)
		return url
	}
	return ""
}
