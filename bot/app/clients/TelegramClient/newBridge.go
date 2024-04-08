package TelegramClient

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"kz_bot/clients/restapi"
	"kz_bot/models"
	"strconv"
)

func (t *Telegram) filterNewBridge(m *tgbotapi.Message, mes models.ToBridgeMessage) {
	mes.Text = m.Text
	mes.Tip = "tg"
	mes.MesId = strconv.Itoa(m.MessageID)
	mes.GuildId = strconv.FormatInt(m.Chat.ID, 10)
	mes.TimestampUnix = m.Time().Unix()

	if m.From.UserName != "" {
		mes.Sender = m.From.UserName
	} else {
		mes.Sender = m.From.FirstName + " " + m.From.LastName
	}
	mes.Avatar = t.getAvatarIsExist(m.From.ID)

	err := t.handleDownloadBridge(&mes, m)
	if err != nil {
		t.log.ErrorErr(err)
	}

	// handle forwarded messages
	t.handleForwarded(&mes, m)

	// quote the previous message
	t.handleQuoting(&mes, m)

	if mes.Text != "" || len(mes.Extra) > 0 {
		err = restapi.SendBridgeApp(mes)
		if err != nil {
			t.log.ErrorErr(err)
		}
	}
}
