package logic

import (
	"compendium/models"
	"fmt"
	"regexp"
	"strings"
)

//TODO NEED TRANSLATE

func (c *Hs) logicRoles(m models.IncomingMessage) bool {
	cutPrefix, _ := strings.CutPrefix(m.Text, "%")
	// Компиляция регулярного выражения
	regex, err := regexp.Compile(`^role (create|delete) (\w+)$`)
	if err != nil {
		c.log.ErrorErr(err)
		return false
	}

	matches := regex.FindStringSubmatch(cutPrefix)
	if matches != nil {
		action := matches[1]
		roleName := matches[2]
		ifExistRole := c.guildsRole.GuildRoleExist(m.MultiGuild.GuildId(), roleName)

		if action == "create" {
			if ifExistRole {
				c.sendChat(m, roleName+" роль уже существует")
			} else {
				c.guildsRole.GuildRoleCreate(m.MultiGuild.GuildId(), roleName)
				c.sendChat(m, roleName+" роль создана")
			}
		}
		if action == "delete" {
			if ifExistRole {
				err = c.guildsRole.GuildRoleDelete(m.MultiGuild.GuildId(), roleName)
				if err != nil {
					c.log.ErrorErr(err)
					return false
				}
				c.sendChat(m, roleName+" роль удалена")
			} else {
				c.sendChat(m, roleName+" не существует")
			}
		}

		return true
	}

	// Компиляция регулярного выражения
	regex, err = regexp.Compile(`^role (s|u) (\w+)$`)
	if err != nil {
		c.log.ErrorErr(err)
		return false
	}

	matches = regex.FindStringSubmatch(cutPrefix)
	if matches != nil {
		action := matches[1]
		roleName := matches[2]
		existRole := c.guildsRole.GuildRoleExist(m.MultiGuild.GuildId(), roleName)
		existSubscribe := c.guildsRole.GuildRolesExistSubscribe(m.MultiGuild.GuildId(), roleName, m.NameId)
		if action == "s" {
			if existRole {
				if existSubscribe {
					c.sendChat(m, "ты уже подписан на роль "+roleName)
				} else {
					err = c.guildsRole.GuildRolesSubscribe(m.MultiGuild.GuildId(), roleName, m.Name, m.NameId)
					if err != nil {
						c.log.ErrorErr(err)
						return false
					}
					c.sendChat(m, "подписался на роль "+roleName)
				}
			} else {
				c.sendChat(m, fmt.Sprintf("роли %s не существует, сначала создай роль\n команда: %%role create %s", roleName, roleName))
			}
		}
		if action == "u" {
			if existRole {
				if existSubscribe {
					err = c.guildsRole.GuildRolesDeleteSubscribe(m.MultiGuild.GuildId(), roleName, m.NameId)
					if err != nil {
						c.log.ErrorErr(err)
						return false
					}
					c.sendChat(m, "отписался от роли "+roleName)
				} else {
					c.sendChat(m, "ты не подписан на "+roleName)
				}
			} else {
				c.sendChat(m, "не существует роли "+roleName)
			}
		}
		return true
	}

	reSubs := regexp.MustCompile(`^role subs (\w+)((?: @\w+)+)$`)

	matches = reSubs.FindStringSubmatch(cutPrefix)
	if matches != nil && len(matches) > 2 {
		text := "попытка оформить подписку \n"
		roleName := matches[1]
		usernames := matches[2]
		if !c.guildsRole.GuildRoleExist(m.MultiGuild.GuildId(), roleName) {
			text += "сначала создай роль " + roleName
			c.sendChat(m, text)
			return true
		}

		split := strings.Split(usernames, " ")
		for _, s := range split {
			var user *models.User
			//если с упоминанием
			after, found := strings.CutPrefix(s, "@")

			if found {
				user, _ = c.users.UsersGetByUserName(after)
			} else {
				user, _ = c.UsersFindByAltsNameOrGameName(s)
			}
			fmt.Println(user)

			if user != nil && user.ID != "" {
				err = c.guildsRole.GuildRolesSubscribe(m.MultiGuild.GuildId(), roleName, user.Username, user.ID)
				if err != nil {
					fmt.Println(err)
					return false
				}
				text += fmt.Sprintf("%s подписан\n", after)
			} else {
				text += fmt.Sprintf("%s не подписан данные не найдены \n", after)
			}
		}
		c.sendChat(m, text)
		return true
	}

	reUnSubs := regexp.MustCompile(`^role unsubs (\w+)((?: @\w+)+)$`)

	matches = reUnSubs.FindStringSubmatch(cutPrefix)
	if matches != nil && len(matches) > 2 {
		text := "попытка отменить подписку \n"
		roleName := matches[1]
		usernames := matches[2]
		if !c.guildsRole.GuildRoleExist(m.MultiGuild.GuildId(), roleName) {
			text += "сначала создай роль " + roleName
			c.sendChat(m, text)
			return true
		}

		split := strings.Split(usernames, " ")
		for _, s := range split {
			var user *models.User
			//если с упоминанием
			after, found := strings.CutPrefix(s, "@")

			if found {
				user, _ = c.users.UsersGetByUserName(after)
			} else {
				user, _ = c.UsersFindByAltsNameOrGameName(s)
			}
			fmt.Println(user)

			if user != nil && user.ID != "" {
				err = c.guildsRole.GuildRolesDeleteSubscribeUser(m.MultiGuild.GuildId(), roleName, user.Username, user.ID)
				if err != nil {
					fmt.Println(err)
					return false
				}
				text += fmt.Sprintf("%s отписан\n", after)
			} else {
				text += fmt.Sprintf("%s не отписан данные не найдены \n", after)
			}
		}
		c.sendChat(m, text)
		return true
	}
	return false
}

func (c *Hs) UsersFindByAltsNameOrGameName(AltNameOrGameName string) (*models.User, error) {
	usersGetAll, err := c.users.UsersGetAll()
	if err != nil {
		return nil, err
	}
	for _, user := range usersGetAll {
		if len(user.Alts) > 0 {
			for _, alt := range user.Alts {
				if alt == AltNameOrGameName {
					return &user, nil
				}
			}
		} else if user.GameName == AltNameOrGameName {
			return &user, nil
		}
	}
	return nil, nil
}
