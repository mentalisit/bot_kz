package tg

//const endpoint = "kz_bot"
//
//func (t *Telegram) MarshalDataSendBridgeAsync(message models.BridgeSendToMessenger) []models.MessageIds {
//	data, err := json.Marshal(message)
//	if err != nil {
//		t.log.ErrorErr(err)
//		return nil
//	}
//
//	resp, err := http.Post("http://"+endpoint+"/telegram/send/bridge", "application/json", bytes.NewBuffer(data))
//	if err != nil {
//		resp, err = http.Post("http://192.168.100.131:802/send/bridge", "application/json", bytes.NewBuffer(data))
//		if err != nil {
//			t.log.ErrorErr(err)
//			return nil
//		}
//	}
//	var dataReply []models.MessageIds
//	err = json.NewDecoder(resp.Body).Decode(&dataReply)
//	if err != nil {
//		t.log.Info(fmt.Sprintf("err resp.Body %+v\n", resp.Body))
//		t.log.ErrorErr(err)
//		return nil
//	}
//	return dataReply
//}
//
//func (t *Telegram) MarshalDelTelegram(message models.DeleteMessageStruct) {
//	data, err := json.Marshal(message)
//	if err != nil {
//		t.log.ErrorErr(err)
//		return
//	}
//
//	_, err = http.Post("http://"+endpoint+"/telegram/del", "application/json", bytes.NewBuffer(data))
//	if err != nil {
//		_, err = http.Post("http://192.168.100.131:802/data", "application/json", bytes.NewBuffer(data))
//		if err != nil {
//			t.log.ErrorErr(err)
//			return
//		}
//	}
//}
//func (t *Telegram) MarshalSendDelTelegram(message models.SendTextDeleteSeconds) {
//	data, err := json.Marshal(message)
//	if err != nil {
//		t.log.ErrorErr(err)
//		return
//	}
//
//	_, err = http.Post("http://"+endpoint+"/telegram/SendDel", "application/json", bytes.NewBuffer(data))
//	if err != nil {
//		_, err = http.Post("http://192.168.100.131:802/data", "application/json", bytes.NewBuffer(data))
//		if err != nil {
//			t.log.ErrorErr(err)
//			return
//		}
//	}
//}
