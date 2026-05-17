package webServer

import (
	"fmt"
	"net/http"
	"rs/config"
	"rs/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (s *Server) getMAcc(c *gin.Context) {
	uidString := c.Query("uuid")
	tgId := c.Query("tg_id")
	tgUsername := c.Query("tg_username")
	firstName := c.Query("first_name")

	var multiAccount *models.MultiAccount

	if uidString == "" && tgId == "" {
		c.JSON(400, gin.H{"status": "error", "message": "Параметры обязательны"})
		return
	}
	if tgId != "" {
		multiAccount, _ = s.db.FindMultiAccountByUserId(tgId)
		if multiAccount == nil {
			nickname := tgUsername
			if nickname == "" {
				nickname = firstName
			}
			multiAccount, _ = s.db.CreateMultiAccountWithPlatform(tgId, nickname, "tg", tgUsername)
		}
	} else {
		uid, err2 := uuid.Parse(uidString)
		if err2 != nil {
			c.JSON(400, gin.H{"status": "error", "message": "Ошибка чтения айди пользователя"})
			return
		}
		multiAccount, _ = s.db.FindMultiAccountByUUId(uid.String())
		fmt.Printf("%+v\n", multiAccount)
	}

	if multiAccount == nil {
		c.JSON(400, gin.H{"status": "error", "message": "Аккаунт не найден"})
		s.log.Info("multiAccount == nil")
		return
	}

	allModules := s.db.ModuleCompendiumGetAll(multiAccount.UUID)

	guildName := ""
	gidString := c.Query("gid")
	if gidString != "" && gidString != "00000000-0000-0000-0000-000000000000" {
		if parsedGid, err := uuid.Parse(gidString); err == nil {
			if guild, err := s.db.GuildGet(parsedGid); err == nil && guild != nil {
				guildName = guild.GuildName
			}
		}
	}

	c.JSON(200, gin.H{
		"status":     "success",
		"MAcc":       multiAccount,
		"AllModules": allModules,
		"GuildName":  guildName,
	})
}

type SaveMAccRequest struct {
	UUID          string                  `json:"UUID"`
	Nickname      string                  `json:"Nickname"`
	Alts          []string                `json:"Alts"`
	Data          models.MultiAccountData `json:"data"`
	ActiveAccount string                  `json:"active_account"`
	AllModules    []models.Module         `json:"AllModules"`
}

func (s *Server) postMAcc(c *gin.Context) {
	var req SaveMAccRequest

	// Привязываем JSON к структуре
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Некорректный формат данных"})
		return
	}

	// Проверка UUID
	uid, err := uuid.Parse(req.UUID)
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Ошибка формата UUID"})
		return
	}
	fmt.Printf("%+v\n", req)

	ma, _ := s.db.FindMultiAccountByUUId(uid.String())

	if ma != nil && ma.Nickname != req.Nickname {
		_ = s.db.ModuleCompendiumUpdateNickName(uid, ma.Nickname, req.Nickname)

	}
	if len(req.Data.Merged) != 0 || (ma != nil && ma.Data != nil && len(ma.Data.Merged) != 0) {
		s.db.UpdateMergedAccounts(uid, req.Data.Merged)
	}

	if req.Data.WsPhone != "" {
		fmt.Println("saving wsPhone ", req.Data.WsPhone)
		if ma != nil && ma.Data.WsPhone != req.Data.WsPhone {
			s.db.UpdateWsPhone(uid, req.Data.WsPhone)
		}
	}
	if ma != nil && ma.Data != nil && ma.Data.NotifyPM != req.Data.NotifyPM {
		s.db.UpdateNotifyPM(uid, req.Data.NotifyPM)
	}

	// Проверка и сохранение часового пояса и оффсета
	if req.Data.Timezone != "" {
		if loc, err := time.LoadLocation(req.Data.Timezone); err == nil {
			if ma != nil && ma.Data != nil && ma.Data.Timezone != req.Data.Timezone {
				_, offsetSec := time.Now().In(loc).Zone()
				offsetMin := offsetSec / 60
				if err := s.db.UpdateCorpMemberTimezone(uid, req.Data.Timezone, offsetMin); err != nil {
					s.log.ErrorErr(err)
				}
				if err := s.db.UpdateMultiAccountTimezone(uid, req.Data.Timezone); err != nil {
					s.log.ErrorErr(err)
				}
			}
		}
	}

	err = s.db.UpdateMultiAccountNickAltsActive(uid, req.Nickname, req.Alts, req.ActiveAccount)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": "Ошибка при обновлении профиля"})
		return
	}

	err = s.db.ModuleCompendiumCleanDeletedAlts(uid, req.Nickname, req.Alts)
	if err != nil {
		s.log.ErrorErr(err)
	}
	for _, m := range req.AllModules {
		m.Uid = uid
		err = s.db.ModuleCompendiumInsertUpdate(m)
		if err != nil {
			s.log.ErrorErr(err)
		}
	}
	go s.SaveOwner(req)

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Настройки успешно сохранены",
	})
}

