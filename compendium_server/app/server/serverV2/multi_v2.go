package serverV2

import (
	"compendium_s/models"
	"fmt"
	"strings"
)

//
//import (
//	"compendium_s/models"
//	"fmt"
//	"net/http"
//	"sort"
//	"strings"
//
//	"github.com/gin-gonic/gin"
//	"github.com/google/uuid"
//)

//func (s *ServerV2) SyncTechMultiGuild(c *gin.Context, i *models.IdentityV2, mode, twin string) {
//	userName := i.User.Username
//
//	if twin != "" && twin != "default" {
//		userName = twin
//	}
//
//	userUUID, err := uuid.Parse(i.User.ID)
//	if err != nil {
//		s.log.ErrorErr(err)
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
//		return
//	}
//
//	fmt.Printf("mode %s corporation %s Name %s\n", mode, i.Guild.Name, userName)
//
//	if mode == "get" {
//		sd := models.SyncData{
//			TechLevels: models.TechLevels{},
//			Ver:        2,
//			InSync:     1,
//		}
//		techLevels, err := s.db.TechnologiesGet(userUUID)
//		if err == nil && techLevels != nil {
//			sd.TechLevels = *techLevels
//		}
//		c.JSON(http.StatusOK, sd)
//	} else if mode == "sync" {
//
//		var data models.SyncData
//		if bindErr := c.BindJSON(&data); bindErr != nil {
//			fmt.Println(bindErr)
//			c.JSON(400, gin.H{"error": bindErr.Error()})
//			return
//		}
//		updateErr := s.db.TechnologiesUpdate(userUUID, userName, data.TechLevels)
//		if updateErr != nil {
//			s.log.ErrorErr(updateErr)
//		}
//
//		// Используйте переменную data с полученными данными
//		c.JSON(http.StatusOK, data)
//	}
//}

func getCorpsTypeId(mg *models.MultiAccountGuildV2) (guildDs, guildTg, guildWa []string) {
	for t, channel := range mg.Channels {
		switch t {
		case "wa":
			guildWa = append(guildWa, channel)
		case "ds":
			guildDs = append(guildDs, channel)
		case "tg":
			guildTg = append(guildTg, channel)
		}
	}
	return guildDs, guildTg, guildWa
}

func corpMemberDetectType(cm []models.CorpMember) []models.CorpMember {
	var members []models.CorpMember
	if len(cm) == 0 {
		return members
	}
	for _, m := range cm {
		m.TypeAccount = ""
		uId := m.UserId
		if strings.Contains(m.UserId, "/") {
			split := strings.Split(m.UserId, "/")
			uId = split[0]
		}
		if len(uId) < 12 {
			m.TypeAccount = "tg"
		} else if strings.Contains(uId, "@") {
			m.TypeAccount = "wa"
		} else if len(uId) > 12 && len(uId) < 24 {
			m.TypeAccount = "ds"
		}
		if m.TypeAccount == "" {
			fmt.Printf("TypeAccount len(%d) userid %s \n", len(m.UserId), m.UserId)
		}
		if m.TypeAccount != "" {
			m.Name = fmt.Sprintf("(%s) %s", strings.ToUpper(m.TypeAccount), m.Name)
			members = append(members, m)
		}
	}
	return members
}
