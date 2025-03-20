package dbredis

//
//func (r *Db) ReadCorpData(key string) *models.CorporationsData {
//	// Get hash field-values
//	userInfo, err := r.c.HGetAll(ctx, key).Result()
//	if err != nil {
//		log.ErrorErr(err)
//	}
//	if len(userInfo) == 0 {
//		return nil
//	}
//
//	s1, _ := strconv.Atoi(userInfo["Corporation1Score"])
//	s2, _ := strconv.Atoi(userInfo["Corporation2Score"])
//
//	dateEnded, _ := time.Parse(time.RFC3339, userInfo["DateEnded"]) // Преобразование строки в time.Time
//
//	corp1 := models.CorporationsData{
//		Corporation1Name:  userInfo["Corporation1Name"],
//		Corporation2Name:  userInfo["Corporation2Name"],
//		Corporation1Score: s1,
//		Corporation2Score: s2,
//		Corporation1Id:    userInfo["Corporation1Id"],
//		Corporation2Id:    userInfo["Corporation2Id"],
//		DateEnded:         dateEnded,
//	}
//	//fmt.Printf("read %s %+v\n", key, corp1)
//	return &corp1
//}
