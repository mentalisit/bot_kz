package wa

import (
	"context"
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"whatsapp/models"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	goproto "google.golang.org/protobuf/proto"
)

func (b *Whatsapp) Send(msg Message) (string, error) {
	extendedMsgID, _ := b.parseMessageID(msg.ID)
	msg.ID = extendedMsgID.MessageID

	var mesIds []string
	var lastError error

	// 1. Обработка отправки файлов
	if len(msg.Extra) > 0 {
		// Итерируемся по каждому файлу в Extra
		for _, fi := range msg.Extra {

			// Проверяем, есть ли данные для отправки
			if len(fi.Data) == 0 {
				// Здесь должна быть логика скачивания по URL, если она не реализована выше
				fmt.Printf("File data is empty for %s, skipping or needs download.\n", fi.Name)
				continue
			}

			// !!! КЛЮЧЕВОЕ ИСПРАВЛЕНИЕ: Создаем временную структуру msg,
			// которая содержит ТОЛЬКО текущий файл.
			tempMsg := msg
			tempMsg.Extra = []models.FileInfo{fi}

			// Определяем MIME-тип
			filetype := mime.TypeByExtension(filepath.Ext(fi.Name))
			var ID string

			if filetype == "" {
				filetype = http.DetectContentType(fi.Data)
			}

			switch filetype {
			case "image/jpeg", "image/png", "image/gif":
				fmt.Println("filetype ", filetype)
				ID, lastError = b.PostImageMessage(tempMsg, filetype)
			case "video/mp4", "video/3gpp":
				ID, lastError = b.PostVideoMessage(tempMsg, filetype)
			case "audio/ogg", "application/ogg":
				// ИСПРАВЛЕНИЕ 1: Для аудио лучше использовать универсальную функцию,
				// чтобы избежать дублирования логики.
				// В вашем текущем коде:
				ID, lastError = b.PostAudioMessage(tempMsg, "audio/ogg; codecs=opus")
			case "audio/aac", "audio/mp4", "audio/amr", "audio/mpeg":
				ID, lastError = b.PostAudioMessage(tempMsg, filetype)
			default:
				ID, lastError = b.PostDocumentMessage(tempMsg, filetype)
			}

			// Обработка ошибок
			if lastError != nil {
				// Логируем ошибку, но продолжаем, чтобы попробовать отправить другие файлы
				b.log.ErrorErr(fmt.Errorf("error sending file %s: %v", fi.Name, lastError))
				continue
			}

			// Сохраняем ID успешно отправленного сообщения
			if ID != "" {
				mesIds = append(mesIds, ID)
			}
		}
	}

	// 2. Обработка возвращаемого ID
	// Если было отправлено медиа, возвращаем ID первого сообщения.
	if len(mesIds) > 0 {
		return mesIds[0], nil
	}

	// 3. Отправка текстового сообщения (если медиа не было или не было отправлено)
	// ... (Остальная логика отправки текста остается без изменений)
	var message waE2E.Message
	text := msg.Username + msg.Text

	// Логика ответа (ParentID)
	if msg.ParentID != "" {
		// ... (логика ExtendedTextMessage)
		replyContext, err := b.getNewReplyContext(msg.ParentID)
		if err == nil {
			message = waE2E.Message{
				ExtendedTextMessage: &waE2E.ExtendedTextMessage{
					Text:        &text,
					ContextInfo: replyContext,
				},
			}
			return b.sendMessage(msg, &message)
		}
	}

	// Логика обычного текста
	message.Conversation = &text
	return b.sendMessage(msg, &message)
}

func (b *Whatsapp) sendMessage(rmsg Message, message *waE2E.Message) (string, error) {
	groupJID, _ := b.ParseJID(rmsg.Channel)
	ID := b.wc.GenerateMessageID()

	_, err := b.wc.SendMessage(context.Background(), groupJID, message, whatsmeow.SendRequestExtra{ID: ID})

	return getMessageIdFormat(*b.wc.Store.ID, ID), err
}

