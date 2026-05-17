package bot2

import (
	"fmt"
	"net/http"
	"rs/models"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (b *Bot) accessChat(in *models.InMessageV2) {
	prefix := strings.HasPrefix(in.Text, ".")
	if prefix {
		if strings.HasPrefix(in.Text, ".setup") {
			b.accessRSBot(in)
			return
		} else if strings.HasPrefix(in.Text, ".invite ") {
			b.accessInviteRSBot(in)
			return
		} else if strings.HasPrefix(in.Text, ".link") {
			b.accessLinkCode(in)
			return
		} else if strings.HasPrefix(in.Text, ".setting") {
			b.accessSetting(in)
			return
		}
	}
}

func (b *Bot) accessSetting(in *models.InMessageV2) {
	if in.MAcc == nil {
		b.log.Info("MAcc not found ")
		return
	}

	// Проверяем, является ли это командой формата .setting ws
	if strings.HasPrefix(in.Text, ".setting ws ") {
		b.handleWsSetting(in)
		return
	}

	// Формируем URL для мобильной веб-страницы
	serverURL := "https://mentalisit.myds.me"

	// Проверяем доступность сайта перед отправкой
	if !b.CheckSiteAvailability(serverURL) {
		serverURL = "https://mentalisit.tsl.rocks"
	}

	setupURL := fmt.Sprintf("%s/rs/settings/settings.html?uuid=%s", serverURL, in.MAcc.UUID.String())

	// Отправляем сообщение со ссылкой
	message := fmt.Sprintf("🔧 **Настройка бота через веб-интерфейс**\n\n"+
		"Для настройки своего профиля перейдите по временной ссылке:\n%s", setupURL)
	text := "Ссылка отправлена в ЛС"
	if in.Tip == tg {
		text = fmt.Sprintf("%s\n Если вы не получили сообщение сначала нажмите на бота и выполните Старт", text)
	}
	b.sendTypeMessenger(in.Tip, in.Messenger.ChannelId, text, 30)

	if in.Tip == "ds" {
		b.client.Ds.SendDmText(message, in.UserId)
	} else if in.Tip == "tg" {
		b.client.Tg.SendChannelId(in.UserId, message)
	} else {
		b.log.Info("not implement ")
	}
}

// accessLinkCode генерирует 6-значный код для привязки аккаунта на сайте
func (b *Bot) accessLinkCode(in *models.InMessageV2) {
	fmt.Println("accessLinkCode")
	b.deleteInMessage(in)

	// Генерируем 6-значный числовой код
	code := fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)

	if b.AddLinkCodeFunc != nil {
		// Передаем данные в веб-сервер
		b.AddLinkCodeFunc(code, in.UserId, in.Username, in.Tip)

		// Отправляем сообщение пользователю
		msg := fmt.Sprintf("🔐 **Код для привязки аккаунта**\n\n"+
			"Ваш код подтверждения: `%s`\n"+
			"Введите его на странице настроек в течение 10 минут.", code)

		// Отправляем личным сообщением или в чат (в зависимости от настроек)
		b.sendTextAfterDeleteSecond(in, msg, 600) // Сообщение удалится через 10 минут
	} else {
		b.sendTextAfterDeleteSecond(in, "Ошибка: Веб-сервер не инициализирован", 30)
	}
}

// accessRSBot отправляет пользователю ссылку на мобильную веб-страницу настроек
func (b *Bot) accessRSBot(in *models.InMessageV2) {
	fmt.Println("accessRSBot")
	b.deleteInMessage(in)
	p := models.Other{
		//Uuid:     uuid.New().String(),
		DataType: "setup",
		Data:     in.Messenger,
	}

	//rsBot
	if p.Uuid == "" {
		exist, conf := b.checkConfig(in)
		if exist {
			p.Uuid = conf.Uid
		}
	}

	//bridge
	if p.Uuid == "" {
		bridge2Config, existBridge := b.storage.ReadBridgeConfigByChannelId(in.Messenger.ChannelId)
		if existBridge {
			ok, u := checkUUID(bridge2Config.NameRelay)
			if !ok {
				_ = b.storage.UpdateBridgeConfigNameRelay(bridge2Config.NameRelay, u.String())
			}
			p.Uuid = u.String()
		}
	}

	//news
	if p.Uuid == "" {
		configNews, exist := b.storage.IsSubscribedToNews(in.Messenger.ChannelId)
		if exist {
			p.Uuid = configNews.Uid
		}
	}

	//scoreboard
	if p.Uuid == "" {
		scoreboardParamsV2 := b.storage.ScoreboardReadByChannelId(in.Messenger.ChannelId)
		if scoreboardParamsV2 != nil && len(scoreboardParamsV2.Channels) != 0 {
			p.Uuid = scoreboardParamsV2.Uid
		}
	}

	//if new
	if p.Uuid == "" {
		p.Uuid = uuid.New().String()
	}

	// Формируем URL для мобильной веб-страницы
	serverURL := "https://mentalisit.myds.me"

	// Проверяем доступность сайта перед отправкой
	if !b.CheckSiteAvailability(serverURL) {
		serverURL = "https://mentalisit.tsl.rocks"
	}

	setupURL := fmt.Sprintf("%s/rs/api/setup/config.html?uuid=%s", serverURL, p.Uuid)

	// Отправляем сообщение со ссылкой
	message := fmt.Sprintf("🔧 **Настройка бота через веб-интерфейс**\n\n"+
		"Для управления всеми параметрами бота перейдите по временной ссылке:\n%s", setupURL)
	p.Data.MessageId = b.SendTextReturnId(in, message)
	//b.sendTextAfterDeleteSecond(in, message, 30)
	b.storage.InsertOther(p)
}

