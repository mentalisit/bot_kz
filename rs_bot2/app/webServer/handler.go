package webServer

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) getMobileSetting(c *gin.Context) {
	// Получаем UUID из query параметра
	uuidStr := c.Query("uuid")
	if uuidStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "UUID parameter is required",
		})
		return
	}

	// Парсим UUID
	uid, err := uuid.Parse(uuidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid UUID format",
			"uuid":  uuidStr,
		})
		return
	}

	answer := make(map[string]interface{})
	answer["uuid"] = uuidStr
	answer["mode"] = nil

	// Получаем данные из базы данных
	data, err := s.db.GetOtherByUUID(uid)
	if err != nil || data == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Data not found",
			"uuid":  uuidStr,
		})
		return
	}

	//rs bot
	conf := s.db.ReadConfigV2Uid(uid.String())
	if conf != nil && conf.Uid != "" {
		answer["mode"] = "rs"
		conf.GetGameLevelInCorpInfo(s.db.ReadCorpInfoByCorpID)
		answer["rs_config"] = conf
	}
	if conf != nil && conf.Channels[data.Data.ChannelId] != nil {
		answer["data"] = *conf.Channels[data.Data.ChannelId]
	} else {
		answer["data"] = data.Data
	}

	//bridge
	if answer["mode"] == nil {
		bridge, exist := s.db.ReadBridgeConfigByChannelId(data.Data.ChannelId)
		if !exist {
			bridge, exist = s.db.ReadBridgeConfigByNameRelay(data.Uuid)
		}
		if exist {
			answer["mode"] = "bridge"
			answer["bridge_config"] = bridge
		}
	}

	//news
	if answer["mode"] == nil {
		configNews, exist := s.db.IsSubscribedToNews(data.Data.ChannelId)
		if exist {
			answer["mode"] = "news"
			answer["news_config"] = configNews
		}
	}

	//scoreboard
	if answer["mode"] == nil {
		scoreboardParamsV2 := s.db.ScoreboardReadByUid(uuidStr)
		if scoreboardParamsV2 != nil && len(scoreboardParamsV2.Channels) != 0 {
			answer["mode"] = "scoreboard"
			answer["scoreboard_param"] = scoreboardParamsV2
		}
	}

	if data.Data.TypeMessenger == "ds" {
		s.cl.Ds.DeleteMessage(data.Data.ChannelId, data.Data.MessageId)
	} else if data.Data.TypeMessenger == "tg" {
		s.cl.Tg.DeleteMessage(data.Data.ChannelId, data.Data.MessageId)
	} else {
		s.log.Info("Unknown typeMessenger " + data.Data.TypeMessenger)
	}
	if data.Read {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Already read",
			"uuid":  uuidStr,
		})
		return
	}

	// Возвращаем данные в формате JSON
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   answer,
	})
	s.db.UpdateReadOtherByUUID(uuidStr)
}

// HealthCheckHandler проверяет здоровье сервиса
func HealthCheckHandler(c *gin.Context) {
	// Если все проверки пройдены успешно
	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"message":   "Service is healthy",
		"timestamp": time.Now().Unix(),
		"cors":      "enabled",
	})
}

func (s *Server) getWebhook(c *gin.Context) {
	// Получаем UUID из query параметра
	uuidStr := c.Query("uuid")
	if uuidStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "UUID parameter is required",
		})
		return
	}

	// Парсим UUID
	uid, err := uuid.Parse(uuidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid UUID format",
			"uuid":  uuidStr,
		})
		return
	}

	corpId := c.Query("corpId")
	corpanme := s.getCorps(corpId)

	answer := make(map[string]interface{})
	answer["uuid"] = uuidStr

	// Получаем данные из базы данных
	data, err := s.db.GetOtherByUUID(uid)
	if err != nil || data == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Data not found",
			"uuid":  uuidStr,
		})
		return
	}

	corpInfo, _ := s.db.ReadCorpInfoByCorpID(corpId)
	if corpInfo != nil && corpInfo.Webhook {
		answer["webhook"] = "уже используется"
		answer["chat_id"] = "not required"
	} else {
		webhookUrl, chatId := s.cl.Ds.GetOrCreateWebhookGame(formatDiscordChannelName(corpanme))
		answer["webhook"] = webhookUrl
		answer["chat_id"] = chatId
	}

	// Возвращаем данные в формате JSON
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   answer,
	})
}

