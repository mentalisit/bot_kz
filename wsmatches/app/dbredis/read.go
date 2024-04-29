package dbredis

import (
	"strconv"
	"time"
	"ws/models"
)

func (r *Db) ReadCorpData(key string) *models.CorpsData {
	// Get hash field-values
	userInfo, err := r.c.HGetAll(ctx, key).Result()
	if err != nil {
		log.ErrorErr(err)
	}
	if len(userInfo) == 0 {
		return nil
	}

	Corp1Score := userInfo["Corp1Score"]
	Corp2Score := userInfo["Corp2Score"]
	s1, _ := strconv.Atoi(Corp1Score)
	s2, _ := strconv.Atoi(Corp2Score)

	dateEnded, _ := time.Parse(time.RFC3339, userInfo["DateEnded"]) // Преобразование строки в time.Time

	corp1 := models.CorpsData{
		Corp1Name:  userInfo["Corp1Name"],
		Corp2Name:  userInfo["Corp2Name"],
		Corp1Score: s1,
		Corp2Score: s2,
		DateEnded:  dateEnded,
	}
	//fmt.Printf("read %s %+v\n", key, corp1)
	return &corp1
}