func (b *Bot) accessInviteRSBot(in *models.InMessageV2) {
	fmt.Println("accessInviteRSBot")
	b.deleteInMessage(in)

	after, found := strings.CutPrefix(in.Text, ".invite rs ")
	if found {
		uid, err := uuid.Parse(after)
		if err != nil {
			b.log.ErrorErr(err)
			return
		}
		conf := b.storage.ReadConfigV2Uid(uid.String())
		if conf != nil && conf.Uid != "" {
			if conf.Channels[in.Messenger.ChannelId] != nil {
				fmt.Printf("conf.Channels[in.Messenger.ChannelId]!=nil\n")
				b.sendTextAfterDeleteSecond(in, "already added ", 30)
				return
			}
			conf.Channels[in.Messenger.ChannelId] = &in.Messenger
			if conf.Channels[in.Messenger.ChannelId].Game == nil {
				conf.Channels[in.Messenger.ChannelId].Game = &models.GameSettings{}
			}
			if conf.Channels[in.Messenger.ChannelId].Corp == nil {
				conf.Channels[in.Messenger.ChannelId].Corp = &models.CorpSettings{}
			}

			b.storage.UpdateConfigV2Channels(*conf)
			b.sendTextAfterDeleteSecond(in, "добавлено", 30)
			in.Config = *conf
			confNew := b.SendHelpInMessenger(in)
			b.storage.UpdateConfigV2HelpMessage(confNew)
		}
	}

	after, found = strings.CutPrefix(in.Text, ".invite bridge ")
	if found {
		uid, err := uuid.Parse(after)
		if err != nil {
			b.log.ErrorErr(err)
			return
		}
		br := models.Bridge2Configs{
			ChannelId:       in.Messenger.ChannelId,
			GuildId:         in.Messenger.GuildId,
			CorpChannelName: fmt.Sprintf("%s - %s", in.Messenger.GuildName, in.Messenger.ChannelName),
			AliasName:       "",
			MappingRoles:    make(map[string]string),
		}
		inType := in.Messenger.TypeMessenger
		relay, exist := b.storage.ReadBridgeConfigByNameRelay(uid.String())
		if exist {
			if relay.Channel[inType] == nil {
				relay.Channel[inType] = []models.Bridge2Configs{}
			}
			relay.Channel[inType] = append(relay.Channel[inType], br)
		} else {
			conf := b.storage.ReadConfigV2Uid(uid.String())
			if conf != nil && conf.Uid != "" {
				return
			}
			other, errOther := b.storage.GetOtherByUUID(uid)
			if errOther != nil || other == nil {
				return
			}
			newBr := models.Bridge2Config{
				NameRelay:         other.Uuid,
				HostRelay:         fmt.Sprintf("%s - %s", other.Data.GuildName, other.Data.ChannelName),
				Role:              []string{},
				Channel:           make(map[string][]models.Bridge2Configs),
				ForbiddenPrefixes: []string{},
			}
			oldBr := models.Bridge2Configs{
				ChannelId:       other.Data.ChannelId,
				GuildId:         other.Data.GuildId,
				CorpChannelName: newBr.HostRelay,
				AliasName:       "",
				MappingRoles:    make(map[string]string),
			}
			if newBr.Channel[other.Data.TypeMessenger] == nil {
				newBr.Channel[other.Data.TypeMessenger] = []models.Bridge2Configs{}
			}
			newBr.Channel[other.Data.TypeMessenger] = append(newBr.Channel[other.Data.TypeMessenger], oldBr)
			newBr.Channel[inType] = append(newBr.Channel[inType], br)

			err = b.storage.InsertBridgeConfig(newBr)
			if err != nil {
				b.log.ErrorErr(err)
				return
			}
			b.sendTextAfterDeleteSecond(in, "добавлено", 30)

		}
	}

	after, found = strings.CutPrefix(in.Text, ".invite scoreboard ")
	if found {
		uid, err := uuid.Parse(after)
		if err != nil {
			b.log.ErrorErr(err)
			return
		}
		params := b.storage.ScoreboardReadByUid(uid.String())
		ch := models.ChannelsInfo{
			GuildChannelName: fmt.Sprintf("%s - %s", in.Messenger.GuildName, in.Messenger.ChannelName),
			TypeMessenger:    in.Messenger.TypeMessenger,
			ChannelId:        in.Messenger.ChannelId,
		}
		text := ""
		if params != nil {
			params.Channels = append(params.Channels, ch)
			b.storage.ScoreboardUpdateParamChannels(*params)
			text = "Табло лидеров ивента подключено"
		} else {
			text = "конфиг не найден"
		}
		b.sendTextAfterDeleteSecond(in, text, 60)

	}
}
func checkUUID(uid string) (ok bool, u uuid.UUID) {
	parse, err := uuid.Parse(uid)
	if err != nil {
		return false, uuid.New()
	}
	return true, parse

}

