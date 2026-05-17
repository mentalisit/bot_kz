package webServer

import (
	"net/http"
	"rs/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) getStudy(c *gin.Context) {
	uidString := c.Query("uuid")
	uid, err := uuid.Parse(uidString)
	name := c.Query("name")
	if uidString == "" || name == "" || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing uuid"})
		return
	}

	study, err := s.db.GetStudy(uid, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, study)
}

func (s *Server) postStudy(c *gin.Context) {
	var req models.Study

	// Привязываем JSON к структуре
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "Некорректный формат данных"})
		return
	}

	var err error
	if len(req.Studies) == 0 {
		err = s.db.DeleteStudyRecord(req)
	} else {
		err = s.db.InsertStudy(req)
	}
	if err != nil {
		s.log.ErrorErr(err)
		c.JSON(400, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(200, "ok")
}
