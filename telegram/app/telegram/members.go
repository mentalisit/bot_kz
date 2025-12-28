package telegram

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"telegram/models"
	"time"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

func (t *Telegram) SendChannelDelSecondRsMention(chatid string, text string, parsemode string, second int) (bool, error) {
	if !strings.HasPrefix(text, "MENTION: UserId") {
		return false, nil
	}
	parts := strings.Fields(text)
	if len(parts) != 5 {
		return false, nil
	}
	userID, _ := strconv.ParseInt(parts[2], 10, 64)
	RsTypeLevel := parts[4]
	if !strings.Contains(RsTypeLevel, "rs") {
		return false, nil
	}
	chat, threadID := t.chat(chatid)
	roleId, _ := t.Storage.Db.GetRoleByName(context.Background(), RsTypeLevel, chat)
	if roleId == 0 {
		return false, nil
	}
	users, _ := t.Storage.Db.GetRolesUsers(context.Background(), chat, roleId)
	if len(users) == 0 {
		return false, nil
	}
	var u []models.User
	for _, user := range users {
		if user.ID != userID {
			u = append(u, user)
		}
	}
	var mId string
	send := func(mention []string, parseMode bool) {
		mentionText := strings.Join(mention, " ")
		fullMessage := fmt.Sprintf("🔔 %s\n%s", RsTypeLevel, mentionText)
		m := tgbotapi.NewMessage(chat, fullMessage)
		if parseMode {
			m.ParseMode = "MarkdownV2"
			m.Text = fmt.Sprintf("🔔 %s без никнейма\n%s", RsTypeLevel, mentionText)
		}
		m.MessageThreadID = threadID
		tMessage, err := t.t.Send(m)
		if err != nil {
			if err.Error() != tgbotapi.ErrAPIForbidden {
				t.log.ErrorErr(err)
				t.log.Info(fmt.Sprintf("chatid '%s', text '%s'", chatid, m.Text))
			}
		}
		mId = strconv.Itoa(tMessage.MessageID)
		tu := int(time.Now().UTC().Unix())
		t.Storage.Db.TimerInsert(models.Timer{
			Tip:    "tg",
			ChatId: chatid,
			MesId:  mId,
			Timed:  tu + second,
		})
	}
	if len(u) > 0 {
		var mentions []string
		var mentionsWithoutNickName []string

		for _, member := range u {
			if member.UserName != "" {
				mentions = append(mentions, member.FormatMention())
			} else {
				mentionsWithoutNickName = append(mentionsWithoutNickName, member.FormatMention())
			}

		}

		if len(mentions) > 0 {
			send(mentions, false)
		}
		if len(mentionsWithoutNickName) > 0 {
			send(mentionsWithoutNickName, true)
		}
	}

	if mId != "" {
		return true, nil
	}
	return false, nil
}

func (t *Telegram) logicMention(m *tgbotapi.Message, edit bool) {
	if edit {
		//todo need create logic if edit
		return
	}
	if strings.Contains(m.Text, "@") {
		re := regexp.MustCompile(`@\S+`)
		mentions := re.FindAllString(m.Text, -1)
		if len(mentions) > 0 {
			ThreadID := m.MessageThreadID
			if !m.IsTopicMessage && ThreadID != 0 {
				ThreadID = 0
			}
			ChatId := strconv.FormatInt(m.Chat.ID, 10) + fmt.Sprintf("/%d", ThreadID)
			roles, _ := t.Storage.Db.GetChatsRoles(context.Background(), m.Chat.ID)
			if roles != nil {
				us := make(map[string][]models.User)
				for _, mention := range mentions {
					if mention[1:] == "all" && t.CheckAdminTg(ChatId, m.From.UserName) {
						users, _ := t.Storage.Db.GetChatUsers(context.Background(), m.Chat.ID)
						if len(users) > 0 {
							us["all"] = users
						}
					} else {
						for _, role := range roles {
							if role.Name == mention[1:] {
								users, _ := t.Storage.Db.GetRolesUsers(context.Background(), m.Chat.ID, role.ID)
								if len(users) > 0 {
									us[role.Name] = users
								}
							}
						}
					}
				}
				if len(us) > 0 {
					t.MentionMembersRoles(ChatId, m.MessageID, us)
				}
			}
		}
	}
}

