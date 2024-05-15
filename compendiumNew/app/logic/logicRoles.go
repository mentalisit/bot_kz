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
		ifExistRole := c.guildsRole.GuildRoleExist(m.GuildId, roleName)

		if action == "create" {
			if ifExistRole {
				c.sendChat(m, roleName+" роль уже существует")
			} else {
				c.guildsRole.GuildRoleCreate(m.GuildId, roleName)
				c.sendChat(m, roleName+" роль создана")
			}
		}
		if action == "delete" {
			if ifExistRole {
				err = c.guildsRole.GuildRoleDelete(m.GuildId, roleName)
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
		existRole := c.guildsRole.GuildRoleExist(m.GuildId, roleName)
		existSubscribe := c.guildsRole.GuildRolesExistSubscribe(m.GuildId, roleName, m.NameId)
		if action == "s" {
			if existRole {
				if existSubscribe {
					c.sendChat(m, "ты уже подписан на роль "+roleName)
				} else {
					err = c.guildsRole.GuildRolesSubscribe(m.GuildId, roleName, m.Name, m.NameId)
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
					err = c.guildsRole.GuildRolesDeleteSubscribe(m.GuildId, roleName, m.NameId)
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
		split := strings.Split(usernames, " ")
		for _, s := range split {
			after, found := strings.CutPrefix(s, "@")
			if found {
				if c.guildsRole.GuildRoleExist(m.GuildId, roleName) {
					user, errget := c.users.UsersGetByUserName(after)
					if errget != nil {
						c.log.ErrorErr(errget)
						return false
					}

					if user.ID != "" {
						err = c.guildsRole.GuildRolesSubscribe(m.GuildId, roleName, user.Username, user.ID)
						if err != nil {
							return false
						}
						text += fmt.Sprintf("%s подписан\n", after)
					} else {
						text += fmt.Sprintf("%s не подписан данные не найдены \n", after)
					}
				} else {
					text += "сначала создай роль " + roleName
					break
				}
			}
		}
		c.sendChat(m, text)
		return true
	}
	return false
}
