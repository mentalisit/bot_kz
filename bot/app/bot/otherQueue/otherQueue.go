package otherQueue

import (
	"fmt"
	"github.com/mentalisit/logger"
	"kz_bot/models"
	"time"
)

type OtherQ struct {
	log *logger.Logger
}

func NewOtherQ(log *logger.Logger) *OtherQ {
	return &OtherQ{log: log}
}

func (o *OtherQ) MyQueue() string {
	time.Sleep(5 * time.Second)
	queueAll, err := GetQueueAll()
	if err != nil {
		o.log.ErrorErr(err)
		return ""
	}
	var text string
	for s, structs := range queueAll {
		text += fmt.Sprintf("⚠️ %s\n", s)
		for _, queueStruct := range structs {
			text += fmt.Sprintf("🔥 %s кз %s %d\n", queueStruct.CorpName, queueStruct.Level, queueStruct.Count)
		}
	}
	return text
}

func (o *OtherQ) ReadingQueueByLevel(level string) (text string, err error) {
	queueLevel, err := GetQueueLevel(level)
	if err != nil {
		return "", err
	}
	if len(queueLevel) > 0 {
		var q []models.QueueStruct
		//nameBot
		for _, queues := range queueLevel {
			for _, i := range queues {
				q = append(q, i)
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
