package server

import (
	"discord/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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
		q.ArrString = s.ds.Send(m.Channel, m.Text)
		fmt.Printf("answer %+v\n", q)
	case funcDeleteMessage:
		fmt.Printf("channel %s mid %s\n", m.Channel, m.MessageId)
		q.ArrBool = true
		s.ds.DeleteMessage(m.Channel, m.MessageId)
		fmt.Printf("answer %+v\n", q)
	case funcSendDel:
		fmt.Printf("channel %s text %s second %d\n", m.Channel, m.Text, m.Second)
		q.ArrBool = true
		s.ds.SendChannelDelSecond(m.Channel, m.Text, m.Second)
		fmt.Printf("answer %+v\n", q)
	case funcDeleteMessageSecond:
		fmt.Printf("channel %s mid %s second %d\n", m.Channel, m.MessageId, m.Second)
		s.ds.DeleteMesageSecond(m.Channel, m.MessageId, m.Second)
		q.ArrBool = true
		fmt.Printf("answer %+v\n", q)
	case funcCheckAdmin:
		//q.ArrBool = s.ds.CheckAdmin(m.UserId, m.Channel)
	case funcChatTyping:
		s.ds.ChannelTyping(m.Channel)
		q.ArrBool = true
	case funcSendHelp:
		//q.ArrString = s.ds.SendHelp(m.Channel, m.ParseMode, m.Text, m.Levels)
	case funcSendEmbed:
		//q.ArrInt = s.ds.SendEmbed(m.LevelRs, m.Channel, m.Text)
		fmt.Printf("answer embed %+v\n", q)
	case funcSendEmbedTime:
		q.ArrString = s.ds.SendEmbedTime(m.Channel, m.Text)
	case funcGetAvatarUrl:
		q.ArrString = s.ds.GetAvatarUrl(m.UserName)

	default:
		{
			q.ArrString = "case not found"
			return http.StatusBadRequest, q
		}
	}
	return http.StatusOK, q
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
	messageTg := s.ds.SendBridgeFuncRest(m)
	c.JSON(http.StatusOK, messageTg)
}