// CheckSiteAvailability проверяет доступность сайта по URL
func (b *Bot) CheckSiteAvailability(url string) bool {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// handleWsSetting обрабатывает команды формата .setting ws {gid} {action}
func (b *Bot) handleWsSetting(in *models.InMessageV2) {
	b.deleteInMessage(in)

	// Разбираем команду: .setting ws {gid} {action} [params...]
	parts := strings.Fields(in.Text)
	if len(parts) < 4 {
		b.sendTypeMessenger(in.Tip, in.Messenger.ChannelId,
			"❌ **Неверный формат команды**\n\nИспользуйте: `.setting ws {gid} {action}`\n\nПримеры:\n- `.setting ws af3ebcc4-77a1-4e9f-9aa9-6f2b4686c6f5 poll`\n- `.setting ws af3ebcc4-77a1-4e9f-9aa9-6f2b4686c6f5 бз1 coordination`",
			30)
		return
	}

	gid := parts[2]
	action := parts[3]

	// Проверяем, что gid похож на UUID (базовая проверка)
	if len(gid) < 20 {
		b.sendTypeMessenger(in.Tip, in.Messenger.ChannelId,
			"❌ **Неверный формат GID**\n\nGID должен быть в формате UUID",
			30)
		return
	}

	// Обрабатываем различные действия
	switch action {
	case "poll":
		b.handlePollChannel(in, gid)
	case "coordination", "discussion":
		if len(parts) < 5 {
			b.sendTypeMessenger(in.Tip, in.Messenger.ChannelId,
				"❌ **Неверный формат команды**\n\nДля БЗ каналов используйте: `.setting ws {gid} {bz_key} {channel_type}`\n\nПример: `.setting ws af3ebcc4-77a1-4e9f-9aa9-6f2b4686c6f5 бз1 coordination`",
				30)
			return
		}
		bzKey := parts[4]
		b.handleBzChannel(in, gid, bzKey, action)
	default:
		b.sendTypeMessenger(in.Tip, in.Messenger.ChannelId,
			fmt.Sprintf("❌ **Неизвестное действие**\n\nДействие `%s` не поддерживается\n\nПоддерживаемые действия:\n- `poll`\n- `coordination`\n- `discussion`", action),
			30)
	}
}

// handlePollChannel обрабатывает добавление канала опросов
func (b *Bot) handlePollChannel(in *models.InMessageV2, gid string) {
	// Добавляем канал в базу данных
	channelFullName := fmt.Sprintf("%s.%s", in.Messenger.GuildName, in.Messenger.ChannelName)
	err := b.storage.SavePollChannel(gid, in.Messenger.ChannelId, channelFullName, in.Tip)
	if err != nil {
		message := fmt.Sprintf("❌ **Ошибка при добавлении канала опросов**\n\nОшибка: %s", err.Error())
		b.sendTypeMessenger(in.Tip, in.Messenger.ChannelId, message, 30)
		return
	}

	message := fmt.Sprintf("✅ **Канал опросов успешно добавлен**\n\n"+
		"GID: `%s`\n"+
		"Канал: %s\n"+
		"Мессенджер: %s",
		gid, in.Messenger.ChannelId, in.Tip)

	b.sendTypeMessenger(in.Tip, in.Messenger.ChannelId, message, 30)
}

// handleBzChannel обрабатывает добавление канала БЗ
func (b *Bot) handleBzChannel(in *models.InMessageV2, gid, bzKey, channelType string) {
	channelFullName := fmt.Sprintf("%s.%s", in.Messenger.GuildName, in.Messenger.ChannelName)
	// Добавляем канал БЗ в базу данных
	err := b.storage.SaveBzChannel(gid, bzKey, channelType, in.Messenger.ChannelId, channelFullName, in.Tip)
	if err != nil {
		message := fmt.Sprintf("❌ **Ошибка при добавлении канала БЗ**\n\nОшибка: %s", err.Error())
		b.sendTypeMessenger(in.Tip, in.Messenger.ChannelId, message, 30)
		return
	}

	message := fmt.Sprintf("✅ **Канал БЗ успешно добавлен**\n\n"+
		"GID: `%s`\n"+
		"БЗ: `%s`\n"+
		"Тип канала: `%s`\n"+
		"Канал: %s\n"+
		"Мессенджер: %s",
		gid, bzKey, channelType, in.Messenger.ChannelId, in.Tip)

	b.sendTypeMessenger(in.Tip, in.Messenger.ChannelId, message, 30)
}
