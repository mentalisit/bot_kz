package TelegramClient

import (
	"fmt"
	"kz_bot/models"
	"strconv"
)

func (t *Telegram) Help(Channel string, lang string) {
	text := fmt.Sprintf("%s\n%s ", t.getLanguage(lang, "information"), t.getLanguage(lang, "info_help_text"))
	mid := t.SendHelp(Channel, text)
	t.DelMessageSecond(Channel, strconv.Itoa(mid), 180)
}
func (t *Telegram) Help1(Channel string, lang string) string {
	text := fmt.Sprintf("%s\n%s ", t.getLanguage(lang, "information"), t.getLanguage(lang, "info_help_text"))
	mid := t.SendHelp(Channel, text)
	return strconv.Itoa(mid)
}

// команда хелп
func (t *Telegram) help(chatid string, mesid string) {
	t.DelMessageSecond(chatid, mesid, 10)
	t.SendChannelDelSecond(chatid, models.Help, 60)
}

// очередь кз
func (t *Telegram) helpQueue(chatid string, mesid string) {
	go t.DelMessageSecond(chatid, mesid, 10)
	t.SendChannelDelSecond(chatid, models.HelpQueue, 60)
}

// Уведомления
func (t *Telegram) helpNotification(chatid string, mesid string) {
	go t.DelMessageSecond(chatid, mesid, 10)
	t.SendChannelDelSecond(chatid, models.HelpNotification, 60)
}

// Ивент кз
func (t *Telegram) helpEvent(chatid string, mesid string) {
	go t.DelMessageSecond(chatid, mesid, 10)
	t.SendChannelDelSecond(chatid, models.HelpEvent, 60)
}

// Топ лист
func (t *Telegram) helpTop(chatid string, mesid string) {
	go t.DelMessageSecond(chatid, mesid, 10)

	t.SendChannelDelSecond(chatid, models.HelpTop, 60)
}

// Работа с иконками
func (t *Telegram) helpIcon(chatid string, mesid string) {
	go t.DelMessageSecond(chatid, mesid, 10)

	t.SendChannelDelSecond(chatid, models.HelpIcon, 60)
}
