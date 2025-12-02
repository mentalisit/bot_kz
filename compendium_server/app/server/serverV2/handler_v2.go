package serverV2

import (
	"compendium_s/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Frontend (JavaScript/TypeScript)
//headers: {
//'X-Sync-Mode': 'get',
//'X-Alt-Name': 'character123'
//}

//func (s *ServerV2) CheckIdentityHandler(c *gin.Context) {
//	c.Header("Access-Control-Allow-Origin", "*")
//	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
//	c.Header("Access-Control-Allow-Headers", "Authorization")
//	code := c.GetHeader("authorization")
//
//	// Проверка наличия кода в запросе и его длины
//	if code == "" || len(code) != 14 {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid code"})
//		return
//	}
//
//	identity := s.CheckCode(code)
//
//	if identity.Token != "" {
//		c.JSON(http.StatusOK, identity)
//		return
//	}
//	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid code"})
//}

func (s *ServerV2) CheckConnectHandler(c *gin.Context) {
	token := c.GetHeader("authorization")

	// Локальная проверка
	i := s.GetTokenIdentity(token)
	if i != nil {
		c.JSON(http.StatusOK, i)
		return
	}
	c.JSON(http.StatusForbidden, nil)
}

func (s *ServerV2) CheckSyncTechHandler(c *gin.Context) {
	mode := c.GetHeader("X-Sync-Mode")
	altName := c.GetHeader("X-Alt-Name")
	token := c.GetHeader("authorization")

	i := s.GetTokenIdentity(token)

	if i == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}

	userName := i.MultiAccount.Nickname
	if altName != "" && altName != "default" {
		userName = altName
	}

	tl := models.TechLevels{}

	if mode == "get" {
		techLevels, err := s.db.TechnologiesGet(i.MultiAccount.UUID, userName)
		if err == nil && techLevels != nil {
			tl = *techLevels
		}
		c.JSON(http.StatusOK, tl)
	} else if mode == "sync" {
		if bindErr := c.BindJSON(&tl); bindErr != nil {
			fmt.Println(bindErr)
			c.JSON(400, gin.H{"error": bindErr.Error()})
			return
		}
		updateErr := s.db.TechnologiesUpdate(i.MultiAccount.UUID, userName, tl)
		if updateErr != nil {
			s.log.ErrorErr(updateErr)
		}
		c.JSON(http.StatusOK, tl)
	}
}

func (s *ServerV2) CheckCorpDataHandler(c *gin.Context) {
	token := c.GetHeader("authorization")
	roleId := c.GetHeader("X-Role-ID")
	mGuild := c.GetHeader("X-Corp-ID")

	if token == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "missing token"})
		return
	}

	cacheKey := token + ":" + roleId + ":" + mGuild

	// Проверяем кэш
	s.cacheMutex.Lock()
	entry, exists := s.cacheReq[cacheKey]
	if exists && time.Since(entry.timestamp) <= 2*time.Second {
		// Есть свежий кэш, возвращаем
		s.cacheMutex.Unlock()
		c.JSON(http.StatusOK, entry.data)
		return
	}
	s.cacheMutex.Unlock()

	// Кэш отсутствует или устарел — получаем заново
	i := s.GetTokenIdentity(token)
	if i != nil {
		var result *models.CorpDataV2
		if mGuild != "" {
			// Получаем данные по конкретному ID корпорации
			var guild *models.MultiAccountGuildV2

			if gid, err := uuid.Parse(mGuild); err == nil {
				multiGuild, err := s.db.GuildGet(&gid)
				if err == nil && multiGuild != nil {
					guild = multiGuild
				}
			}

			if guild == nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "guild not found"})
				return
			}

			// Создаем временный Identity с выбранной гильдией
			tempIdentity := &models.IdentityV2{
				Token:        i.Token,
				MultiAccount: i.MultiAccount,
				GuildId:      mGuild,
			}

			result = s.GetCorpDataInternal(tempIdentity, roleId)
		}

		// Сохраняем в кэш
		s.cacheMutex.Lock()
		s.cacheReq[cacheKey] = cacheEntry{
			data:      result,
			timestamp: time.Now(),
		}
		s.cacheMutex.Unlock()

		c.JSON(http.StatusOK, result)
		return
	}

	c.JSON(http.StatusForbidden, gin.H{"error": "invalid token"})
}

