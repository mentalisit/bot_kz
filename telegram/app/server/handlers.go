package server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"telegram/models"
	"time"
)

const (
	funcSend                = "send"
	funcDeleteMessage       = "delete_message"
	funcEditMessage         = "edit_message"
	funcEditMessageTextKey  = "edit_message_text_key"
	funcSendDel             = "send_del"
	funcDeleteMessageSecond = "delete_message_second"
	funcCheckAdmin          = "check_admin"
	funcChatTyping          = "chat_typing"
	funcSendHelp            = "send_help"
	funcSendEmbed           = "send_embed"
	funcSendEmbedTime       = "send_embed_time"
	funcGetAvatarUrl        = "get_avatar_url"
)

type apiRs struct {
	FuncApi   string   `json:"funcApi"`
	Text      string   `json:"text"`
	Channel   string   `json:"channel"`
	MessageId string   `json:"messageId"`
	ParseMode string   `json:"parseMode"`
	Second    int      `json:"second"`
	LevelRs   string   `json:"levelRs"`
	Levels    []string `json:"levels"`
	UserName  string   `json:"userName"`
}

func (s *Server) funcInbox(c *gin.Context) {
	s.PrintGoroutine()
	var rawData json.RawMessage

	// Читаем тело запроса как необработанные JSON-данные
	if err := c.ShouldBindJSON(&rawData); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("funcInbox", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Пробуем десериализовать
	var m apiRs
	err := json.Unmarshal(rawData, &m)
	if err == nil {
		code, q := s.selectFunc(m)
		c.JSON(code, q)
		return
	}

	if len(rawData) > 4 {
		// просто выводим полученные данные как есть
		var otherData interface{}
		if err := json.Unmarshal(rawData, &otherData); err == nil {
			s.log.InfoStruct("	Received other data", otherData)
			c.JSON(http.StatusOK, "OK")
			return
		}
	}
	fmt.Println(len(rawData))
	fmt.Println(string(rawData))
	// Если ничего не удалось разобрать, выводим ошибку
	s.log.Error("Не удалось разобрать данные")
	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
}

type answer struct {
	ArrString string    `json:"arrString"`
	ArrInt    int       `json:"arrInt"`
	ArrBool   bool      `json:"arrBool"`
	ArrError  error     `json:"arrError"`
	Time      time.Time `json:"time"`
}

func (s *Server) selectFunc(m apiRs) (code int, q answer) {
	fmt.Printf("selectFunc %+v\n", m)
	q.Time = time.Now()
	switch m.FuncApi {
	case funcSend:
		fmt.Printf("channel %s text %s\n", m.Channel, m.Text)
		q.ArrString, q.ArrError = s.tg.SendChannel(m.Channel, m.Text)
		q.ArrInt, _ = strconv.Atoi(q.ArrString)
		if q.ArrError != nil {
			return http.StatusForbidden, q
		}
		fmt.Printf("answer %+v\n", q)
	case funcDeleteMessage:
		fmt.Printf("channel %s mid %s\n", m.Channel, m.MessageId)
		id, _ := strconv.Atoi(m.MessageId)
		q.ArrBool = s.tg.DelMessage(m.Channel, id)
		fmt.Printf("answer %+v\n", q)
	case funcEditMessage:
		fmt.Printf("channel %s text %s mid %s parse %s\n", m.Channel, m.Text, m.MessageId, m.ParseMode)
		mID, _ := strconv.Atoi(m.MessageId)
		q.ArrError = s.tg.EditTextParseMode(m.Channel, mID, m.Text, m.ParseMode)
		fmt.Printf("answer %+v\n", q)
	case funcSendDel:
		fmt.Printf("channel %s text %s second %d\n", m.Channel, m.Text, m.Second)
		q.ArrBool = s.tg.SendChannelDelSecond(m.Channel, m.Text, m.Second)
		fmt.Printf("answer %+v\n", q)
	case funcEditMessageTextKey:
		fmt.Printf("channel %s text %s mid %s levelrs %s\n", m.Channel, m.Text, m.MessageId, m.LevelRs)
		mID, _ := strconv.Atoi(m.MessageId)
		s.tg.EditMessageTextKey(m.Channel, mID, m.Text, m.LevelRs)
		q.ArrBool = true
	case funcDeleteMessageSecond:
		fmt.Printf("channel %s mid %s second %d\n", m.Channel, m.MessageId, m.Second)
		s.tg.DelMessageSecond(m.Channel, m.MessageId, m.Second)
		q.ArrBool = true
		fmt.Printf("answer %+v\n", q)
	case funcCheckAdmin:
		q.ArrBool = s.tg.CheckAdminTg(m.Channel, m.UserName)
	case funcChatTyping:
		s.tg.ChatTyping(m.Channel)
		q.ArrBool = true
	case funcSendHelp:
		q.ArrString = s.tg.SendHelp(m.Channel, m.Text, m.MessageId)
	case funcSendEmbed:
		q.ArrInt = s.tg.SendEmbed(m.LevelRs, m.Channel, m.Text)
		fmt.Printf("answer embed %+v\n", q)
	case funcSendEmbedTime:
		q.ArrInt = s.tg.SendEmbedTime(m.Channel, m.Text)
	case funcGetAvatarUrl:
		q.ArrString = s.tg.GetAvatarUrl(m.UserName)

	default:
		{
			q.ArrString = "case not found"
			return http.StatusBadRequest, q
		}
	}
	return http.StatusOK, q
}
func (s *Server) telegramSendPic(c *gin.Context) {
	s.PrintGoroutine()
	var m models.SendPic
	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("telegramSendPic", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := s.tg.SendPic(m.Channel, m.Text, m.Pic)
	if err != nil {
		s.log.ErrorErr(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Message sent to Telegram successfully"})
}
func (s *Server) telegramSendBridge(c *gin.Context) {
	s.PrintGoroutine()
	var m models.BridgeSendToMessenger
	if err := c.ShouldBindJSON(&m); err != nil {
		s.log.ErrorErr(err)
		s.log.InfoStruct("telegramSendBridge", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	messageTg := s.tg.SendBridgeFuncRest(m)
	c.JSON(http.StatusOK, messageTg)
}
func (s *Server) telegramGetAvatarUrl(c *gin.Context) {
	s.PrintGoroutine()
	userid := c.Query("userid")
	if userid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userid must not be empty"})
		return
	}

	urlAvatar := s.tg.GetAvatarUrl(userid)
	c.JSON(http.StatusOK, urlAvatar)
}