// Функция для упоминания участников по ролям
func (t *Telegram) MentionMembersRoles1(ChatId string, replyId int, trackedMembers map[string][]models.User) {
	text := ""
	text2 := ""
	for roleName, users := range trackedMembers {
		if len(users) == 0 {
			continue
		}

		var mentions []string
		var mentionsWithoutNickName []string

		for _, member := range users {
			if member.UserName != "" {
				mentions = append(mentions, member.FormatMention())
			} else {
				mentionsWithoutNickName = append(mentionsWithoutNickName, member.FormatMention())
			}
		}

		if len(mentions) > 0 {
			// Формируем сообщение
			mentionText := strings.Join(mentions, " ")
			text = fmt.Sprintf("%s\n%s\n%s\n", text, roleName, mentionText)
		}
		if len(mentionsWithoutNickName) > 0 {
			mentionText := strings.Join(mentionsWithoutNickName, " ")
			text2 = fmt.Sprintf("%s\n%s\n%s\n", text2, roleName, mentionText)
		}
	}
	if text != "" {
		_, _ = t.SendChannelReply(ChatId, text, "", replyId)
	}
	if text2 != "" {
		_, _ = t.SendChannelReply(ChatId, text2, "MarkdownV2", replyId)
	}
}
func (t *Telegram) MentionMembersRoles(ChatId string, replyId int, trackedMembers map[string][]models.User) {
	if len(trackedMembers) == 0 {
		return
	}

	var roles []string
	mentions := make(map[string]struct{})
	mentionsNoNick := make(map[string]struct{})

	for roleName, users := range trackedMembers {
		roles = append(roles, roleName)
		for _, member := range users {
			mention := member.FormatMention()
			if member.UserName != "" {
				mentions[mention] = struct{}{}
			} else {
				mentionsNoNick[mention] = struct{}{}
			}
		}
	}

	textRoles := strings.Join(roles, " ")

	// Обработка обычных упоминаний
	if len(mentions) > 0 {
		text := fmt.Sprintf("%s\n\n", textRoles)
		//mList := make([]string, 0, len(mentions))
		for m := range mentions {
			text = fmt.Sprintf("%s %s", text, m)
			//mList = append(mList, m)
		}
		//text := textRoles + "\n" + strings.Join(mList, " ") + "\n"
		_, _ = t.SendChannelReply(ChatId, text, "", replyId)
	}

	// Обработка упоминаний без ника (MarkdownV2)
	if len(mentionsNoNick) > 0 {
		mListNoNick := make([]string, 0, len(mentionsNoNick))
		for m := range mentionsNoNick {
			mListNoNick = append(mListNoNick, m)
		}
		text2 := textRoles + "\n" + strings.Join(mListNoNick, " ") + "\n"
		_, _ = t.SendChannelReply(ChatId, text2, "MarkdownV2", replyId)
	}
}
func (t *Telegram) Unsubscribe(userId int64, argRoles string, guildId int64) int {
	roleId, err := t.Storage.Db.RoleExists(context.Background(), guildId, argRoles)
	if err != nil {
		t.log.Info("role not found")
		return 1
	}
	subscribedRole, _ := t.Storage.Db.IsUserSubscribedToRole(context.Background(), userId, roleId)
	if !subscribedRole {
		return 0
	}
	err = t.Storage.Db.LeaveRole(context.Background(), userId, roleId)
	if err == nil {
		return 2
	}
	return 3
}
func (t *Telegram) Subscribe(userId int64, argRoles string, guildId int64) int {
	roleId, err := t.Storage.Db.RoleExists(context.Background(), guildId, argRoles)
	if err != nil && roleId == 0 {
		_ = t.Storage.Db.CreateRole(context.Background(), &models.Role{ChatID: guildId, Name: argRoles})
		roleId, err = t.Storage.Db.RoleExists(context.Background(), guildId, argRoles)
	}
	if roleId == 0 {
		return 2
	}
	err = t.Storage.Db.JoinRole(context.Background(), userId, roleId, guildId)
	if err != nil {
		return 2
	}

	return 0
}
