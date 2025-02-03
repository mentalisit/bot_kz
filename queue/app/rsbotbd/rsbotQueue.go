package rsbotbd

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mentalisit/logger"
	"queue/config"
	"queue/models"
	"sort"
	"strconv"
)

type Queue struct {
	corpName map[int]string
	log      *logger.Logger
}

func NewQueue(log *logger.Logger) *Queue {
	q := &Queue{
		log:      log,
		corpName: make(map[int]string),
	}

	q.corpName[-1001242024247] = "Союз"
	q.corpName[-1001582192116] = "Союз Академия"
	q.corpName[-1001265143636] = "Конгломерат"
	q.corpName[-1001386882184] = "Неизбежный рок"
	q.corpName[-1001295995727] = "RUS"
	q.corpName[-1001685747025] = "Best"
	q.corpName[-1002098812155] = "Zvezdec"
	q.corpName[-1002075054059] = "Дом Датэ"
	q.corpName[-1002467616555] = "СССР"

	return q
}
func (q *Queue) GetDBQueue() (tt []models.Tumcha) {

	db, err := sql.Open("mysql", config.Instance.MySQl)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	rows, err := db.Query("select name,nameid,lvlkz,vid,chatid,timedown from sborkz WHERE active = 0")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Tumcha
		err = rows.Scan(&t.Name, &t.NameId, &t.Level, &t.Vid, &t.Chatid, &t.Timedown)
		if err != nil {
			fmt.Println(err)
			continue
		}
		tt = append(tt, t)
	}
	sort.Slice(tt, func(i, j int) bool {
		return tt[i].Chatid < tt[j].Chatid
	})

	return
}

func (q *Queue) getname(chatid int) string {
	text := q.corpName[chatid]
	if text == "" {
		return strconv.Itoa(chatid)
	}
	return text
}
