package main

//
//import (
//	"fmt"
//	"github.com/celestix/gotgproto/types"
//	"github.com/gotd/td/tg"
//	"strconv"
//	"strings"
//)
//
//func chat(chatid string) (chatId int64, threadID int) {
//	a := strings.SplitN(chatid, "/", 2)
//	chatId, _ = strconv.ParseInt(a[0], 10, 64)
//
//	if len(a) > 1 {
//		threadID, _ = strconv.Atoi(a[1])
//	}
//	fmt.Printf("chatId %d threadID %d\n", chatId, threadID)
//	return chatId, threadID
//}
//
//func getChatId(m *types.Message) string {
//	ChatId := ""
//	if m != nil {
//		if m.PeerID != nil {
//			switch c := m.PeerID.(type) {
//			case *tg.PeerChannel:
//				ChatId = strconv.FormatInt(c.ChannelID, 10)
//			case *tg.PeerChat:
//				ChatId = strconv.FormatInt(c.ChatID, 10)
//			case *tg.PeerUser:
//				ChatId = strconv.FormatInt(c.UserID, 10)
//
//			default:
//				fmt.Println(m.PeerID)
//			}
//		}
//		if m.ReplyTo != nil {
//			value, ok := m.GetReplyTo()
//			if ok {
//				switch v := value.(type) {
//				case *tg.MessageReplyHeader:
//					if v.ForumTopic {
//						ChatId = fmt.Sprintf("%s/%d", ChatId, v.ReplyToMsgID)
//					}
//				default:
//					fmt.Println(value)
//				}
//			}
//		}
//	}
//	return ChatId
//}
