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
	"time"
)

func (d *Discord) logicMixWebhook(m *discordgo.Message) {
	if m.Attachments != nil && len(m.Attachments) > 0 && m.Attachments[0].Filename == "data.json" {
		scoreboard := d.storage.Scoreboard.ScoreboardReadWebhookChannel(m.ChannelID)
		if scoreboard != nil {
			err := d.FetchJSON(m.Attachments[0].URL, scoreboard, m.Timestamp.Unix())
			if err != nil {
				d.log.ErrorErr(err)
				return
			}
			scoreboard.LastMessageID = m.ID
			d.storage.Scoreboard.ScoreboardUpdateParamLastMessageId(*scoreboard)
		} else {
			d.log.Info(fmt.Sprintf("not found setting for scoreboard channel %s\n", m.ChannelID))
		}
	}
}

func (d *Discord) FetchJSON(url string, params *models.ScoreboardParams, tsUnix int64) error {
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

	d.storage.Db.InsertWebhook(tsUnix, params.Name, string(body))

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
			return err
		}
		d.RedStarStart(rss, params)

	case "RedStarEnded":
		var rse models.RedStarEvent
		if err = rse.UnmarshalJSON(formattedJSON); err != nil {
			fmt.Println(err)
		}
		err = d.RedStarEnded(rse, params)
		if err != nil {
			return err
		}

	case "WhiteStarStarted":
		var wss models.WhiteStarStarted
		err = json.Unmarshal(formattedJSON, &wss)
		if err != nil {
			return err
		}
		d.WhiteStarStarted(wss, params)

	case "WhiteStarEnded":
		var wse models.WhiteStarEnded
		err = json.Unmarshal(formattedJSON, &wse)
		if err != nil {
			return err
		}
		d.WhiteStarEnded(wse, params)

	default:
		d.log.InfoStruct("default "+params.Name, jsonData)
		fmt.Println(string(body))
	}
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
	if len(rse.Players) > 0 {
		text := "RedStarEnded "
		if rse.DarkRedStar {
			text = text + fmt.Sprintf("ТКЗ%d\n", rse.StarLevel)
		} else {
			text = text + fmt.Sprintf("КЗ%d\n", rse.StarLevel)
		}
		for i, player := range rse.Players {
			text = text + fmt.Sprintf("%d. %s\n", i+1, player.PlayerName)
		}
		if rse.RSEventPoints != 0 {
			text = text + fmt.Sprintf("Получено %d", rse.RSEventPoints)
		}

		d.SendWebhook(text, rse.Players[0].PlayerName, params.ChannelWebhook, "")
	} else {
		d.SendWebhook(rse.EventType+" \nучастники не обнаружены", params.Name, params.ChannelWebhook, "")
	}
	if params.ChannelScoreboard != "" {
		if rse.RSEventPoints != 0 {
			nextDateStart, nextDateStop, message := d.storage.Scoreboard.ReadEventScheduleAndMessage()
			date1 := time.Now().UTC().Format("02-01-2006")
			date2 := time.Now().UTC().Add(24 * time.Hour).Format("02-01-2006")
			eventId := 0
			if date1 == nextDateStart || date2 == nextDateStop {
				eventId = getSeasonNumber(message)
			}

			points := rse.RSEventPoints / len(rse.Players)
			for _, player := range rse.Players {
				err := d.storage.Battles.BattlesInsert(models.Battles{
					EventId:  eventId,
					CorpName: params.Name,
					Name:     combineNames(player.PlayerName),
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
					Name:     combineNames(player.PlayerName),
					Level:    rse.StarLevel,
					Count:    1, //+1
				})
				if err != nil {
					fmt.Println(err)
					return err
				}
			}
		}
	}

	return nil
}
func (d *Discord) RedStarStart(rse models.RedStarEvent, params *models.ScoreboardParams) {
	if len(rse.Players) > 0 {
		text := "RedStarStart "
		if rse.DarkRedStar {
			text = text + fmt.Sprintf("ТКЗ%d\n", rse.StarLevel)
		} else {
			text = text + fmt.Sprintf("КЗ%d\n", rse.StarLevel)
		}
		for i, player := range rse.Players {
			text = text + fmt.Sprintf("%d. %s\n", i+1, player.PlayerName)
		}
		if rse.RSEventPoints != 0 {
			text = text + fmt.Sprintf("Получено %d", rse.RSEventPoints)
		}

		d.SendWebhook(text, rse.Players[0].PlayerName, params.ChannelWebhook, "")
	} else {
		d.SendWebhook(rse.EventType+" \nучастники не обнаружены", params.Name, params.ChannelWebhook, "")
	}

}
func (d *Discord) WhiteStarStarted(wss models.WhiteStarStarted, params *models.ScoreboardParams) {
	if wss.Opponent.CorporationName != "" {
		text := wss.EventType
		if len(wss.OurParticipants) > 0 {
			text = text + fmt.Sprintf("\n  Участники: \n")
			for i, participant := range wss.OurParticipants {
				text = text + fmt.Sprintf("%d. %s\n", i+1, participant.PlayerName)
			}
		}
		if len(wss.OpponentParticipants) > 0 {
			text = text + fmt.Sprintf("\n  Противник %s: \n", wss.Opponent.CorporationName)
			for i, participant := range wss.OpponentParticipants {
				text = text + fmt.Sprintf("%d. %s\n", i+1, participant.PlayerName)
			}
		}
		d.SendWebhook(text, params.Name, params.ChannelWebhook, "")
	}
}
func (d *Discord) WhiteStarEnded(wse models.WhiteStarEnded, params *models.ScoreboardParams) {
	text := fmt.Sprintf("%s\nПротивник: %s \n Opponent %d - Our %d \nXPGained %d\n",
		wse.EventType, wse.Opponent.CorporationName, wse.OpponentScore, wse.OurScore, wse.XPGained)
	d.SendWebhook(text, params.Name, params.ChannelWebhook, "")

}
func combineNames(r string) string {
	switch r {
	case "Mchuleft", "Valenvaryon":
		return "Mchuleft"
	case "Overturned", "Overturned-1.1":
		return "Overturned"
	case "RedArrow", "Light Matter", "Dark Matter", "Drake":
		return "RedArrow"
	case "Коньячный ЗАВОД", "falcon_2":
		return "falcon_2"
	case "Silent_Noise", "WarySamurai1055":
		return "Silent_Noise"
	case "arsenium23", "Tabu 666", "Psyker", "kozlovskiu":
		return "Tabu"
	case "delov@r", "delovar", "Plague":
		return "delovar"
	case "iVanCoMik", "eVanCoMik", "VanCoMik":
		return "VanCoMik"
	case "Альтаир", "АЛЬТАИР", "Storm":
		return "Альтаир"
	case "Гэндальф серый", "Ёжик71":
		return "Ёжик71"
	case "Джонни_De", "JonnyDe":
		return "JonnyDe"
	case "N@N", "ChubbChubbs":
		return "ChubbChubbs"
	case "Nixonblade", "TimA":
		return "Nixonblade"
	case "Widowmaker":
		return "retresh90"
	case "Persil", "Pasis", "ILTS":
		return "Shishu"

	case "Iterius", "Furia":
		return "Iterius"
	case "Gennadiy Reng", "Mr.Reng":
		return "Mr.Reng"
	//case "DevilMyCry":
	//	return "Wait and Bleed"

	default:
		return r
	}
}
