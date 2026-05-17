package webServer

//
//import (
//	"bridge/models"
//	"fmt"
//	"strconv"
//	"strings"
//	"time"
//
//	"github.com/google/uuid"
//)
//
//func (s *Server) ifPoll() {
//	after, found := strings.CutPrefix(b.in.Text, ".poll")
//
//	if found {
//		arg := strings.Split(after, `"`)
//
//		// Проверяем, был ли создан новый опрос
//		if len(arg) > 3 {
//			p := models.Poll2Struct{
//				Author:      b.in.Sender,
//				Question:    arg[1],
//				CreateTime:  time.Now().Unix(),
//				Config:      *b.in.Config,
//				PollMessage: make(map[string]string),
//			}
//			guildV2, err := b.db.GuildGetChannel(b.in.GuildId)
//			if err == nil && guildV2 != nil {
//				p.Gid = guildV2.GId.String()
//			}
//
//			for _, s := range arg[2:] {
//				if len(s) > 2 {
//					p.Options = append(p.Options, s)
//				}
//			}
//			fmt.Printf("poll2 %+v\n %+v\n", p, p.Config)
//
//			// Генерация ссылки для результатов
//			p.UrlPoll = fmt.Sprintf("https://mentalisit.myds.me/web/poll.html?id=%d", p.CreateTime)
//
//			m := make(map[string]string)
//			m["author"] = p.Author
//			m["question"] = p.Question
//			m["url"] = p.UrlPoll
//			m["createTime"] = strconv.FormatInt(p.CreateTime, 10)
//
//			for t, configs := range p.Config.Channel {
//				for _, config := range configs {
//					m["chatid"] = config.ChannelId
//					if t == "ds" {
//						p.PollMessage[config.ChannelId] = b.discord.SendPollChannel(m, p.Options)
//					} else if t == "tg" {
//						p.PollMessage[config.ChannelId] = b.telegram.SendPollChannel(m, p.Options)
//					}
//				}
//			}
//			err = b.db.CreatePoll(p)
//			if err != nil {
//				b.log.ErrorErr(err)
//			}
//			return
//		}
//
//		// Обработка выбора ответа
//		split := strings.Split(after, ".")
//		tsId := strings.TrimSpace(split[0]) // Убирает пробелы, переносы строк и табы
//		choice := split[1]
//
//		pollById, err := b.db.GetPollById(tsId)
//		fmt.Printf("poll %+v\n %+v\n", pollById, err)
//		if err != nil || pollById.Question == "" {
//			b.ifPollOld(after)
//			return
//		}
//
//		v := models.Votes2{
//			Type:     b.in.Tip,
//			Channel:  b.in.ChatId,
//			UserName: b.in.Sender,
//			Answer:   choice,
//		}
//
//		uid, err := b.db.FindMultiAccountUidByUserId(b.in.SenderId)
//		if err == nil && uid != uuid.Nil {
//			v.Uid = uid.String()
//		}
//
//		choiceInt, _ := strconv.Atoi(choice)
//		b.ifTipSendDel(b.in.Sender + " внесено " + pollById.Options[choiceInt-1])
//
//		fmt.Printf("Poll2Struct after vote: %+v\n", pollById)
//
//		err = b.db.UpsertVote(pollById.CreateTime, v)
//		if err != nil {
//			b.log.ErrorErr(err)
//			return
//		}
//	}
//}