// Post a document message from the bridge to WhatsApp
func (b *Whatsapp) PostDocumentMessage(msg Message, filetype string) (string, error) {
	groupJID, _ := b.ParseJID(msg.Channel)

	fi := msg.Extra[0]

	caption := msg.Username

	resp, err := b.wc.Upload(context.Background(), fi.Data, whatsmeow.MediaDocument)
	if err != nil {
		return "", err
	}

	// Post document message
	var message waE2E.Message
	var ctx *waE2E.ContextInfo
	if msg.ParentID != "" {
		ctx, _ = b.getNewReplyContext(msg.ParentID)
	}

	message.DocumentMessage = &waE2E.DocumentMessage{
		Title:         &fi.Name,
		FileName:      &fi.Name,
		Mimetype:      &filetype,
		Caption:       &caption,
		MediaKey:      resp.MediaKey,
		FileEncSHA256: resp.FileEncSHA256,
		FileSHA256:    resp.FileSHA256,
		FileLength:    goproto.Uint64(resp.FileLength),
		URL:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		ContextInfo:   ctx,
	}

	ID := b.wc.GenerateMessageID()
	_, err = b.wc.SendMessage(context.TODO(), groupJID, &message, whatsmeow.SendRequestExtra{ID: ID})

	return ID, err
}

// Post an image message from the bridge to WhatsApp
// Handle, for sure image/jpeg, image/png and image/gif MIME types
func (b *Whatsapp) PostImageMessage(msg Message, filetype string) (string, error) {
	fi := msg.Extra[0]

	caption := msg.Text
	fmt.Println("msg.Text ", msg.Text)

	resp, err := b.wc.Upload(context.Background(), fi.Data, whatsmeow.MediaImage)
	if err != nil {
		return "", err
	}

	var message waE2E.Message
	var ctx *waE2E.ContextInfo
	if msg.ParentID != "" {
		ctx, _ = b.getNewReplyContext(msg.ParentID)
	}

	message.ImageMessage = &waE2E.ImageMessage{
		Mimetype:      &filetype,
		Caption:       &caption,
		MediaKey:      resp.MediaKey,
		FileEncSHA256: resp.FileEncSHA256,
		FileSHA256:    resp.FileSHA256,
		FileLength:    goproto.Uint64(resp.FileLength),
		URL:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		ContextInfo:   ctx,
	}

	return b.sendMessage(msg, &message)
}

// Post a video message from the bridge to WhatsApp
func (b *Whatsapp) PostVideoMessage(msg Message, filetype string) (string, error) {
	fi := msg.Extra[0]

	caption := msg.Username

	resp, err := b.wc.Upload(context.Background(), fi.Data, whatsmeow.MediaVideo)
	if err != nil {
		return "", err
	}

	var message waE2E.Message
	var ctx *waE2E.ContextInfo
	if msg.ParentID != "" {
		ctx, _ = b.getNewReplyContext(msg.ParentID)
	}

	message.VideoMessage = &waE2E.VideoMessage{
		Mimetype:      &filetype,
		Caption:       &caption,
		MediaKey:      resp.MediaKey,
		FileEncSHA256: resp.FileEncSHA256,
		FileSHA256:    resp.FileSHA256,
		FileLength:    goproto.Uint64(resp.FileLength),
		URL:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		ContextInfo:   ctx,
	}

	return b.sendMessage(msg, &message)
}

// Post audio inline
func (b *Whatsapp) PostAudioMessage(msg Message, filetype string) (string, error) {
	//groupJID, _ := b.ParseJIDmsg.Channel)

	fi := msg.Extra[0]

	resp, err := b.wc.Upload(context.Background(), fi.Data, whatsmeow.MediaAudio)
	if err != nil {
		return "", err
	}

	var message waE2E.Message
	var ctx *waE2E.ContextInfo
	if msg.ParentID != "" {
		ctx, _ = b.getNewReplyContext(msg.ParentID)
	}

	message.AudioMessage = &waE2E.AudioMessage{
		Mimetype:      &filetype,
		MediaKey:      resp.MediaKey,
		FileEncSHA256: resp.FileEncSHA256,
		FileSHA256:    resp.FileSHA256,
		FileLength:    goproto.Uint64(resp.FileLength),
		URL:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		ContextInfo:   ctx,
		PTT:           goproto.Bool(true)}

	ID, err := b.sendMessage(msg, &message)
	if err != nil {
		fmt.Println(err)
	}

	var captionMessage waE2E.Message
	text := msg.Username + msg.Text + "\u2B06" // the char on the end is upwards arrow emoji

	captionMessage.Conversation = &text
	_, err = b.sendMessage(msg, &captionMessage)
	if err != nil {
		fmt.Println(err)
	}

	return ID, err
}
