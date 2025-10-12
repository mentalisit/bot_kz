package telegram

import (
	"fmt"
	"strconv"
	"strings"
	"telegram/models"
	"telegram/telegram/helper"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

func (t *Telegram) handleForwarded(rmsg *models.ToBridgeMessage, message *tgbotapi.Message) {
	if message.ForwardOrigin == nil {
		return
	}

	if message.ForwardOrigin.Chat != nil {
		rmsg.Text = "Forwarded from " + message.ForwardOrigin.Chat.Title + ": " + rmsg.Text
		return
	}

	usernameForward := message.ForwardOrigin.SenderUser.String()

	if usernameForward == "" {
		usernameForward = "unknown"
	}

	rmsg.Text = "Forwarded from " + usernameForward + ": " + rmsg.Text
}

func (t *Telegram) handleQuoting(rmsg *models.ToBridgeMessage, message *tgbotapi.Message) {
	if message.ReplyToMessage != nil && (!message.IsTopicMessage || message.ReplyToMessage.MessageID != message.MessageThreadID) {
		usernameReply := ""
		if message.ReplyToMessage.From != nil {
			f := message.ReplyToMessage.From
			if f.UserName != "" {
				usernameReply = f.UserName
			} else {
				usernameReply = f.FirstName + " " + f.LastName
			}
		}
		if usernameReply == "" {
			usernameReply = "unknown"
		}

		quote := message.ReplyToMessage.Text
		if quote == "" {
			quote = message.ReplyToMessage.Caption
		}
		fmt.Println(usernameReply)
		fmt.Println(quote)

		reply := &models.BridgeMessageReply{
			TimeMessage: message.ReplyToMessage.Time().Unix(),
			Text:        quote,
			Avatar:      t.getAvatarIsExist(message.ReplyToMessage.From.ID),
			UserName:    usernameReply,
		}
		rmsg.ReplyMap = make(map[string]string)
		rmsg.ReplyMap[rmsg.ChatId] = strconv.Itoa(message.ReplyToMessage.MessageID)
		rmsg.Reply = reply
		fmt.Printf("%+v\n", rmsg.Reply)
	}
}

func (t *Telegram) handleDownloadBridge(rmsg *models.ToBridgeMessage, message *tgbotapi.Message) error {
	if message.Text == "" && message.Caption != "" {
		rmsg.Text = message.Caption
	}
	size := int64(0)
	var url, name, FileID string
	switch {
	case message.Sticker != nil:
		FileID = message.Sticker.FileID
		name, url = t.getDownloadInfo(FileID, ".webp", true)
		size = int64(message.Sticker.FileSize)
	case message.Voice != nil:
		FileID = message.Voice.FileID
		name, url = t.getDownloadInfo(FileID, ".ogg", true)
		size = message.Voice.FileSize
	case message.Video != nil:
		FileID = message.Video.FileID
		name, url = t.getDownloadInfo(FileID, "", true)
		size = message.Video.FileSize
	case message.Audio != nil:
		FileID = message.Audio.FileID
		name, url = t.getDownloadInfo(FileID, "", true)
		size = message.Audio.FileSize
	case message.Document != nil:
		FileID = message.Document.FileID
		_, url = t.getDownloadInfo(FileID, "", false)
		size = message.Document.FileSize
		name = message.Document.FileName
	case message.Photo != nil:
		photos := message.Photo
		size = int64(photos[len(photos)-1].FileSize)
		FileID = photos[len(photos)-1].FileID
		name, url = t.getDownloadInfo(FileID, "", true)
	}

	// if name is empty we didn't match a thing to download
	if name == "" {
		return nil
	}
	if int(size) > 10*1024*1024 {
		return nil
	}
	// if we have a file attached, download it (in memory) and put a pointer to it in msg.Extra
	data, err := helper.DownloadFile(url)
	if err != nil {
		return err
	}

	if strings.HasSuffix(name, ".tgs.webp") {
		//b.maybeConvertTgs(&name, data)
	} else if strings.HasSuffix(name, ".webp") {
		t.maybeConvertWebp(&name, &data)
	}

	if strings.HasSuffix(name, ".oga") && message.Audio != nil {
		name = strings.Replace(name, ".oga", ".ogg", 1)
	}

	rmsg.Extra = append(rmsg.Extra, models.FileInfo{
		Name:   name,
		Data:   data,
		URL:    url,
		Size:   size,
		FileID: FileID,
	})
	return nil
}

func (t *Telegram) getDownloadInfo(id string, suffix string, urlpart bool) (string, string) {
	url := t.getFileDirectURL(id)
	name := ""
	if urlpart {
		urlPart := strings.Split(url, "/")
		name = urlPart[len(urlPart)-1]
	}
	if suffix != "" && !strings.HasSuffix(name, suffix) && !strings.HasSuffix(name, ".webm") {
		name += suffix
	}
	return name, url
}
func (t *Telegram) getFileDirectURL(id string) string {
	res, err := t.t.GetFileDirectURL(id)
	if err != nil {
		return ""
	}
	return res
}

//	func (b *Telegram) maybeConvertTgs(name *string, data *[]byte) {
//		format := b.GetString("MediaConvertTgs")
//		if helper.SupportsFormat(format) {
//			b.Log.Debugf("Format supported by %s, converting %v", helper.LottieBackend(), name)
//		} else {
//			// Otherwise, no conversion was requested. Trying to run the usual webp
//			// converter would fail, because '.tgs.webp' is actually a gzipped JSON
//			// file, and has nothing to do with WebP.
//			return
//		}
//		err := helper.ConvertTgsToX(data, format, b.Log)
//		if err != nil {
//			b.Log.Errorf("conversion failed: %v", err)
//		} else {
//			*name = strings.Replace(*name, "tgs.webp", format, 1)
//		}
//	}
func (t *Telegram) maybeConvertWebp(name *string, data *[]byte) {
	err := helper.ConvertWebPToPNG(data)
	if err != nil {
		t.log.ErrorErr(err)
	} else {
		*name = strings.Replace(*name, ".webp", ".png", 1)
	}
}
