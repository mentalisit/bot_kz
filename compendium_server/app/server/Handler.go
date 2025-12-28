package server

import (
	"compendium_s/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *Server) CheckIdentityHandler(c *gin.Context) {
	code := c.GetHeader("authorization")

	// Проверка наличия кода в запросе и его длины
	if code == "" || len(code) != 14 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid code"})
		return
	}

	identity := s.CheckCode(code)

	if identity.Token != "" {
		c.JSON(http.StatusOK, identity)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid code"})
}

func (s *Server) CheckConnectHandler(c *gin.Context) {
	token := c.GetHeader("authorization")

	// Локальная проверка
	i := s.GetTokenIdentity(token)
	if i != nil {
		if i.User.GameName != "" {
			i.User.Username = i.User.GameName
		}
		c.JSON(http.StatusOK, i)
		return
	}
	c.JSON(http.StatusForbidden, nil)
}

func (s *Server) CheckCorpDataHandler(c *gin.Context) {
	token, roleId, mGuild, err := extractAndValidateCheckHeaders(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		s.log.ErrorErr(err)
		return
	}
	cacheKey := token + ":" + roleId + ":" + mGuild

	if data, ok := s.getFreshCache(cacheKey); ok {
		c.JSON(http.StatusOK, data)
		return
	}
	result, status, err := s.fetchCorpData(token, roleId, mGuild)
	if err != nil {
		s.log.ErrorErr(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	s.setCache(cacheKey, result)
	c.JSON(http.StatusOK, result)
}

func (s *Server) CheckRefreshHandler(c *gin.Context) {
	token := c.GetHeader("authorization")
	if token == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "missing token"})
		return
	}

	i := s.refreshToken2(token)
	if i != nil {
		if i.User.GameName != "" {
			i.User.Username = i.User.GameName
		}
		c.JSON(http.StatusOK, i)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
}

func (s *Server) CheckSyncTechHandler(c *gin.Context) {
	mode := c.Param("mode")
	twin := c.DefaultQuery("twin", "")
	token := c.GetHeader("authorization")
	if strings.HasPrefix(token, "my_compendium_") {
		i2 := s.GetTokenIdentity(token)
		if i2 != nil && i2.MAccount.Nickname != "" {
			if mode == "" {
				mode = c.GetHeader("X-Sync-Mode")
			}
			if twin == "" {
				twin = c.GetHeader("X-Alt-Name")
			}

			userName := i2.MAccount.Nickname
			if twin != "" && twin != "default" {
				userName = twin
			}

			sd := models.SyncData{
				TechLevels: models.TechLevels{},
				Ver:        2,
				InSync:     1,
			}

			if mode == "get" {
				techLevels, err := s.dbV2.TechnologiesGet(i2.MAccount.UUID, userName)
				if err == nil && techLevels != nil {
					sd.TechLevels = *techLevels
				}
				c.JSON(http.StatusOK, sd)
			} else if mode == "sync" {
				if bindErr := c.BindJSON(&sd); bindErr != nil {
					fmt.Println(bindErr)
					c.JSON(400, gin.H{"error": bindErr.Error()})
					return
				}
				updateErr := s.dbV2.TechnologiesUpdate(i2.MAccount.UUID, userName, sd.TechLevels)
				if updateErr != nil {
					s.log.ErrorErr(updateErr)
				}
				c.JSON(http.StatusOK, sd)
			}

			return
		}
	}

	i := s.GetTokenIdentity(token)
	if i == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}
	if i.MultiAccount != nil && i.MGuild != nil {
		s.SyncTechMulti(c, i, mode, twin)
		return
	}
	userId := i.User.ID
	userName := i.User.Username
	//guildId := i.MultiGuild.GId.String()
	if twin != "" && twin != "default" {
		userName = twin
	}
	if userName == i.User.GameName {
		userName = i.User.Username
	}

	fmt.Printf("mode %s corporation %s Name %s\n", mode, i.Guild.Name, userName)

	if mode == "get" {
		sd := models.SyncData{
			TechLevels: models.TechLevels{},
			Ver:        1,
			InSync:     1,
		}
		techBytes, err := s.db.TechGet(userName, userId)
		if err == nil && len(techBytes) > 0 {
			sd.TechLevels = sd.TechLevels.ConvertToTech(techBytes)
		}
		c.JSON(http.StatusOK, sd)
	} else if mode == "sync" {

		var data models.SyncData
		if err := c.BindJSON(&data); err != nil {
			fmt.Println(err)
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		//read old tech
		techBytes, _ := s.db.TechGet(userName, userId)
		m := make(map[int]models.TechLevel)
		err := json.Unmarshal(techBytes, &m)
		//comparison of data by date
		for id, l := range data.TechLevels {
			if m[id].Ts > l.Ts {
				//ignore data
				text := fmt.Sprintf("%s %s module %d data old %d new %d Ok>\n",
					i.User.Username, i.Guild.Name, id, m[id].Ts, l.Ts)
				s.log.Info(text)
				data.TechLevels[id] = m[id]
			}
		}

		bytes, err := json.Marshal(data.TechLevels)
		if err != nil {
			s.log.ErrorErr(err)
		}
		err = s.db.TechUpdate(userName, userId, bytes)
		if err != nil {
			s.log.ErrorErr(err)
		}

		// Используйте переменную data с полученными данными
		c.JSON(http.StatusOK, data)
	}
}

