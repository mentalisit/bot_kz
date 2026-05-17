package webServer

import (
	"context"
	"fmt"
	"net/http"
	"rs/models"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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

func (s *Server) compendiumTech(c *gin.Context) {
	uidString := c.Query("uuid")
	name := c.Query("name")
	uid, _ := uuid.Parse(uidString)

	if c.Request.Method == "GET" {
		technologiesGet, err := s.db.TechnologiesGet(uid, name)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, technologiesGet)

	} else if c.Request.Method == "POST" {
		var req models.CompendiumTechReq

		// Привязываем JSON к структуре
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"status": "error", "message": "Некорректный формат данных"})
			return
		}

		techLevels, err := s.db.TechnologyInsertUpdate(req)
		if err != nil {
			c.JSON(400, gin.H{"status": "error", "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, techLevels)

		fmt.Printf("req %+v\n", req)
	}
}

func (s *Server) compendiumCorps(c *gin.Context) {
	uidString := c.Query("uuid")
	uid, _ := uuid.Parse(uidString)
	if uid == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}

	mAcc, _ := s.db.FindMultiAccountByUUId(uid.String())
	corporations, err := s.db.UserCorporationsGet(mAcc)
	if err != nil {
		s.log.ErrorErr(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user corporations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":         mAcc,
		"corporations": corporations,
	})
}

func (s *Server) compendiumAccessibleCorporations(c *gin.Context) {
	uidString := c.Query("uuid")
	uid, _ := uuid.Parse(uidString)
	if uid == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}

	ginH := make(gin.H)

	mAcc, _ := s.db.FindMultiAccountByUUId(uid.String())
	ginH["user"] = mAcc

	if mAcc.TelegramID != "" {
		userId, _ := strconv.ParseInt(mAcc.TelegramID, 10, 64)
		memberTG := s.cl.Tg.GetChatsMember(userId, false)
		ginH["chats_tg"] = memberTG
	}
	if mAcc.DiscordID != "" {
		//getDiscordChats
		//ginH["chats_ds"] = memberDS
	}

	c.JSON(http.StatusOK, ginH)
}

func (s *Server) compendiumMultiCorp(c *gin.Context) {
	uidString := c.Query("uuid")
	uid, _ := uuid.Parse(uidString)
	if uid == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}
	chatId := c.Query("chat_id")
	//platform := c.Query("platform")
	//status := c.Query("status")

	ginH := make(gin.H)

	mAcc, _ := s.db.FindMultiAccountByUUId(uid.String())
	ginH["user"] = mAcc

	guild, err := s.db.GuildGetChannel(context.Background(), chatId)
	if err != nil {
		s.log.ErrorErr(err)
	}
	ginH["guild"] = guild

	corpMember, err := s.db.CorpMemberByUId(context.Background(), mAcc.UUID)
	if err != nil {
		s.log.ErrorErr(err)
	}
	ginH["corp_member"] = corpMember

	c.JSON(http.StatusOK, ginH)
}

