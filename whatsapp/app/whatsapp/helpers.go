package wa

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"go.mau.fi/whatsmeow/proto/waE2E"
	goproto "google.golang.org/protobuf/proto"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type ProfilePicInfo struct {
	URL    string `json:"eurl"`
	Tag    string `json:"tag"`
	Status int16  `json:"status"`
}

func (b *Whatsapp) reloadContacts() {
	// Fix: Add context.Background() as first parameter
	if _, err := b.wc.Store.Contacts.GetAllContacts(context.Background()); err != nil {
		b.log.ErrorErr(err)
	}

	// Fix: Add context.Background() as first parameter
	allcontacts, err := b.wc.Store.Contacts.GetAllContacts(context.Background())
	if err != nil {
		b.log.ErrorErr(err)
	}

	if len(allcontacts) > 0 {
		b.contacts = allcontacts
	}
}

func (b *Whatsapp) getSenderName(info types.MessageInfo) string {
	// Parse AD JID
	var senderJid types.JID
	senderJid.User, senderJid.Server = info.Sender.User, info.Sender.Server

	sender, exists := b.contacts[senderJid]

	if !exists || (sender.FullName == "" && sender.FirstName == "") {
		b.GetGroupMemberNames(info.Chat)
		b.reloadContacts() // Contacts may need to be reloaded
		//b.GetGroupMemberNames(info.Chat)
		sender, exists = b.contacts[senderJid]
	}

	if exists && sender.FullName != "" {
		return sender.FullName
	}

	if info.PushName != "" {
		return info.PushName
	}

	if exists && sender.FirstName != "" {
		return sender.FirstName
	}

	return "Someone"
}

func (b *Whatsapp) getSenderNameFromJID(senderJid types.JID) string {
	sender, exists := b.contacts[senderJid]

	if !exists || (sender.FullName == "" && sender.FirstName == "") {
		b.reloadContacts() // Contacts may need to be reloaded
		sender, exists = b.contacts[senderJid]
	}

	if exists && sender.FullName != "" {
		return sender.FullName
	}

	if exists && sender.FirstName != "" {
		return sender.FirstName
	}

	if sender.PushName != "" {
		return sender.PushName
	}

	return "Someone"
}

// groupJID - JID группы (например, types.NewJID("1234567890", types.GroupServer))
func (b *Whatsapp) GetGroupMemberNames(groupJID types.JID) (defaultName string) {
	defaultName = "Someone"
	// Запрашиваем информацию о группе у сервера WhatsApp
	groupInfo, err := b.wc.GetGroupInfo(context.Background(), groupJID)
	if err != nil {
		fmt.Printf("Ошибка при получении информации о группе: %v\n", err)
		return
	}

	//fmt.Printf("Имя группы: %s\n", groupInfo.Name)

	// 2. Извлечение имен участников
	for _, participant := range groupInfo.Participants {
		// JID участника (например, 79991234567@s.whatsapp.net)
		jid := participant.JID

		// Отображаемое имя (nickname), которое пользователь установил для этой группы,
		// или имя из его профиля, видимое вашему клиенту.
		displayName := participant.DisplayName //.PushName

		// Если PushName пуст, попробуйте использовать Name
		if displayName == "" {
			displayName = participant.JID.User //.Name
		}

		// Если и то, и другое пусто, используем JID
		if displayName == "" {
			displayName = jid.String()
		}

		//fmt.Printf("Участник: %s (JID: %s)\n", displayName, jid.String())
		return displayName
	}
	return ""
}

func (b *Whatsapp) getSenderNotify(senderJid types.JID) string {
	sender, exists := b.contacts[senderJid]

	if !exists || (sender.FullName == "" && sender.PushName == "" && sender.FirstName == "") {
		b.reloadContacts() // Contacts may need to be reloaded
		sender, exists = b.contacts[senderJid]
	}

	if !exists {
		return "someone"
	} else if sender.FullName != "" {
		return sender.FullName
	} else if sender.PushName != "" {
		return sender.PushName
	} else if sender.FirstName != "" {
		return sender.FirstName
	}

	return "someone"
}

