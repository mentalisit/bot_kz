package server

import (
	"compendium_s/models"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

func (s *Server) GetTokenIdentity(token string) *models.Identity {
	userid, guildid, err := s.db.ListUserGetUserIdAndGuildId(token)
	if err != nil {
		token2 := s.db.ListUserGetByMatch(token)
		userid, guildid, err = s.db.ListUserGetUserIdAndGuildId(token2)
		if err != nil {
			s.log.Info("get user by token: " + token + " " + err.Error())
			return nil
		}
	}
	var i models.Identity
	i.Token = token
	user, err := s.db.UsersGetByUserId(userid)
	if err != nil {
		s.log.ErrorErr(err)
		s.log.Info("get user by userid :" + userid)
		return nil
	}
	i.User = *user
	guild, err := s.db.GuildGet(guildid)
	if err != nil {
		s.log.ErrorErr(err)
		s.log.Info("get guild by guildid:" + guildid)
		return nil
	}
	i.Guild = *guild

	return &i
}

func (s *Server) GetCorpData(i *models.Identity, roleId string) *models.CorpData {
	c := models.CorpData{}
	c.Members = []models.CorpMember{}

	if i.Guild.ID != "" {
		c.Roles = s.getRoles(i)
		cm, err := s.db.CorpMembersRead(i.Guild.ID)
		if err != nil {
			s.log.ErrorErr(err)
			return nil
		}
		for _, member := range cm {
			if i.Guild.Type == "tg" {
				if roleId == "" {
					c.Members = append(c.Members, member)
				} else {
					roles, er := s.db.GuildRolesRead(i.Guild.ID)
					if er != nil {
						s.log.ErrorErr(er)
						return nil
					}
					for _, role := range roles {
						if role.Id == roleId {
							if s.db.GuildRolesExistSubscribe(i.Guild.ID, role.Name, member.UserId) {
								c.Members = append(c.Members, member)
							}
						}
					}
				}
			} else if i.Guild.Type == "ds" {
				uid := member.UserId
				if strings.Contains(member.UserId, "/") {
					split := strings.Split(member.UserId, "/")
					uid = split[0]
				}
				if CheckRoleDs(i.Guild.ID, uid, roleId) {
					c.Members = append(c.Members, member)
				}
			}
		}
	}
	return &c
}
func (s *Server) getRoles(i *models.Identity) []models.CorpRole {
	if i.Guild.Type == "tg" {
		everyone := []models.CorpRole{{
			Id:   "",
			Name: "Telegram",
		}}
		roles, err := s.db.GuildRolesRead(i.Guild.ID)
		if err != nil {
			s.log.ErrorErr(err)
			return nil
		}
		if len(roles) > 0 {
			everyone = append(everyone, roles...)
		}
		return everyone
	} else {
		roles, err := GetRoles(i.Guild.ID)
		if err != nil {
			s.log.ErrorErr(err)
			return nil
		}
		return roles
	}
}

func (s *Server) CheckCode(code string) models.Identity {
	var i models.Identity
	coder, err := s.db.CodeGet(code)
	if err != nil {
		fmt.Println("CheckCode " + err.Error())
	}

	if coder != nil && coder.Code == code {
		if time.Now().Unix() < coder.Timestamp+600 {
			i = coder.Identity
		}
	}
	if code == "test-test-test" {
		i = models.Identity{
			User: models.User{
				ID:       "111111111",
				Username: "TestUser",
				Alts:     []string{"alt1", "alt2"},
			},
			Guild: models.Guild{
				ID:   "22222222222",
				Name: "TestGuild",
				Type: "tg",
			},
			Token: "gGUBIlUAU1uTKWd8HssP27ojG0DugoAaPslwFGTDSAbEM6UM",
		}

	}
	return i
}
func (s *Server) refreshToken(token string) string {
	if len(token) < 60 {
		newToken := s.checkPrefixToken(token)
		err := s.db.ListUserUpdateToken(token, newToken)
		if err != nil {
			return token
		}
		return newToken
	}
	return token
}
func GenerateToken() string {
	// Вычисляем необходимый размер байт для указанной длины токена
	tokenBytes := make([]byte, 174)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return ""
	}

	// Кодируем байты в строку base64
	token := base64.URLEncoding.EncodeToString(tokenBytes)
	return token
}

func (s *Server) checkPrefixToken(token string) string {
	userid, guildid, err := s.db.ListUserGetUserIdAndGuildId(token)
	if err != nil || userid == "" || guildid == "" {
		s.log.Info("err || userid == nil || guildid == nil")
		return token
	}
	guildGet, err := s.db.GuildGet(guildid)
	if err != nil || guildGet.Name == "" {
		s.log.Info("err || guildGet.Name == nil")
		return token
	}

	newToken := guildGet.Type + guildid + "." + userid + GenerateToken()
	return newToken
}