func (s *Server) compendiumCorpData(c *gin.Context) {
	roleId := c.Query("roleId")
	gidString := c.Query("gid")
	uidString := c.Query("uuid")
	if uidString == "" || gidString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing uuid"})
		return
	}

	gid, err := uuid.Parse(gidString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid corpId"})
		return
	}

	mg, err := s.db.GuildGet(gid)
	if err != nil || mg == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "corporation not found"})
		return
	}

	guildDs, guildTg, guildWa := getCorpsTypeIdV2(mg)

	corp := models.CorpData{
		Members: []models.CorpMember{},
		Roles: []models.CorpRole{{
			ID:   "",
			Name: "@everyone",
		}},
		MGuild: *mg,
	}

	members, _ := s.db.CorpMembersReadMulti(&gid)

	appendRolesByType := func(roles []models.CorpRole, nameMSG string) {
		if len(roles) > 0 {
			for _, role := range roles {
				if role.Name == "@everyone" || role.Name == "all" {
					continue
				}
				corp.Roles = append(corp.Roles, models.CorpRole{
					ID:       role.ID,
					Name:     fmt.Sprintf("%s (%s) ", role.Name, strings.ToUpper(nameMSG)),
					TypeRole: nameMSG,
				})
			}
		}
	}

	CheckRoleDs := func(targetRoleId string) {
		for _, gId := range guildDs {
			s.roles.LoadGuild(gId) // Используем rolesHelper
			for _, member := range members {
				if member.GetType() == "ds" {
					uid := member.Multi.DiscordID
					if strings.Contains(member.Multi.DiscordID, "/") {
						split := strings.Split(member.Multi.DiscordID, "/")
						uid = split[0]
					}
					if s.roles.CheckRoleDs(gId, uid, targetRoleId) {
						corp.Members = append(corp.Members, member)
					}
				}
			}
		}
	}

	CheckRoleTg := func(targetRoleId string) {
		for _, member := range members {
			if member.GetType() == "tg" {
				roleInt, _ := strconv.ParseInt(targetRoleId, 10, 64)
				exist, _ := s.db.IsUserSubscribedToRole(member.Multi.GetTelegramChatId(), roleInt)
				if exist {
					corp.Members = append(corp.Members, member)
				}
			}
		}
	}

	if len(guildTg) != 0 {
		corp.AppendEveryone("tg")
		for _, tgId := range guildTg {
			gTg, _ := strconv.ParseInt(tgId, 10, 64)
			roles, err := s.db.GetChatsRoles(gTg)
			if err != nil {
				s.log.ErrorErr(err)
			}
			if len(roles) != 0 {
				appendRolesByType(roles, "tg")
			}
		}
	}

	if len(guildDs) != 0 {
		corp.AppendEveryone("ds")
		for _, gds := range guildDs {
			s.roles.LoadGuild(gds)
			roles, err := s.cl.Ds.GetRoles(gds)
			if err != nil {
				s.log.ErrorErr(err)
			}
			if len(roles) != 0 {
				appendRolesByType(roles, "ds")
			}
		}
	}

	// Note: s.cl.Ds in rs_bot2 doesn't have GetRoles method in function.go but it might be in discord.pb.go
	// Let's assume for now we need a way to get roles.

	if len(guildWa) != 0 {
		corp.AppendEveryone("wa")
	}

	if roleId == "" {
		corp.Members = members
	} else if roleId == "tg" || roleId == "ds" || roleId == "ma" {
		for _, member := range members {
			if member.GetType() == roleId {
				corp.Members = append(corp.Members, member)
			}
		}
	} else {
		// Specific role filtering
		var foundRole models.CorpRole
		for _, r := range corp.Roles {
			if roleId == r.ID {
				foundRole = r
				break
			}
		}

		if foundRole.TypeRole == "ds" {
			CheckRoleDs(roleId)
		} else if foundRole.TypeRole == "tg" {
			CheckRoleTg(roleId)
		} else {
			// If role not found in corp.Roles (e.g. it's a global role ID), try all types?
			// Compendium_server logic:
			CheckRoleDs(roleId)
			CheckRoleTg(roleId)
		}
	}

	if len(corp.Members) != 0 {
		sort.Slice(corp.Members, func(i, j int) bool {
			nameI := corp.Members[i].Name
			nameJ := corp.Members[j].Name
			if len(nameI) >= 5 {
				nameI = nameI[5:]
			}
			if len(nameJ) >= 5 {
				nameJ = nameJ[5:]
			}
			return nameI < nameJ
		})
	}

	c.JSON(http.StatusOK, corp)
}

func getCorpsTypeIdV2(mg *models.MultiAccountGuildV2) (guildDs, guildTg, guildWa []string) {
	for t, channel := range mg.Channels {
		switch t {
		case "wa":
			guildWa = append(guildWa, channel...)
		case "ds":
			guildDs = append(guildDs, channel...)
		case "tg":
			guildTg = append(guildTg, channel...)
		}
	}
	return guildDs, guildTg, guildWa
}

type regMultiCorp struct {
	Uuid     string `json:"uuid"`
	ChatId   string `json:"chat_id"`
	Platform string `json:"platform"`
	ChatName string `json:"chat_name"`
	Status   string `json:"status"`
}

func (s *Server) compendiumMultiCorpRegister(c *gin.Context) {
	var req regMultiCorp

	// Привязываем JSON к структуре
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Некорректный формат данных"})
		return
	}
	mg := models.MultiAccountGuildV2{
		GId:       uuid.UUID{},
		GuildName: req.ChatName,
		Channels:  make(models.GuildChannels),
	}
	mg.Channels[req.Platform] = append(mg.Channels[req.Platform], req.ChatId)
	saveGuild, err := s.db.GuildSave(context.Background(), mg)
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, saveGuild)

}

type participate struct {
	Uuid         string `json:"uuid"`
	ParentId     string `json:"parent_id"`
	ParentStatus string `json:"parent_status"`
	Platform     string `json:"platform"`
	CorpId       string `json:"corp_id"`
	Active       bool   `json:"active"`
}

func (s *Server) compendiumToggleParticipation(c *gin.Context) {
	var req participate

	// Привязываем JSON к структуре
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Некорректный формат данных"})
		return
	}

	uid, _ := uuid.Parse(req.Uuid)

	corpMember, err := s.db.CorpMemberByUId(context.Background(), uid)
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": err.Error()})
		return
	}
	gid, _ := uuid.Parse(req.CorpId)

	if corpMember == nil {
		cm := models.MultiAccountCorpMember{
			Uid:        uid,
			GuildIds:   models.UUIDArray{},
			TimeZona:   "",
			ZonaOffset: 0,
			AfkFor:     "",
		}
		if req.Active {
			cm.GuildIds = append(cm.GuildIds, gid)
		}

		corpMember = &cm
		err = s.db.CorpMemberInsert(context.Background(), cm)
		if err != nil {
			c.JSON(400, gin.H{"status": "error", "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, "ok")
		return
	}
	if req.Active {
		corpMember.GuildIds = append(corpMember.GuildIds, gid)
	} else if !req.Active {
		corpMember.GuildIds.Remove(gid)
	}
	err = s.db.CorpMemberUpdate(context.Background(), *corpMember)
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "ok")
}
