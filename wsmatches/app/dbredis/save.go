package dbredis

import (
	"fmt"
	"time"
	"ws/models"
)

func (r *Db) SaveCorpDate(key string, corp models.CorpsData) {

	// Set hash field-values
	err := r.c.HSet(ctx, key, map[string]interface{}{
		"Corp1Name":  corp.Corp1Name,
		"Corp2Name":  corp.Corp2Name,
		"Corp1Score": corp.Corp1Score,
		"Corp2Score": corp.Corp2Score,
		"DateEnded":  corp.DateEnded.Format(time.RFC3339), // Сохраняем время как строку
	}).Err()
	if err != nil {
		log.ErrorErr(err)
	}
	fmt.Printf("Saved %s Date %+v\n", key, corp)
}
