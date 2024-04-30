package Compendium

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

func (c *Compendium) logicRoles() bool {
	cutPrefix, _ := strings.CutPrefix(c.in.Text, "%")
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

		ifExistRole := c.db.Temp.ExistRole(c.in.GuildId, roleName)

		if action == "create" {
			if ifExistRole {
				c.sendChat(roleName + " роль уже существует")
			} else {
				c.db.Temp.CreateRole(c.in.GuildId, roleName)
				c.sendChat(roleName + " роль создана")
			}
		}
		if action == "delete" {
			if ifExistRole {
				c.db.Temp.DeleteRole(c.in.GuildId, roleName)
				c.sendChat(roleName + " роль удалена")
			} else {
				c.sendChat(roleName + " не существует")
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
		existRole := c.db.Temp.ExistRole(c.in.GuildId, roleName)
		existSubscribe := c.db.Temp.ExistSubscribe(c.in.GuildId, roleName, c.in.NameId)
		if action == "s" {
			if existRole {
				if existSubscribe {
					c.sendChat("ты уже подписан на роль " + roleName)
				} else {
					c.db.Temp.RoleSubscribe(c.in.GuildId, roleName, c.in.Name, c.in.NameId)
					c.sendChat("подписался на роль " + roleName)
				}
			} else {
				c.sendChat(fmt.Sprintf("роли %s не существует, сначала создай роль\n команда: %%role create %s", roleName, roleName))
			}
		}
		if action == "u" {
			if existRole {
				if existSubscribe {
					c.db.Temp.DeleteSubscribe(c.in.GuildId, roleName, c.in.NameId)
					c.sendChat("отписался от роли " + roleName)
				} else {
					c.sendChat("ты не подписан на " + roleName)
				}
			} else {
				c.sendChat("не существует роли " + roleName)
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
				if c.db.Temp.ExistRole(c.in.GuildId, roleName) {
					member := c.db.Temp.CorpMemberReadByNameByGuildId(context.Background(), after, c.in.GuildId)
					if member.UserId != "" {
						c.db.Temp.RoleSubscribe(c.in.GuildId, roleName, c.in.NameId, member.UserId)
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
		c.sendChat(text)
		return true
	}
	return false
}
