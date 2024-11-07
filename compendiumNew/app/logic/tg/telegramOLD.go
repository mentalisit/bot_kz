package tg

//
//type apiRs struct {
//	FuncApi   string `json:"funcApi"`
//	Text      string `json:"text"`
//	Channel   string `json:"channel"`
//	MessageId string `json:"messageId"`
//	ParseMode string `json:"parseMode"`
//	//Second    int      `json:"second"`
//	//LevelRs   string   `json:"levelRs"`
//	//Levels    []string `json:"levels"`
//	//UserName  string   `json:"userName"`
//}
//type answer struct {
//	ArrString string    `json:"arrString"`
//	ArrInt    int       `json:"arrInt"`
//	ArrBool   bool      `json:"arrBool"`
//	ArrError  error     `json:"arrError"`
//	Time      time.Time `json:"time"`
//}
//
//func convertAndSend(m any) (a answer, err error) {
//	fmt.Printf("convertAndSend %+v\n", m)
//	data, err := json.Marshal(m)
//	if err != nil {
//		return
//	}
//
//	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//	defer cancel()
//
//	req, err := http.NewRequestWithContext(ctx, "POST", "http://telegram/func", bytes.NewBuffer(data))
//	if err != nil {
//		return
//	}
//	req.Header.Set("Content-Type", "application/json")
//
//	client := &http.Client{}
//
//	resp, err := client.Do(req)
//	if err != nil {
//		return
//	}
//	defer resp.Body.Close()
//
//	err = json.NewDecoder(resp.Body).Decode(&a)
//	if err != nil {
//		printBody(resp.Body)
//	}
//	return a, err
//}
//func printBody(r io.Reader) {
//	body, _ := ioutil.ReadAll(r)
//	fmt.Println("printBody ", string(body))
//}
//
//const (
//	send          = "send"
//	deleteMessage = "delete_message"
//	editMessage   = "edit_message"
//	urlTg         = "http://telegram" //192.168.100.14  telegram
//)
//
//// отправка сообщения в телегу
//func Send(chatId string, text string) (string, error) {
//	m := apiRs{
//		FuncApi: send,
//		Text:    text,
//		Channel: chatId,
//	}
//	fmt.Printf("convertAndSendSend %+v\n", m)
//	data, err := json.Marshal(m)
//	if err != nil {
//		return "", err
//	}
//
//	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//	defer cancel()
//
//	req, err := http.NewRequestWithContext(ctx, "POST", "http://telegram/func", bytes.NewBuffer(data))
//	if err != nil {
//		return "", err
//	}
//	req.Header.Set("Content-Type", "application/json")
//
//	client := &http.Client{}
//
//	resp, err := client.Do(req)
//	if err != nil {
//		return "", err
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode == http.StatusForbidden {
//		return "", errors.New("forbidden")
//	}
//	var a answer
//	err = json.NewDecoder(resp.Body).Decode(&a)
//	if err != nil {
//		printBody(resp.Body)
//	}
//
//	return a.ArrString, a.ArrError
//}
//func SendPic(channelId string, text string, pic []byte) error {
//	m := models.SendPic{
//		Text:    text,
//		Channel: channelId,
//		Pic:     pic,
//	}
//
//	data, err := json.Marshal(m)
//	if err != nil {
//		return err
//	}
//
//	resp, err := http.Post(urlTg+"/sendPic", "application/json", bytes.NewBuffer(data))
//	if err != nil {
//		resp, err = http.Post("http://192.168.100.131:802/data", "application/json", bytes.NewBuffer(data))
//		if err != nil {
//			return err
//		}
//	}
//	defer resp.Body.Close()
//	return nil
//}
//func DeleteMessage(ChatId string, messageID string) error {
//	m := apiRs{
//		FuncApi:   deleteMessage,
//		Channel:   ChatId,
//		MessageId: messageID,
//	}
//	a, err := convertAndSend(m)
//	if err != nil {
//		return err
//	}
//	fmt.Printf("DeleteMessage %+v\n", a)
//	return nil
//}
//func EditMessage(channel, messageID string, text, parse string) error {
//	m := apiRs{
//		FuncApi:   editMessage,
//		Text:      text,
//		Channel:   channel,
//		MessageId: messageID,
//		ParseMode: parse,
//	}
//	a, err := convertAndSend(m)
//	if err != nil {
//		return err
//	}
//	fmt.Printf("EditMessage %+v\n", a)
//	return nil
//}
