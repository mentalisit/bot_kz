package tg

//func (t *Telegram) SendBridgeFunc(chatid []string, text string, extra []models.FileInfo, resultChannel chan<- models.MessageIds, wg *sync.WaitGroup) {
//
//	defer wg.Done()
//
//	m := models.BridgeSendToMessenger{
//		Text:      text,
//		ChannelId: chatid,
//		Extra:     extra,
//	}
//	MessageIds := t.MarshalDataSendBridgeAsync(m)
//
//	for _, id := range MessageIds {
//		resultChannel <- id
//	}
//
//}