func (s *Server) getGameData(c *gin.Context) {
	// Получаем ID из query параметра
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID parameter is required",
		})
		return
	}

	corpInfo, _ := s.db.ReadCorpInfoByCorpID(id)
	if corpInfo != nil && corpInfo.CorpID == id {
		corpInfo.Level = corpInfo.GetLevelByXP()
		// Возвращаем данные в формате JSON
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   corpInfo,
		})
		return
	}

	// Возвращаем данные в формате JSON
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"id":     id,
		"data":   "not found",
	})
}

// FormatDiscordChannelName приводит строку к стандарту текстового канала Discord
func formatDiscordChannelName(input string) string {
	// 1. Переводим в нижний регистр
	name := strings.ToLower(input)

	// 2. Заменяем пробелы и подчеркивания на дефисы
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "_", "-")

	// 3. Удаляем всё, что не буквы, цифры или дефисы
	// Мы оставляем поддержку Unicode букв (чтобы работал русский язык)
	reg := regexp.MustCompile(`[^a-z0-9а-яё\-]+`)
	name = reg.ReplaceAllString(name, "")

	// 4. Убираем двойные дефисы, если они возникли (напр. "hello  world" -> "hello--world")
	regDoubleDash := regexp.MustCompile(`-+`)
	name = regDoubleDash.ReplaceAllString(name, "-")

	// 5. Убираем дефисы в начале и в конце
	name = strings.Trim(name, "-")

	return name
}

func (s *Server) getCorps(corpId string) string {
	file, err := os.ReadFile("docker/ws/corps.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return ""
	}
	type Corporation struct {
		Name string `json:"Name"`
		Id   string `json:"Id"`
		Win  int    `json:"Win"`
		Loss int    `json:"Loss"`
		Draw int    `json:"Draw"`
		Elo  int    `json:"Elo"`
	}

	var corps []Corporation
	err = json.Unmarshal(file, &corps)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return ""
	}

	for _, corporation := range corps {
		if corporation.Id == corpId {
			return corporation.Name
		}
	}

	return ""
}

func GetCorpBonusRemote(corpId string) (*CorpBonus, error) {
	url := "https://raw.githubusercontent.com/CapricanDRJ/HadesStarCommunityData/refs/heads/main/corpBonus.csv"

	// 1. Устанавливаем таймаут, чтобы запрос не висел вечно
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка сервера: %d", resp.StatusCode)
	}

	// 2. Используем csv.Reader для потокового чтения
	reader := csv.NewReader(resp.Body)

	// Пропускаем заголовок
	_, _ = reader.Read()

	for {
		record, err := reader.Read() // Читаем одну строку

		if err == io.EOF {
			break // Достигли конца файла, выходим из цикла
		}
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения строки: %v", err)
		}
		// 3. Проверяем ID (полный или короткий — зависит от твоей логики)
		// Если хочешь искать по последним 3 символам:
		// currentId := record[0]
		// if len(currentId) > 3 && currentId[len(currentId)-3:] == corpId { ... }
		fmt.Println(record)
		//if record[0] == corpId {
		//	bEnd, _ := strconv.ParseInt(record[1], 10, 64)
		//	pct, _ := strconv.Atoi(record[2])
		//
		//	return &CorpBonus{
		//		CorporationId: record[0],
		//		BonusEnd:      bEnd,
		//		Percent:       pct,
		//	}, nil
		//}
	}

	return nil, fmt.Errorf("корпорация %s не найдена", corpId)
}

type CorpBonus struct {
	CorporationId string `json:"corporation_id"`
	BonusEnd      int64  `json:"bonus_end"`
	Percent       int    `json:"percent"`
}

func (s *Server) getChatChannels(c *gin.Context) {
	gidStr := c.Query("gid")
	if gidStr == "" {
		gidStr = "00000000-0000-0000-0000-000000000000"
	}

	channels, err := s.db.GetChatChannels(gidStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get channels"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   channels,
	})
}
