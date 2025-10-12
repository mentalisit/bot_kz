package tg

import (
	"bridge/models"
	"context"
	"fmt"
	"sync"

	"github.com/mentalisit/logger"
	"google.golang.org/grpc"
)

type Client struct {
	conn   *grpc.ClientConn
	client TelegramServiceClient
	log    *logger.Logger
}

func NewClient(log *logger.Logger) *Client {
	conn, err := grpc.Dial("telegram:50051", grpc.WithInsecure())
	if err != nil {
		log.ErrorErr(err)
		return nil
	}
	fmt.Println("connect to grpc discord ok")
	return &Client{
		conn:   conn,
		client: NewTelegramServiceClient(conn),
		log:    log,
	}
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) DeleteMessage(ChatId string, messageID string) {
	er, err := c.client.DeleteMessage(context.Background(), &DeleteMessageRequest{
		Chatid: ChatId,
		Mesid:  messageID,
	})
	if err != nil {
		c.log.ErrorErr(err)
	} else if er.GetErrorMessage() != "" {
		c.log.Error(er.GetErrorMessage())
	}
}

func (c *Client) SendChannelDelSecond(chatId, text string, second int) {
	_, err := c.client.SendChannelDelSecond(context.Background(), &SendMessageRequest{
		Text:   text,
		ChatID: chatId,
		Second: int32(second),
	})
	if err != nil {
		c.log.ErrorErr(err)
		return
	}
}
func (c *Client) SendPollChannel(m map[string]string, options []string) string {
	poll, err := c.client.SendPoll(context.Background(), &SendPollRequest{
		Data:    m,
		Options: options,
	})
	if err != nil {
		c.log.ErrorErr(err)
		return ""
	}
	return poll.GetText()
}
func (c *Client) sendBridgeArrayMessages(inMessenger models.BridgeSendToMessenger) (MessageIds []models.MessageIds) {
	req := &SendBridgeArrayMessagesRequest{
		Text:      inMessenger.Text,
		ChannelID: inMessenger.ChannelId,
		ReplyMap:  inMessenger.ReplyMap,
	}
	if len(inMessenger.Extra) > 0 {
		for _, i := range inMessenger.Extra {
			req.Extra = append(req.Extra, &FileInfo{
				Name:   i.Name,
				Data:   i.Data,
				Url:    i.URL,
				Size:   i.Size,
				FileID: i.FileID,
			})
		}

	}

	arrayMessages, err := c.client.SendBridgeArrayMessages(context.Background(), req)
	if err != nil {
		c.log.ErrorErr(err)
		return
	}
	var mids []models.MessageIds
	for _, id := range arrayMessages.GetMessageIds() {
		mids = append(mids, models.MessageIds{
			MessageId: id.MessageId,
			ChatId:    id.ChatId,
		})
	}
	return mids
}
func (c *Client) SendBridgeArrayMessage(resultChannel chan<- models.MessageIds, wg *sync.WaitGroup, inMessenger models.BridgeSendToMessenger) {
	defer wg.Done()

	ids := c.sendBridgeArrayMessages(inMessenger)
	for _, id := range ids {
		resultChannel <- id
	}
}

//
//type Telegram struct {
//	log *logger.Logger
//}
//
//func NewTelegram(log *logger.Logger) *Telegram {
//	return &Telegram{log: log}
//}
//
//const (
//	deleteMessage = "delete_message"
//	sendDel       = "send_del"
//)
//
//func (t *Telegram) DeleteMessage(ChatId string, messageID string) error {
//	m := apiRs{
//		FuncApi:   deleteMessage,
//		Channel:   ChatId,
//		MessageId: messageID,
//	}
//	a, err := convertAndSend(m)
//	if err != nil {
//		t.log.ErrorErr(err)
//		return err
//	}
//	fmt.Printf("DelMessageSecond %+v\n", a)
//	return nil
//}
//func (t *Telegram) SendChannelDelSecond(chatId, text string, second int) {
//	m := apiRs{
//		FuncApi: sendDel,
//		Text:    text,
//		Channel: chatId,
//		Second:  second,
//	}
//
//	a, err := convertAndSend(m)
//	if err != nil {
//		t.log.ErrorErr(err)
//		return
//	}
//	fmt.Printf("SendChannelDelSecond %+v\n", a)
//}
//func (t *Telegram) SendPollChannel(m map[string]string, options []string) string {
//	s := models.Request{
//		Data:    m,
//		Options: options,
//	}
//
//	data, err := json.Marshal(s)
//	if err != nil {
//		t.log.ErrorErr(err)
//		return ""
//	}
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	req, err := http.NewRequestWithContext(ctx, "POST", "http://telegram/send_poll", bytes.NewBuffer(data))
//	if err != nil {
//		t.log.ErrorErr(err)
//		return ""
//	}
//	req.Header.Set("Content-Type", "application/json")
//
//	client := &http.Client{}
//
//	resp, err := client.Do(req)
//	if err != nil {
//		t.log.ErrorErr(err)
//		return ""
//	}
//	defer resp.Body.Close()
//
//	var a string
//	err = json.NewDecoder(resp.Body).Decode(&a)
//	if err != nil {
//		body, _ := ioutil.ReadAll(resp.Body)
//		t.log.ErrorErr(err)
//		t.log.Info(string(body))
//		return ""
//	}
//
//	fmt.Printf("SendPoll %+v\n", a)
//	return a
//}
//func (t *Telegram) SendBridgeArrayMessage(chatid []string, text string, extra []models.FileInfo, resultChannel chan<- models.MessageIds, wg *sync.WaitGroup) {
//
//	defer wg.Done()
//
//	m := models.BridgeSendToMessenger{
//		Text:      text,
//		ChannelId: chatid,
//		Extra:     extra,
//	}
//
//	data, err := json.Marshal(m)
//	if err != nil {
//		t.log.ErrorErr(err)
//		return
//	}
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	url := "http://telegram/bridge"
//
//	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
//	if err != nil {
//		t.log.ErrorErr(err)
//		return
//	}
//	req.Header.Set("Content-Type", "application/json")
//	client := &http.Client{}
//
//	resp, err := client.Do(req)
//	if err != nil {
//		t.log.ErrorErr(err)
//		return
//	}
//	defer resp.Body.Close()
//
//	var MessageIds []models.MessageIds
//	err = json.NewDecoder(resp.Body).Decode(&MessageIds)
//	if err != nil {
//		t.log.Info(fmt.Sprintf("err resp.Body %+v\n", resp.Body))
//		t.log.ErrorErr(err)
//		return
//	}
//
//	for _, id := range MessageIds {
//		resultChannel <- id
//	}
//}

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
//type apiRs struct {
//	FuncApi   string `json:"funcApi"`
//	Text      string `json:"text"`
//	Channel   string `json:"channel"`
//	MessageId string `json:"messageId"`
//	//ParseMode string `json:"parseMode"`
//	Second int `json:"second"`
//	//LevelRs   string `json:"levelRs"`
//	//Levels    []string `json:"levels"`
//	//UserName  string `json:"userName"`
//}
//type answer struct {
//	ArrString string    `json:"arrString"`
//	ArrInt    int       `json:"arrInt"`
//	ArrBool   bool      `json:"arrBool"`
//	ArrError  error     `json:"arrError"`
//	Time      time.Time `json:"time"`
//}