func (s *Server) CheckUserCorporationsHandler(c *gin.Context) {
	token := c.GetHeader("authorization")
	if token == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "missing token"})
		return
	}

	// Получаем информацию о пользователе по токену
	i := s.GetTokenIdentity(token)
	if strings.HasPrefix(token, "my_compendium_") {
		if i != nil && i.MAccount.Nickname != "" {
			// Получаем список корпораций пользователя
			corporations, err := s.dbV2.UserCorporationsGet(i)
			if err != nil {
				s.log.ErrorErr(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user corporations"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"user":         i.MAccount,
				"corporations": corporations,
			})
			return
		}
	}

	if i != nil {
		// Получаем список корпораций пользователя
		corporations, err := s.db.UserCorporationsGet(i)
		if err != nil {
			s.log.ErrorErr(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user corporations"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user":         i.User,
			"corporations": corporations,
		})
		return
	}

	c.JSON(http.StatusForbidden, gin.H{"error": "invalid token"})
}

func (s *Server) links(c *gin.Context) {
	htmlContent := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Links</title>
    <style>
        html, body {
            height: 100%;
            margin: 0;
            padding: 0;
            display: flex;
            justify-content: center;
            align-items: center;
            font-family: Arial, sans-serif;
            background-color: #f4f4f4; // Светлый фон всей страницы
        }
        .centered-content {
            width: 100%; // Ширина контента равна ширине страницы
            max-width: 600px; // Максимальная ширина контента
            text-align: center;
            box-shadow: 0 4px 8px rgba(0,0,0,0.1);
            background-color: white;
            padding: 20px;
            border-radius: 8px;
        }
        ul {
            list-style-type: none;
            padding: 0;
        }
        li {
            margin: 10px 0;
        }
        a {
            text-decoration: none;
            color: #3366cc;
        }
        a:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>
    <div class="centered-content">
        <h1>Select link:</h1>
        <ul>
            <li><a href="https://discord.com/oauth2/authorize?client_id=909526127305953290&scope=bot+applications.commands&permissions=141533113424" target="_blank">Invite Discord Bot</a></li>
            <li><a href="https://t.me/gote1st_bot" target="_blank">Invite Telegram Bot</a></li>
            <li><a href="https://discord.com/users/582882137842122773" target="_blank">Message Discord Bot Author</a></li>
            <li><a href="https://t.me/mentalisit" target="_blank">Message Telegram Bot Author</a></li>
        </ul>
    </div>
</body>
</html>
`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlContent))
}

func (s *Server) api(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "empty token"})
		return
	}
	i := s.GetTokenIdentity(token)
	if i == nil || i.Token == "" || i.Guild.ID != "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid token"})
		return
	}
	userid := c.Query("userid")
	if userid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userid must not be empty"})
		return
	}

	read, _ := s.db.CorpMembersRead(i.Guild.ID)
	if len(read) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "guildid empty members"})
		return
	}
	var memb []models.CorpMember
	for _, member := range read {
		if strings.Contains(member.UserId, userid) {
			memb = append(memb, member)
		}
	}
	if len(memb) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "member not found"})
		return
	}
	c.JSON(http.StatusOK, memb)
}
