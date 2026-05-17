package webServer

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) checkGA(c *gin.Context) {
	uidString := c.Query("uuid")
	if uidString == "" {
		c.JSON(400, gin.H{"status": "error", "message": "Параметры обязательны"})
		return
	}
	_, err2 := uuid.Parse(uidString)
	if err2 != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Ошибка чтения айди пользователя"})
		return
	}
	ga := c.Query("ga")
	if ga == "" {
		c.JSON(400, gin.H{"status": "error", "message": "Аккаунт не запрошен"})
		return
	}
	fmt.Printf("Check Game Account = %s\n", ga)
	gameAccounts := s.db.GetMergeGameAccount(ga)
	if len(gameAccounts) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Аккаунт не найден"})
		return
	}
	if len(gameAccounts) != 1 {
		s.log.Info(fmt.Sprintf("game account %s not one , found  %d\n", ga, len(gameAccounts)))
		c.JSON(http.StatusConflict, gin.H{"status": "error", "message": "Аккаунт не один"})
		return
	}
	gameAccount := gameAccounts[0]
	if gameAccount.OwnerUuid != "" {
		owner, _ := s.db.FindMultiAccountByUUId(gameAccount.OwnerUuid)
		if owner != nil {
			gameAccount.CurrentOwner = owner.Nickname
		}
	}

	c.JSON(200, gin.H{
		"status": "success",
		"GA":     gameAccount,
	})
}
