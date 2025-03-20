package DiscordClient

import (
	"discord/models"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (d *Discord) logicScoreboardSetting(m *discordgo.MessageCreate) bool {
	if m.Message.WebhookID != "" {
		d.logicMixWebhook(m)
		return true
	}

	afterScoreboard, found := strings.CutPrefix(m.Content, ".scoreboard")
	if found {
		text := ""
		afterWebhook, foundWebhook := strings.CutPrefix(afterScoreboard, " webhook ")
		if afterWebhook == "name" {
			text = "you can't use the 'name', it's not unique"
			foundWebhook = false
		}
		if foundWebhook {
			scoreboardReadName := d.storage.Scoreboard.ScoreboardReadName(afterWebhook)
			if scoreboardReadName == nil {
				d.storage.Scoreboard.ScoreboardInsertParam(models.ScoreboardParams{
					Name:              afterWebhook,
					ChannelWebhook:    m.ChannelID,
					ChannelScoreboard: "",
				})
				text = "now the bot will wait here for webhooks from the game, connect another channel to display the leaderboard"
			} else {
				text = "this channel is already listened to by a bot to receive webhooks from the game.Name " + scoreboardReadName.Name
				d.log.Info("found " + afterWebhook + " in scoreboard")
			}
		}
		afterHere, fountHere := strings.CutPrefix(afterScoreboard, " here ")
		if afterHere == "name" {
			text = "you can't use the 'name', it's not unique"
			fountHere = false
		}
		if fountHere {
			scoreboard := d.storage.Scoreboard.ScoreboardReadName(afterHere)
			if scoreboard != nil {
				scoreboard.ChannelScoreboard = m.ChannelID
				d.storage.Scoreboard.ScoreboardUpdateParam(*scoreboard)
				text = "now the leaderboard will be displayed here"
			} else {
				d.log.Info("not found " + afterHere + " in scoreboard")
				text = "it is impossible to connect the leaderboard without having a channel of incoming data from the game via webhook"
			}
		}
		if text == "" && !fountHere && !foundWebhook {
			text = "To set up automatic display of red star event leaders, you need to do several things:\n" +
				"1) set up sending webhook to you in the game in a hidden channel in discord\n" +
				"2) execute the command to connect the bot to listening to this channel, come up with a unique name or use the corporation name in the command '.scoreboard webhook name' where name is a unique name that will link the data in the bot.\n" +
				"3) execute the command in the channel open for viewing your corporation where the leaderboard will be displayed '.scoreboard here name'"
		}
		d.DeleteMesageSecond(m.ChannelID, m.ID, twoDay)
		d.SendChannelDelSecond(m.ChannelID, text, twoDay)
		return true
	}
	return false
}

func (d *Discord) logicMixWebhook(m *discordgo.MessageCreate) {
	if m.Attachments != nil && len(m.Attachments) > 0 && m.Attachments[0].Filename == "data.json" {
		scoreboard := d.storage.Scoreboard.ScoreboardReadWebhookChannel(m.ChannelID)
		if scoreboard != nil {
			err := d.FetchJSON(m.Attachments[0].URL, scoreboard)
			if err != nil {
				d.log.ErrorErr(err)
				return
			}
		} else {
			fmt.Printf("not found setting for scoreboard channel %s\n", m.ChannelID)
		}
	}

	//нужно получить текущий ивент

}

func (d *Discord) FetchJSON(url string, params *models.ScoreboardParams) error {
	// Отправка GET-запроса
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("ошибка при отправке запроса: %v", err)
	}
	defer resp.Body.Close()

	// Проверка статуса ответа
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неожиданный статус ответа: %s", resp.Status)
	}

	// Чтение тела ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка при чтении ответа: %v", err)
	}

	// Декодирование JSON в map для удобного отображения
	var jsonData map[string]interface{}
	if err = json.Unmarshal(body, &jsonData); err != nil {
		return fmt.Errorf("ошибка при парсинге JSON: %v", err)
	}

	// Форматированный вывод JSON
	formattedJSON, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка при форматировании JSON: %v", err)
	}
	switch jsonData["EventType"] {
	case "RedStarStarted":
		var rss models.RedStarEvent
		err = rss.UnmarshalJSON(formattedJSON)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Sprintf("RedStarStarted %+v\n", rss) //not usage
	case "RedStarEnded":
		var rse models.RedStarEvent
		if err = rse.UnmarshalJSON(formattedJSON); err != nil {
			fmt.Println(err)
		}
		err = d.RedStarEnded(rse, params)
		if err != nil {
			fmt.Println(err)
			return err
		}
	default:
		fmt.Println(string(formattedJSON))
	}

	//fmt.Println(len(string(formattedJSON)))
	return nil
}
func getSeasonNumber(text string) int {
	re := regexp.MustCompile(`Season (\d+)`)
	matches := re.FindStringSubmatch(text)

	if len(matches) < 2 {
		return 0
	}
	seasonNumber, _ := strconv.Atoi(matches[1])
	return seasonNumber
}

func (d *Discord) RedStarEnded(rse models.RedStarEvent, params *models.ScoreboardParams) error {
	if rse.RSEventPoints != 0 {
		nextDateStart, nextDateStop, message := d.storage.Scoreboard.ReadEventScheduleAndMessage()
		date1 := time.Now().UTC().Format("02-01-2006")
		date2 := time.Now().UTC().Add(24 * time.Hour).Format("02-01-2006")
		eventId := time.Now().UTC().Day() - 10 //test id
		if date1 == nextDateStart || date2 == nextDateStop {
			eventId = getSeasonNumber(message)
		}

		points := rse.RSEventPoints / len(rse.Players)
		for _, player := range rse.Players {
			err := d.storage.Battles.BattlesInsert(models.Battles{
				EventId:  eventId,
				CorpName: params.Name,
				Name:     player.PlayerName,
				Level:    rse.StarLevel,
				Points:   points,
			})
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
	} else {
		for _, player := range rse.Players {
			err := d.storage.Battles.BattlesTopInsert(models.BattlesTop{
				CorpName: params.Name,
				Name:     player.PlayerName,
				Level:    rse.StarLevel,
				Count:    1, //+1
			})
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
	}
	return nil
}
