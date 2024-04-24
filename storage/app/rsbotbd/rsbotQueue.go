package rsbotbd

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"sort"
	"storage/config"
	"strconv"
)

type tumcha struct {
	name     string
	level    int
	vid      string
	chatid   int
	timedown int
}

func GetQueue() string {
	corpName = make(map[int]string)
	corpName[-1001242024247] = "Союз"
	corpName[-1001582192116] = "Союз Академия"
	corpName[-1001265143636] = "Конгломерат"
	corpName[-1001386882184] = "Неизбежный рок"
	corpName[-1001295995727] = "RUS"
	corpName[-1001685747025] = "Best"

	log.Println()
	db, err := sql.Open("mysql", config.Instance.MySQl)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	rows, err := db.Query("select name,lvlkz,vid,chatid,timedown from sborkz WHERE active = 0")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var tt []tumcha

	for rows.Next() {
		var t tumcha
		err := rows.Scan(&t.name, &t.level, &t.vid, &t.chatid, &t.timedown)
		if err != nil {
			fmt.Println(err)
			continue
		}
		tt = append(tt, t)
	}
	sort.Slice(tt, func(i, j int) bool {
		return tt[i].chatid < tt[j].chatid
	})
	//go preparingToSendChat(tt)
	Chat := make(map[int][]tumcha)
	for _, t := range tt {
		Chat[t.chatid] = append(Chat[t.chatid], t)
	}
	var finalText string
	for _, tumchas := range Chat {
		//preparingToSendLevel(tumchas)

		level := make(map[int][]tumcha)
		for _, t := range tumchas {
			level[t.level] = append(level[t.level], t)
		}

		text := "⚠️ " + getname(tumchas[0].chatid) + "\n"
		for i, tumchasl := range level {
			text += fmt.Sprintf("очередь на %d\n", i)
			for id, t := range tumchasl {
				text += fmt.Sprintf("%d. %s  🕒%d\n", id+1, t.name, t.timedown)
			}
		}
		finalText += text
		fmt.Println(text)
	}
	return finalText

}

var corpName map[int]string

func getname(chatid int) string {
	text := corpName[chatid]
	if text == "" {
		return strconv.Itoa(chatid)
	}
	return text
}
