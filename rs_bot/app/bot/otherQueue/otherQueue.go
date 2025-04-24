package otherQueue

import (
	"fmt"
	"github.com/mentalisit/logger"
	"rs/models"
	"strconv"
	"strings"
)

type OtherQ struct {
	log *logger.Logger
}

func NewOtherQ(log *logger.Logger) *OtherQ {
	return &OtherQ{log: log}
}

func (o *OtherQ) MyQueue() string {
	queueAll, err := GetQueueAll()
	if err != nil {
		fmt.Println(err)
		return "Сервис временно недоступен "
	}
	var text string
	for s, structs := range queueAll {
		text += fmt.Sprintf("⚠️ %s\n", s)
		for _, queueStruct := range structs {
			text += fmt.Sprintf("	🔥 %s кз %s %d\n", queueStruct.CorpName, queueStruct.Level, queueStruct.Count)
		}
	}
	if text == "" {
		text = "нет активных очередей"
	}
	return text
}

func (o *OtherQ) ReadingQueueByLevel(level, corp string) (text string, err error) {
	queueLevel, err := GetQueueLevel(level)
	if err != nil {
		return "", err
	}
	if len(queueLevel) > 0 {
		var q []models.QueueStruct
		//nameBot
		for _, queues := range queueLevel {
			for _, i := range queues {
				//fmt.Printf("corp '%s' deleteChannel '%s'", deleteChannelName(corp), i.CorpName)
				if i.CorpName != deleteChannelName(corp) {
					q = append(q, i)
				}
			}
		}
		if len(q) > 0 {
			text = "Другие активные очереди"
		}
		for _, queueStruct := range q {
			text += fmt.Sprintf("\n%s в очереди %d", queueStruct.CorpName, queueStruct.Count)
		}
	}

	return text, nil
}
func deleteChannelName(Corpname string) string {
	split := strings.Split(Corpname, ".")
	if len(split) == 2 {
		return split[0]
	} else {
		i := strings.Split(Corpname, "/")
		if len(i) == 2 {
			return i[0]
		}
	}

	return Corpname
}
func (o *OtherQ) NeedRemoveOtherQueue(ss []string) {
	uId, err := GetUseridTumcha()
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(ss) > 0 {
		for _, s := range ss {
			for _, i := range uId {
				if strconv.FormatInt(i, 10) == s {
					o.log.InfoStruct("NeedRemoveOtherQueue", s)
				}
			}
		}
	}
}