func (s *ServerV2) CheckRefreshHandler(c *gin.Context) {
	token := c.GetHeader("authorization")
	if token == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "missing token"})
		return
	}

	token = s.refreshToken(token)

	i := s.GetTokenIdentity(token)
	if i != nil {
		c.JSON(http.StatusOK, i)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
}

//func (s *ServerV2) links(c *gin.Context) {
//	htmlContent := `
//<!DOCTYPE html>
//<html lang="en">
//<head>
//    <meta charset="UTF-8">
//    <meta name="viewport" content="width=device-width, initial-scale=1.0">
//    <title>Links</title>
//    <style>
//        html, body {
//            height: 100%;
//            margin: 0;
//            padding: 0;
//            display: flex;
//            justify-content: center;
//            align-items: center;
//            font-family: Arial, sans-serif;
//            background-color: #f4f4f4; // Светлый фон всей страницы
//        }
//        .centered-content {
//            width: 100%; // Ширина контента равна ширине страницы
//            max-width: 600px; // Максимальная ширина контента
//            text-align: center;
//            box-shadow: 0 4px 8px rgba(0,0,0,0.1);
//            background-color: white;
//            padding: 20px;
//            border-radius: 8px;
//        }
//        ul {
//            list-style-type: none;
//            padding: 0;
//        }
//        li {
//            margin: 10px 0;
//        }
//        a {
//            text-decoration: none;
//            color: #3366cc;
//        }
//        a:hover {
//            text-decoration: underline;
//        }
//    </style>
//</head>
//<body>
//    <div class="centered-content">
//        <h1>Select link:</h1>
//        <ul>
//            <li><a href="https://discord.com/oauth2/authorize?client_id=909526127305953290&scope=bot+applications.commands&permissions=141533113424" target="_blank">Invite Discord Bot</a></li>
//            <li><a href="https://t.me/gote1st_bot" target="_blank">Invite Telegram Bot</a></li>
//            <li><a href="https://discord.com/users/582882137842122773" target="_blank">Message Discord Bot Author</a></li>
//            <li><a href="https://t.me/mentalisit" target="_blank">Message Telegram Bot Author</a></li>
//        </ul>
//    </div>
//</body>
//</html>
//`
//	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlContent))
//}

//func (s *ServerV2) api(c *gin.Context) {
//	token := c.Query("token")
//	if token == "" {
//		c.JSON(http.StatusForbidden, gin.H{"error": "empty token"})
//		return
//	}
//	i := s.GetTokenIdentity(token)
//	if i == nil || i.Token == "" {
//		c.JSON(http.StatusForbidden, gin.H{"error": "invalid token"})
//		return
//	}
//	userid := c.Query("userid")
//	if userid == "" {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "userid must not be empty"})
//		return
//	}
//
//	read, _ := s.db.CorpMembersReadMulti(i.GetGuildUUID())
//	if len(read) == 0 {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "guildid empty members"})
//		return
//	}
//	var memb []models.CorpMember
//	for _, member := range read {
//		if strings.Contains(member.UserId, userid) {
//			memb = append(memb, member)
//		}
//	}
//	if len(memb) == 0 {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "member not found"})
//		return
//	}
//	c.JSON(http.StatusOK, memb)
//}

func (s *ServerV2) CheckUserCorporationsHandler(c *gin.Context) {
	token := c.GetHeader("authorization")
	if token == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "missing token"})
		return
	}

	// Получаем информацию о пользователе по токену
	i := s.GetTokenIdentity(token)
	if i == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid token"})
		return
	}
	fmt.Printf("		corp User %+v\n", i.MultiAccount)
	// Получаем список корпораций пользователя
	corporations, err := s.db.UserCorporationsGet(i)
	if err != nil {
		s.log.ErrorErr(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user corporations"})
		return
	}
	fmt.Printf("		corp Corps %+v\n", corporations)

	c.JSON(http.StatusOK, gin.H{
		"user":         i.MultiAccount,
		"corporations": corporations,
	})
}
