package logic

import (
	"bridge/models"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func (b *Bridge) ifPoll() {
	after, found := strings.CutPrefix(b.in.Text, ".poll")

	if found {
		arg := strings.Split(after, `"`)

		// Проверяем, был ли создан новый опрос
		if len(arg) > 3 {
			p := models.PollStruct{
				Author:      b.in.Sender,
				Question:    arg[1],
				CreateTime:  time.Now().Unix(),
				Config:      *b.in.Config,
				PollMessage: make(map[string]string),
			}

			for _, s := range arg[2:] {
				if len(s) > 2 {
					p.Options = append(p.Options, s)
				}
			}
			fmt.Printf("poll %+v\n %+v\n", p, p.Config)

			// Генерация ссылки для результатов
			p.UrlPoll = fmt.Sprintf("https://ws.mentalisit.myds.me/poll/%d", p.CreateTime)

			m := make(map[string]string)
			m["author"] = p.Author
			m["question"] = p.Question
			m["url"] = p.UrlPoll
			m["createTime"] = strconv.FormatInt(p.CreateTime, 10)

			for t, configs := range p.Config.Channel {
				for _, config := range configs {
					m["chatid"] = config.ChannelId
					if t == "ds" {
						p.PollMessage[config.ChannelId] = b.discord.SendPollChannel(m, p.Options)
					} else if t == "tg" {
						p.PollMessage[config.ChannelId] = b.telegram.SendPollChannel(m, p.Options)
					}
				}
			}
			// Отправка опроса в Discord
			//if len(p.Config.ChannelDs) > 0 {
			//	for _, ds := range p.Config.ChannelDs {
			//		m["chatid"] = ds.ChannelId
			//		p.PollMessage[ds.ChannelId] = b.discord.SendPollChannel(m, p.Options)
			//	}
			//}
			//
			//// Отправка опроса в Telegram
			//if len(p.Config.ChannelTg) > 0 {
			//	for _, tg := range p.Config.ChannelTg {
			//		m["chatid"] = tg.ChannelId
			//		p.PollMessage[tg.ChannelId] = b.telegram.SendPollChannel(m, p.Options)
			//	}
			//}

			bytes, err := json.Marshal(p)
			if err != nil {
				b.log.ErrorErr(err)
				return
			}
			pathFile := fmt.Sprintf("docker/poll/%d", p.CreateTime)
			err = os.WriteFile(pathFile, bytes, 0666)
			if err != nil {
				b.log.ErrorErr(err)
				return
			}
			return
		}

		// Обработка выбора ответа
		split := strings.Split(after, ".")
		pathFile := strings.Replace(split[0], " ", "docker/poll/", 1)
		choice := split[1]
		choiceint, _ := strconv.Atoi(choice)

		if !strings.HasPrefix(split[0], " 17") {
			return
		}

		oFile, err := os.ReadFile(pathFile)
		if err != nil {
			b.log.ErrorErr(err)
			return
		}
		var r models.PollStruct
		err = json.Unmarshal(oFile, &r)
		if err != nil {
			b.log.ErrorErr(err)
			return
		}

		// Проверяем, голосовал ли пользователь ранее
		userVoted := false
		for i, vote := range r.Votes {
			if vote.UserName == b.in.Sender {
				// Если пользователь проголосовал, обновляем его ответ
				r.Votes[i].Answer = choice
				userVoted = true
				b.ifTipSendDel(b.in.Sender + " выбор обновлен " + r.Options[choiceint-1])
				break
			}
		}

		// Если пользователь не голосовал, добавляем новый голос
		if !userVoted {
			r.Votes = append(r.Votes, models.Votes{
				Type:     b.in.Tip,
				Channel:  b.in.ChatId,
				UserName: b.in.Sender,
				Answer:   choice,
			})
			b.ifTipSendDel(b.in.Sender + " внесено " + r.Options[choiceint-1])
		}

		fmt.Printf("PollStruct after vote: %+v\n", r)

		bytes, err := json.Marshal(r)
		if err != nil {
			b.log.ErrorErr(err)
			return
		}

		err = os.WriteFile(pathFile, bytes, 0666)
		if err != nil {
			b.log.ErrorErr(err)
			return
		}
	}
}