func (b *Whatsapp) GetProfilePicThumb(jid string) (*types.ProfilePictureInfo, error) {
	pjid, _ := b.ParseJID(jid)

	info, err := b.wc.GetProfilePictureInfo(context.Background(), pjid, &whatsmeow.GetProfilePictureParams{
		Preview: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get avatar: %v", err)
	}

	return info, nil
}

func isGroupJid(identifier string) bool {
	return strings.HasSuffix(identifier, "@g.us") ||
		strings.HasSuffix(identifier, "@temp") ||
		strings.HasSuffix(identifier, "@broadcast")
}

func (b *Whatsapp) getDevice() (*store.Device, error) {
	device := &store.Device{}

	// Fix: Add context.Background() as first parameter and waLog.Noop logger
	storeContainer, err := sqlstore.New(context.Background(), "sqlite", "file:"+b.cfg.Whatsapp.SessionFile+".db?_pragma=foreign_keys(1)&_pragma=busy_timeout=10000", waLog.Noop)
	if err != nil {
		return device, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Fix: Add context.Background() as first parameter
	device, err = storeContainer.GetFirstDevice(context.Background())
	if err != nil {
		return device, fmt.Errorf("failed to get device: %v", err)
	}

	return device, nil
}

func (b *Whatsapp) getNewReplyContext(parentID string) (*waE2E.ContextInfo, error) {
	replyInfo, err := b.parseMessageID(parentID)
	if err != nil {
		return nil, err
	}

	sender := fmt.Sprintf("%s@%s", replyInfo.Sender.User, replyInfo.Sender.Server)
	ctx := &waE2E.ContextInfo{
		StanzaID:      &replyInfo.MessageID,
		Participant:   &sender,
		QuotedMessage: &waE2E.Message{Conversation: goproto.String("")},
	}

	return ctx, nil
}

func (b *Whatsapp) parseMessageID(id string) (*Replyable, error) {
	// No message ID in case action is executed on a message sent before the bridge was started
	// and then the bridge cache doesn't have this message ID mapped
	if id == "" {
		return &Replyable{MessageID: id}, nil
	}

	replyInfo := strings.Split(id, "/")

	if len(replyInfo) == 2 {
		sender, err := b.ParseJID(replyInfo[0])

		if err == nil {
			return &Replyable{
				MessageID: types.MessageID(replyInfo[1]),
				Sender:    sender,
			}, nil
		}
	}

	err := fmt.Errorf("MessageID does not match format of {senderJID}:{messageID} : \"%s\"", id)

	return &Replyable{MessageID: id}, err
}

func (b *Whatsapp) getParentIdFromCtx(ci *waE2E.ContextInfo) string {
	if ci != nil && ci.StanzaID != nil {
		senderJid, err := b.ParseJID(*ci.Participant)

		if err == nil {
			return getMessageIdFormat(senderJid, *ci.StanzaID)
		}
	}

	return ""
}

func getMessageIdFormat(jid types.JID, messageID string) string {
	// we're crafting our own JID str as AD JID format messes with how stuff looks on a webclient
	//if jid.String() == "79991399754@s.whatsapp.net" {
	//	jid, _ = types.ParseJID("85178361896964@lid")
	//}
	jidStr := fmt.Sprintf("%s@%s", jid.User, jid.Server)
	return fmt.Sprintf("%s/%s", jidStr, messageID)
}

type Replyable struct {
	MessageID types.MessageID
	Sender    types.JID
}

type GroupData struct {
	ChannelId   string
	GuildId     string
	GuildName   string
	GuildAvatar string
}

func (b *Whatsapp) getGroupCommunity(info types.MessageInfo) GroupData {
	g := GroupData{
		ChannelId: info.Chat.String(),
		GuildId:   info.Chat.String(),
		GuildName: "waDM",
	}
	if !info.IsGroup {
		if avatarURL, exists := b.userAvatars[info.Sender.String()]; exists {
			found, newUrl := b.SaveAvatarLocalCache(info.Sender.String(), avatarURL)
			if found {
				g.GuildAvatar = newUrl
			}
		}
		return g
	}

	getProfilePicture := func(groupJID types.JID, isCommunity bool) string {
		profilePicInfo, err := b.wc.GetProfilePictureInfo(context.Background(), groupJID, &whatsmeow.GetProfilePictureParams{
			// false (или nil, если это указатель) обычно означает получение полного изображения
			Preview:     false,
			ExistingID:  "", // Оставляем пустым
			IsCommunity: isCommunity,
		})
		if err != nil || profilePicInfo == nil {
			return ""
		}
		found, newUrl := b.SaveAvatarLocalCache(groupJID.String(), profilePicInfo.URL)
		if found {
			return newUrl
		}
		return profilePicInfo.URL
	}

	if info.IsGroup {
		// Запрашиваем полную информацию о группе
		infoGroup, err := b.wc.GetGroupInfo(context.Background(), info.Chat)
		if err != nil {
			fmt.Println(fmt.Errorf("не удалось получить GroupInfo для %s: %w", info.Chat.String(), err))
		} else if infoGroup != nil {
			g.GuildName = infoGroup.Name
			var isCommunity bool

			if infoGroup.LinkedParentJID.String() != "" {
				communityInfo, err := b.wc.GetGroupInfo(context.Background(), infoGroup.LinkedParentJID)
				if err == nil {
					isCommunity = true
					g.GuildName = communityInfo.Name
					g.ChannelId = fmt.Sprintf("%s/%s", communityInfo.JID.String(), info.Chat.String())
					g.GuildAvatar = getProfilePicture(infoGroup.LinkedParentJID, isCommunity)
					g.GuildId = communityInfo.JID.String()
				}
			} else {
				g.GuildAvatar = getProfilePicture(info.Chat, isCommunity)
			}
		}

	}
	return g
}

func (b *Whatsapp) SaveAvatarLocalCache(userID, url string) (bool, string) {
	if url == "" {
		return false, ""
	}
	folder := "docker/compendium/avatars"

	// Создаем HTTP-запрос
	resp, err := http.Get(url)
	if err != nil {
		b.log.ErrorErr(fmt.Errorf("ошибка при выполнении запроса: %v", err))
		return false, ""
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		b.log.ErrorErr(fmt.Errorf("неправильный статус-код: %d", resp.StatusCode))
		return false, ""
	}

	// Получаем расширение файла из URL
	fileExt := filepath.Ext(url)
	if fileExt == "" {
		b.log.ErrorErr(fmt.Errorf("не удалось определить расширение файла"))
		return false, ""
	}

	// Создаем директорию, если она не существует
	err = os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		b.log.ErrorErr(fmt.Errorf("ошибка при создании директории: %v", err))
		return false, ""
	}

	// Полный путь к файлу
	filename := userID + fileExt
	filePath := filepath.Join(folder, filename)

	newUrl := "https://compendiumnew.mentalisit.myds.me/compendium/avatars/" + filename

	// Проверяем, существует ли файл
	if fileInfo, err := os.Stat(filePath); err == nil {
		// Файл существует, проверяем его размер
		if fileInfo.Size() == resp.ContentLength {
			// Размеры совпадают, пропускаем скачивание
			return true, newUrl
		}
	}

	// Создаем файл
	file, err := os.Create(filePath)
	if err != nil {
		b.log.ErrorErr(fmt.Errorf("ошибка при создании файла: %v", err))
		return false, ""
	}
	defer file.Close()

	// Копируем содержимое ответа в файл
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		b.log.ErrorErr(fmt.Errorf("ошибка при записи файла: %v", err))
		return false, ""
	}

	return true, newUrl
}

