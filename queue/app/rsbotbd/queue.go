package rsbotbd

import (
	"queue/models"
	"strconv"
	"strings"
)

func (q *Queue) GetQueueLevel(level string) (m map[string][]models.Tumcha) {

	if level == "" {
		return
	}

	tt := q.GetDBQueue()
	var t []models.Tumcha
	m = make(map[string][]models.Tumcha)

	after, found := strings.CutPrefix(level, "drs")

	if found {
		lvl, err := strconv.Atoi(after)
		if err != nil {
			q.log.ErrorErr(err)
		}
		for _, tumcha := range tt {
			if tumcha.Vid == "black" && lvl == tumcha.Level {
				t = append(t, tumcha)
			}
		}
	}
	after, found = strings.CutPrefix(level, "rs")
	if found {
		lvl, err := strconv.Atoi(after)
		if err != nil {
			q.log.ErrorErr(err)
		}
		for _, tumcha := range tt {
			if tumcha.Vid != "black" && lvl == tumcha.Level {
				t = append(t, tumcha)
			}
		}
	}

	for _, tumcha := range t {
		chat := q.getname(tumcha.Chatid)
		m[chat] = append(m[chat], tumcha)
	}

	return m
}

func (q *Queue) GetQueueAll() (m map[string][]models.Tumcha) {
	tt := q.GetDBQueue()
	m = make(map[string][]models.Tumcha)

	for _, tumcha := range tt {
		chat := q.getname(tumcha.Chatid)
		m[chat] = append(m[chat], tumcha)
	}

	return m
}
