package tg

//const apinametg = "kz_bot"
//
//// отправка сообщения в телегу
//func SendOLD(chatId string, text string) (string, error) {
//	m := models.SendText{
//		Text:    text,
//		Channel: chatId,
//	}
//
//	data, err := json.Marshal(m)
//	if err != nil {
//		return "", err
//	}
//
//	resp, err := http.Post("http://"+apinametg+"/telegram/SendText", "application/json", bytes.NewBuffer(data))
//	if err != nil {
//		//_, err = http.Post("http://192.168.100.155:802/data", "application/json", bytes.NewBuffer(data))
//		return "", err
//	}
//	if resp.StatusCode == http.StatusForbidden {
//		return "", errors.New("forbidden")
//	}
//	var mid string
//	err = json.NewDecoder(resp.Body).Decode(&mid)
//	if err != nil {
//		return "", err
//	}
//	return mid, err
//}
//
//func SendPicOLD(channelId string, text string, pic []byte) error {
//	m := models.SendPic{
//		Text:    text,
//		Channel: channelId,
//		Pic:     pic,
//	}
//	data, err := json.Marshal(m)
//	if err != nil {
//		return err
//	}
//
//	resp, err := http.Post("http://"+apinametg+"/telegram/SendPic", "application/json", bytes.NewBuffer(data))
//	if err != nil {
//		fmt.Println(err)
//		//resp, err = http.Post("http://192.168.100.155:802/data", "application/json", bytes.NewBuffer(data))
//		return err
//	}
//	fmt.Println("resp.Status", resp.Status)
//	return nil
//}
//func DeleteMessageOLD(ChatId string, MesId string) error {
//	s := models.DeleteMessageStruct{
//		MessageId: MesId,
//		Channel:   ChatId,
//	}
//	data, err := json.Marshal(s)
//	if err != nil {
//		return err
//	}
//
//	_, err = http.Post("http://"+apinametg+"/telegram/del", "application/json", bytes.NewBuffer(data))
//	if err != nil {
//		return err
//
//	}
//	return nil
//}
//
//func EditMessageOLD(chatId, mid string, text, ParseMode string) error {
//	m := models.EditText{
//		Text:      text,
//		Channel:   chatId,
//		MessageId: mid,
//	}
//
//	url := "http://" + apinametg + "/telegram/edit"
//	if ParseMode != "" {
//		url = url + "?parse=" + ParseMode
//	}
//	data, err := json.Marshal(m)
//	if err != nil {
//		return err
//	}
//
//	resp, errr := http.Post(url, "application/json", bytes.NewBuffer(data))
//	defer resp.Body.Close()
//	if errr != nil {
//		return errr
//	}
//	if resp.StatusCode != http.StatusOK {
//		body, _ := ioutil.ReadAll(resp.Body)
//		return errors.New(fmt.Sprintf("Code %d body:%s", resp.StatusCode, string(body)))
//	}
//	return nil
//}
