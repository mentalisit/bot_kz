package server

import (
	"compendium_s/models"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type cacheEntry struct {
	data      any       // Сами данные, которые отправляем
	timestamp time.Time // Когда сохранили
}

// Возвращает данные из кэша, если они свежие
func (s *Server) getFreshCache(key string) (any, bool) {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	entry, exists := s.cacheReq[key]
	if !exists {
		return nil, false
	}

	if time.Since(entry.timestamp) <= 2*time.Second {
		return entry.data, true
	}

	delete(s.cacheReq, key)
	return nil, false
}

// Записывает данные в кэш
func (s *Server) setCache(key string, data any) {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	s.cacheReq[key] = cacheEntry{
		data:      data,
		timestamp: time.Now(),
	}
}

func (s *Server) fetchCorpData(token, roleId, mGuild string) (any, int, error) {
	i := s.GetTokenIdentity(token)
	if strings.HasPrefix(token, "my_compendium_") {
		if i != nil && i.MAccount != nil && i.MAccount.Nickname != "" {
			i.MGuild, _ = s.dbV2.GuildGetById(mGuild)
			result := s.GetCorpDataV2Internal(i, roleId)
			return result, http.StatusOK, nil
		}
	}

	if i == nil {
		return nil, http.StatusForbidden, errors.New("invalid token")
	}

	if mGuild == "" {
		// Логика для запроса БЕЗ указания конкретной гильдии
		result := s.GetCorpDataMultiGuild(i, roleId)
		fmt.Printf("resultf %+v\n", result)
		return result, http.StatusOK, nil
	}

	// Логика для запроса С указанием конкретной гильдии (mGuild != "")
	return s.fetchSpecificCorpData(i, roleId, mGuild)
}
func (s *Server) findAndPrepareGuild(mGuild string) (*models.MultiAccountGuildV2, error) {
	var multiGuild *models.MultiAccountGuildV2 // MultiAccountGuild для Identity

	// Попытка 1: Поиск как Multi-Guild по UUID
	if gid, err := uuid.Parse(mGuild); err == nil {
		mg, err := s.dbV2.GuildGet(gid)
		if err == nil && mg != nil {
			return mg, nil
		}
	}

	return multiGuild, nil
}
func (s *Server) fetchSpecificCorpData(i *models.Identity, roleId, mGuild string) (*models.CorpData, int, error) {
	multiGuild, err := s.findAndPrepareGuild(mGuild)
	if err != nil {
		// Статус 404, если гильдия не найдена, 500 для ошибок БД
		status := http.StatusNotFound
		if strings.Contains(err.Error(), "DB") {
			status = http.StatusInternalServerError
		}
		return nil, status, err
	}

	// Создаем временный Identity
	tempIdentity := &models.Identity{
		User:         i.User,
		Token:        i.Token,
		MultiAccount: i.MultiAccount,
		MGuild:       multiGuild, // Используем найденный mGuild
	}

	// Вызываем финальную внутреннюю логику
	result := s.GetCorpDataMultiGuild(tempIdentity, roleId)
	return result, http.StatusOK, nil
}

func (s *Server) GetCorpDataV2Internal(i *models.Identity, roleId string) *models.CorpData {
	c := models.CorpData{}
	c.Initialization()

	members, _ := s.dbV2.CorpMembersReadMulti(&i.MGuild.GId)

	appendRolesByType := func(roles []models.CorpRole, nameMSG string) {
		if len(roles) > 0 {
			for _, role := range roles {
				if role.Name == "@everyone" || role.Name == "all" {
					continue
				}
				c.Roles = append(c.Roles, models.CorpRole{
					Id:       role.Id,
					Name:     fmt.Sprintf("(%s) %s", strings.ToUpper(nameMSG), role.Name),
					TypeRole: nameMSG,
				})
			}
		}
	}

	cOld := s.GetCorpDataMultiGuild(i, roleId)
	if cOld != nil {
		members = append(members, cOld.Members...)
		appendRolesByType(cOld.Roles, "")
	}

	guildDs, guildTg, guildWa := getCorpsTypeIdV2(i.MGuild)

	CheckRoleDs := func(role models.CorpRole) {
		for _, gId := range guildDs {
			fmt.Printf("members total Len %d\nCheckRoleDs gId %+v\n", len(members), gId)
			for _, member := range members {
				if member.GetType() == role.TypeRole {
					uid := member.Multi.DiscordID
					if strings.Contains(member.Multi.DiscordID, "/") {
						split := strings.Split(member.Multi.DiscordID, "/")
						uid = split[0]
					}
					if s.roles.CheckRoleDs(gId, uid, roleId) {
						fmt.Printf("CheckRoleDs OK roleId %s member %+v\n", roleId, member)
						c.Members = append(c.Members, member)
					}
				}
			}
		}
	}
	CheckRoleTg := func(role models.CorpRole) {
		for _, member := range members {
			if member.GetType() == role.TypeRole {
				exist, _ := s.dbV2.IsUserSubscribedToRole(member.Multi.GetTelegramChatId(), role.GetRoleId())
				if exist {
					c.Members = append(c.Members, member)
				}
			}
		}
	}

	if len(guildTg) != 0 {
		c.AppendEveryone("tg")
		for _, tgId := range guildTg {
			gTg, _ := strconv.ParseInt(tgId, 10, 64)
			roles, err := s.dbV2.GetChatsRoles(gTg)
			if err != nil {
				s.log.ErrorErr(err)
			}
			if len(roles) != 0 {
				appendRolesByType(roles, "tg")
			}
		}

	}
	if len(guildDs) != 0 {
		c.AppendEveryone("ds")
		for _, gds := range guildDs {
			s.roles.LoadGuild(gds)
			roles, err := s.roles.ds.GetRoles(gds)
			if err != nil {
				s.log.ErrorErr(err)
			}
			appendRolesByType(roles, "ds")
		}
	}
	if len(guildWa) != 0 {
		c.AppendEveryone("wa")
	}

	if roleId == "" {
		c.Members = members
	} else if roleId == "tg" || roleId == "ds" || roleId == "ma" {
		for _, member := range members {
			if member.GetType() == roleId {
				c.Members = append(c.Members, member)
			}
		}
	} else {
		var role models.CorpRole
		for _, roles := range c.Roles {
			if roleId == roles.Id {
				role = roles
			}
		}
		if role.TypeRole == "" {
			s.log.Info(fmt.Sprintf("role type not exist %+v\n", role))
		}

		if role.TypeRole == "ds" {
			CheckRoleDs(role)
		} else if role.TypeRole == "tg" {
			CheckRoleTg(role)
		} else {
			s.log.Info(fmt.Sprintf("role req %s found %+v\n", roleId, role))
		}
	}
	if len(c.Members) != 0 {
		sort.Slice(c.Members, func(i, j int) bool {
			// Проверка, чтобы индекс не выходил за пределы строки
			nameI := c.Members[i].Name
			nameJ := c.Members[j].Name

			// Игнорируем первые пять символов, если длина имени больше или равна пяти
			if len(nameI) >= 5 {
				nameI = nameI[5:]
			}
			if len(nameJ) >= 5 {
				nameJ = nameJ[5:]
			}

			return nameI < nameJ
		})
	}

	return &c
}
