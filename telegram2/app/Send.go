package main

//
//import (
//	"github.com/celestix/gotgproto/ext"
//	"github.com/gotd/td/tg"
//	"strconv"
//)
//
//func SendText(ctx *ext.Context, ChatId, Text string) (mId string, err error) {
//	chatId, threadID := chat(ChatId)
//	mes := &tg.MessagesSendMessageRequest{
//		Message: Text,
//	}
//
//	if threadID != 0 {
//		r := &tg.InputReplyToMessage{
//			ReplyToMsgID: threadID,
//		}
//		mes.SetReplyTo(r)
//	}
//
//	message, err := ctx.SendMessage(chatId, mes)
//	if err != nil {
//		return "", err
//	}
//	mId = strconv.Itoa(message.ID)
//	return mId, nil
//}