func (s *Server) SaveOwner(req SaveMAccRequest) {
	if s.GameNames == nil || len(s.GameNames) == 0 {
		s.GameNames = s.db.GetMergeGameAccountAll()
	}

	check := func(nick string) {

		old, exists := s.GameNames[nick]

		if exists {
			if old.OwnerUuid == "" {
				old.OwnerUuid = req.UUID
				err := s.db.UpdateGamePlayer(old)
				if err != nil {
					s.log.ErrorErr(err)
				}
				s.GameNames[nick] = old
			} else if old.OwnerUuid != req.UUID {
				s.log.Info(fmt.Sprintf("conflict names %+v\n%+v\n", old, req))
			}
		}
	}

	check(req.Nickname)
	for _, alt := range req.Alts {
		check(alt)
	}
	if req.Data.Merged != nil && len(req.Data.Merged) > 0 {
		for _, m := range req.Data.Merged {
			check(m.PlayerName)
		}
	}

}

func (s *Server) postMAccSendMessage(c *gin.Context) {
	type Req struct {
		PlayerId string `json:"playerId"`
		Message  string `json:"message"`
	}

	var r Req

	// Привязываем JSON к структуре
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Некорректный формат данных"})
		return
	}

	// Проверка UUID
	uid, err := uuid.Parse(r.PlayerId)
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Ошибка формата UUID"})
		return
	}

	ma, err := s.db.FindMultiAccountByUUId(uid.String())

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	if ma == nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if ma.DiscordID != "" {
		s.cl.Ds.SendDmText(r.Message, ma.DiscordID)
	}
	if ma.TelegramID != "" {
		s.cl.Tg.SendChannelId(ma.TelegramID, r.Message)
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "sending message",
	})
}

func (s *Server) getSecret(c *gin.Context) {
	secretToken := c.Query("secretToken")
	if secretToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Параметры обязательны"})
		return
	}

	// Убираем префикс
	after, found := strings.CutPrefix(secretToken, "my_compendium_")
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Неверный формат токена"})
		return
	}

	// Парсим токен
	mapClaims, err := parseToken(after)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Ошибка парсинга токена"})
		return
	}

	// Безопасное извлечение UUID из claims
	// Проверяем, что ключ существует и является строкой, чтобы избежать паники
	uuidStr, ok := mapClaims["uuid"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "UUID не найден в токене"})
		return
	}

	uid, err := uuid.Parse(uuidStr)
	if err != nil || uid == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Некорректный UUID"})
		return
	}

	// Формируем URL для перехода
	// В вашем случае это API запрос, который возвращает данные пользователя
	targetURL := fmt.Sprintf("https://mentalisit.myds.me/rs/settings/compendium.html?uuid=%s", uid.String())

	// Выполняем переадресацию (302 Found)
	c.Redirect(http.StatusFound, targetURL)
}

func parseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.Instance.Postgress.Password), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("cannot parse claims")
	}

	return claims, nil
}
