package dbredis

//
//func (r *Db) SaveCorpDate(key string, corp models.CorporationsData) {
//
//	// Set hash field-values
//	err := r.c.HSet(ctx, key, map[string]interface{}{
//		"Corporation1Name":  corp.Corporation1Name,
//		"Corporation2Name":  corp.Corporation2Name,
//		"Corporation1Score": corp.Corporation1Score,
//		"Corporation2Score": corp.Corporation2Score,
//		"Corporation1Id":    corp.Corporation1Id,
//		"Corporation2Id":    corp.Corporation2Id,
//		"DateEnded":         corp.DateEnded.Format(time.RFC3339), // Сохраняем время как строку
//	}).Err()
//	if err != nil {
//		log.ErrorErr(err)
//	}
//	fmt.Printf("Saved %s Date %+v\n", key, corp)
//}
