package server

import (
	"compendium/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (s *Server) RunServerRestApi() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.POST("/compendium/inbox", s.InboxMessage)
	router.GET("/compendium/api", s.api)
	router.GET("/compendium/api/user", s.apiUserAlts)

	err := router.Run(":80")
	if err != nil {
		s.log.ErrorErr(err)
	}
}

func (s *Server) InboxMessage(c *gin.Context) {
	var data models.IncomingMessage
	if err := c.BindJSON(&data); err != nil {
		s.log.ErrorErr(err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	s.In <- data
	c.JSON(http.StatusOK, "ok")
}

func (s *Server) api(c *gin.Context) {
	userid := c.Query("userid")
	guildid := c.Query("guildid")
	if userid != "" {
		fmt.Println("userid:", userid)
		multiAccount, _ := s.multi.FindMultiAccountByUserId(userid)
		if multiAccount != nil {
			fmt.Printf("multiAccount: %+v\n", multiAccount)
			corpMember, _ := s.multi.TechnologiesGetAllCorpMember(models.CorpMember{
				MultiAccount: multiAccount,
				UserId:       userid})
			if len(corpMember) != 0 {
				fmt.Printf("corpMember: %+v\n", corpMember)
				c.JSON(http.StatusOK, corpMember)
				return
			}

		}
	}
	if userid == "" || guildid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userid and guildid must not be empty"})
		return
	}
	read, err := s.db.CorpMembersApiRead(guildid, userid)
	if len(read) == 0 || read == nil || err != nil {

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

func (s *Server) apiUserAlts(c *gin.Context) {
	userid := c.Query("userid")
	if userid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userid must not be empty"})
		return
	}
	fmt.Println("userid:", userid)
	accountByUserId, _ := s.multi.FindMultiAccountByUserId(userid)
	if accountByUserId != nil && len(accountByUserId.Alts) > 0 {
		fmt.Printf("accountByUserId: %+v\n", accountByUserId.Alts)
		c.JSON(http.StatusOK, accountByUserId.Alts)
		return
	}

	read, _ := s.db.UsersGetByUserId(userid)
	if read == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userid not found"})
		return
	}
	if len(read.Alts) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "alts not found"})
		return
	}
	c.JSON(http.StatusOK, read.Alts)
}
