package reststorage

//
//func (d Db) TimerDeleteMessage() []models.Timer {
//	var br []models.Timer
//	resp, err := http.Get("http://storage/storage/timer/delete")
//	if err != nil {
//		resp, err = http.Get("http://192.168.100.131:804/storage/timer/delete")
//		if err != nil {
//			d.log.ErrorErr(err)
//			return nil
//		}
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusOK {
//		d.log.ErrorErr(errors.New("resp.StatusCode != http.StatusOK"))
//		return nil
//	}
//
//	err = json.NewDecoder(resp.Body).Decode(&br)
//	if err != nil {
//		d.log.ErrorErr(err)
//		return nil
//	}
//	return br
//}
//
//func (d Db) TimerInsert(c models.Timer) {
//	data, err := json.Marshal(c)
//	if err != nil {
//		d.log.ErrorErr(err)
//		return
//	}
//
//	_, err = http.Post("http://storage/storage/timer/insert", "application/json", bytes.NewBuffer(data))
//	if err != nil {
//		_, err = http.Post("http://192.168.100.131:804/storage/timer/insert", "application/json", bytes.NewBuffer(data))
//		if err != nil {
//			d.log.ErrorErr(err)
//			return
//		}
//	}
//}