func (b *Whatsapp) ParseJID(chatId string) (types.JID, error) {
	if strings.Contains(chatId, "/") {
		split := strings.Split(chatId, "/")
		return types.ParseJID(split[1])
	}
	return types.ParseJID(chatId)
}

func (b *Whatsapp) getJIDForName(name string) (types.JID, bool) {
	if b.mapMentionsNames[name] != nil {
		if b.mapMentionsNames[name].HiddenServer != "" {
			jid, err := types.ParseJID(b.mapMentionsNames[name].HiddenServer)
			if err == nil {
				return jid, true
			}
		} else if b.mapMentionsNames[name].DefaultServer == "" {
			jid, err := types.ParseJID(b.mapMentionsNames[name].DefaultServer)
			if err == nil {
				return jid, true
			}
		}
	}
	return types.JID{}, false
}
func (b *Whatsapp) prepareGroupMentionMessage(originalText string) (mentionedJIDs []string) {
	re := regexp.MustCompile(`@([a-zA-Zа-яА-Я0-9_]+)`)
	matches := re.FindAllStringSubmatch(originalText, -1)

	for _, match := range matches {
		// match[0] - "@ИмяУчастника"
		// match[1] - "ИмяУчастника" (захваченная группа)
		name := match[1]

		jid, found := b.getJIDForName(name) // Используйте вашу функцию сопоставления
		if found {
			mentionedJIDs = append(mentionedJIDs, jid.String())
			// Важно: Текст сообщения должен содержать *отображаемое* имя,
			// а не сам JID. WhatsApp выделит имя, соответствующее JID.
			// Если вы хотите, чтобы в сообщении был виден JID, вам нужно
			// отредактировать originalText, но обычно лучше оставить `@ИмяУчастника`.
		}
	}

	return mentionedJIDs
}
func (b *Whatsapp) getJidDefaultServer(jidDefaultServer string) string {
	//if jidDefaultServer == "79991399754@s.whatsapp.net" {
	//	return jidDefaultServer
	//}
	for _, names := range b.mapMentionsNames {
		if names.DefaultServer == jidDefaultServer {
			return names.HiddenServer
		}
	}
	jid, err := b.ParseJID(jidDefaultServer)
	if err != nil {
		return jidDefaultServer
	}
	pn, err := b.wc.Store.LIDs.GetLIDForPN(context.Background(), jid)
	if err != nil {
		return jidDefaultServer
	}
	return pn.String()
}
