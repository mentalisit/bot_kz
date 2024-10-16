package server

//func (s *Server) telegramSendBridge(c *gin.Context) {
//	var m models.BridgeSendToMessenger
//	if err := c.ShouldBindJSON(&m); err != nil {
//		s.log.ErrorErr(err)
//		s.log.InfoStruct("telegramSendBridge", c.Request.Body)
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	messageTg := s.cl.Tg.SendBridgeFuncRest(m)
//	c.JSON(http.StatusOK, messageTg)
//}
//
//func (s *Server) telegramSendText(c *gin.Context) {
//	parseMode := c.DefaultQuery("parse", "")
//	var m models.SendText
//	if err := c.ShouldBindJSON(&m); err != nil {
//		s.log.ErrorErr(err)
//		s.log.InfoStruct("telegramSendText", c.Request.Body)
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	mid, err := s.cl.Tg.Send(m.Channel, m.Text, parseMode)
//	if err != nil {
//		s.log.ErrorErr(err)
//		if err.Error() == "Forbidden: bot can't initiate conversation with a user" {
//			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
//			return
//		} else {
//			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//			return
//		}
//	}
//	c.JSON(http.StatusOK, mid)
//}
//
//func (s *Server) telegramEditMessage(c *gin.Context) {
//	parseMode := c.DefaultQuery("parse", "")
//	var m models.EditText
//	if err := c.ShouldBindJSON(&m); err != nil {
//		s.log.ErrorErr(err)
//		s.log.InfoStruct("telegramEditMessage", c.Request.Body)
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	mid, err := strconv.Atoi(m.MessageId)
//	if err != nil {
//		return
//	}
//	err = s.cl.Tg.EditTextParseMode(m.Channel, mid, m.Text, parseMode)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, err)
//		return
//	}
//	c.JSON(http.StatusOK, "ok")
//}
//
//func (s *Server) telegramSendPic(c *gin.Context) {
//	var m models.SendPic
//	if err := c.ShouldBindJSON(&m); err != nil {
//		s.log.ErrorErr(err)
//		s.log.InfoStruct("telegramSendPic", c.Request.Body)
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	err := s.cl.Tg.SendPic(m.Channel, m.Text, m.Pic)
//	if err != nil {
//		s.log.ErrorErr(err)
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{"status": "Message sent to Telegram successfully"})
//}
//func (s *Server) telegramDel(c *gin.Context) {
//	var m models.DeleteMessageStruct
//	if err := c.ShouldBindJSON(&m); err != nil {
//		s.log.ErrorErr(err)
//		s.log.InfoStruct("telegramDel", c.Request.Body)
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	id, err := strconv.Atoi(m.MessageId)
//	if err != nil {
//		s.log.ErrorErr(err)
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	s.cl.Tg.DelMessage(m.Channel, id)
//	c.JSON(http.StatusOK, gin.H{"status": "Message sent to telegram successfully"})
//}
//
//func (s *Server) telegramSendDelSecond(c *gin.Context) {
//	var m models.SendTextDeleteSeconds
//	if err := c.ShouldBindJSON(&m); err != nil {
//		s.log.ErrorErr(err)
//		s.log.InfoStruct("telegramSendDelSec", c.Request.Body)
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	s.cl.Tg.SendChannelDelSecond(m.Channel, m.Text, m.Seconds)
//	c.JSON(http.StatusOK, gin.H{"status": "Message sent to Discord successfully"})
//}
//
//func (s *Server) telegramGetAvatarUrl(c *gin.Context) {
//	userid := c.Query("userid")
//	if userid == "" {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "userid must not be empty"})
//		return
//	}
//	parseInt, err := strconv.ParseInt(userid, 10, 64)
//	if err != nil {
//		s.log.ErrorErr(err)
//		return
//	}
//
//	urlAvatar := s.cl.Tg.GetAvatarUrl(parseInt)
//	c.JSON(http.StatusOK, urlAvatar)
//	s.log.Info("telegramGetAvatarUrl")
//	c.JSON(http.StatusOK, gin.H{})
//}
