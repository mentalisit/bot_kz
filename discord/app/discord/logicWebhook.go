package DiscordClient

import (
	"bytes"
	"context"
	"discord/config"
	"discord/models"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

func (d *Discord) logicMixWebhook(m *discordgo.Message) {
	if m.Attachments != nil && len(m.Attachments) > 0 && m.Attachments[0].Filename == "data.json" {
		scoreboard := d.storage.Scoreboard.ScoreboardReadWebhookChannel(m.ChannelID)
		if scoreboard == nil {
			CorpName, _ := d.FetchJSON(m.Attachments[0].URL, nil, m.Timestamp.Unix())
			if CorpName != "" {
				scoreboard = &models.ScoreboardParams{
					Name:              CorpName,
					ChannelWebhook:    m.ChannelID,
					ChannelScoreboard: "",
					LastMessageID:     "",
				}
				d.storage.Scoreboard.ScoreboardInsertParam(*scoreboard)
				d.Send(m.ChannelID, fmt.Sprintf("Конфиг создан автоматически, название '%s'", scoreboard.Name))

				g, err := d.S.Guild(m.GuildID)
				if err != nil {
					d.log.ErrorErr(err)
				}
				channel, _ := d.S.Channel(m.ChannelID)
				d.log.Info(fmt.Sprintf("not found setting for scoreboard channel %s Guild %s\n", channel.Name, g.Name))
			}
		}
		if scoreboard == nil {
			d.log.Error("WTF ??")
			return
		}
		_, err := d.FetchJSON(m.Attachments[0].URL, scoreboard, m.Timestamp.Unix())
		if err != nil {
			d.log.ErrorErr(err)
			return
		}
		scoreboard.LastMessageID = m.ID
		d.storage.Scoreboard.ScoreboardUpdateParamLastMessageId(*scoreboard)
	}
}

func (d *Discord) FetchJSON(url string, params *models.ScoreboardParams, tsUnix int64) (string, error) {
	// Отправка GET-запроса
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("ошибка при отправке запроса: %v", err)
	}
	defer resp.Body.Close()

	// Проверка статуса ответа
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("неожиданный статус ответа: %s", resp.Status)
	}

	// Чтение тела ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка при чтении ответа: %v", err)
	}

	// Декодирование JSON в map для удобного отображения
	var jsonData map[string]interface{}
	if err = json.Unmarshal(body, &jsonData); err != nil {
		return "", fmt.Errorf("ошибка при парсинге JSON: %v", err)
	}
	if params == nil {
		corporation := jsonData["Corporation"].(map[string]interface{})
		return corporation["CorporationName"].(string), nil
	}

	// Форматированный вывод JSON
	formattedJSON, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return "", fmt.Errorf("ошибка при форматировании JSON: %v", err)
	}

	switch jsonData["EventType"] {
	case "RedStarStarted":
		var rss models.RedStarEvent
		err = rss.UnmarshalJSON(formattedJSON)
		if err != nil {
			return "", err
		}
		d.RedStarStart(rss, params)

		d.storage.Db.InsertWebhookType(tsUnix, rss.Corporation.CorporationName, rss.EventType, string(body))

	case "RedStarEnded":
		var rse models.RedStarEvent
		if err = rse.UnmarshalJSON(formattedJSON); err != nil {
			fmt.Println(err)
		}
		err = d.RedStarEnded(rse, params)
		if err != nil {
			return "", err
		}
		d.storage.Db.InsertWebhookType(tsUnix, rse.Corporation.CorporationName, rse.EventType, string(body))

		if params.Name == "rus" || params.Name == "soyuz" {
			d.storage.Db.InsertWebhook(tsUnix, params.Name, string(body))
		}

	case "WhiteStarStarted":
		var wss models.WhiteStarStarted
		err = json.Unmarshal(formattedJSON, &wss)
		if err != nil {
			return "", err
		}
		d.WhiteStarStarted(wss, params)
		d.storage.Db.InsertWebhookType(tsUnix, wss.Corporation.CorporationName, wss.EventType, string(body))

		d.sendWhiteStarData(body)

		if params.Name == "rus" || params.Name == "soyuz" {
			d.storage.Db.InsertWebhook(tsUnix, params.Name, string(body))
		}

	case "WhiteStarEnded":
		var wse models.WhiteStarEnded
		err = json.Unmarshal(formattedJSON, &wse)
		if err != nil {
			return "", err
		}
		d.WhiteStarEnded(wse, params)
		d.storage.Db.InsertWebhookType(tsUnix, wse.Corporation.CorporationName, wse.EventType, string(body))

		d.sendWhiteStarData(body)

		if params.Name == "rus" || params.Name == "soyuz" {
			d.storage.Db.InsertWebhook(tsUnix, params.Name, string(body))
		}

	default:
		d.log.InfoStruct("default "+params.Name, jsonData)
		fmt.Println(string(body))
	}
	return "", nil
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

		text = text + d.CombineNamesWithNameAliases(rse.Players)

		if rse.RSEventPoints != 0 {
			text = text + fmt.Sprintf("Получено %d", rse.RSEventPoints)
		}

		d.SendWebhook(text, rse.Corporation.CorporationName, params.ChannelWebhook, rse.Corporation.GetAvatar())
	} else {
		d.SendWebhook(rse.EventType+" \nучастники не обнаружены", rse.Corporation.CorporationName, params.ChannelWebhook, rse.Corporation.GetAvatar())
	}
	if params.Name == "best" || params.ChannelScoreboard != "" {
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
					Name:     d.CombineNames(player.PlayerName),
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
					Name:     d.CombineNames(player.PlayerName),
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

		text = text + d.CombineNamesWithNameAliases(rse.Players)

		if rse.RSEventPoints != 0 {
			text = text + fmt.Sprintf("Получено %d", rse.RSEventPoints)
		}

		d.SendWebhook(text, rse.Corporation.CorporationName, params.ChannelWebhook, rse.Corporation.GetAvatar())
	} else {
		d.SendWebhook(rse.EventType+" \nучастники не обнаружены", rse.Corporation.CorporationName, params.ChannelWebhook, rse.Corporation.GetAvatar())
	}

}
func (d *Discord) WhiteStarStarted(wss models.WhiteStarStarted, params *models.ScoreboardParams) {
	if wss.Opponent.CorporationName != "" {
		text := wss.EventType
		if len(wss.OurParticipants) > 0 {
			text = text + fmt.Sprintf("\n  Участники: \n")
			text = text + d.CombineNamesWithNameAliases(wss.OurParticipants)
		}
		if len(wss.OpponentParticipants) > 0 {
			text = text + fmt.Sprintf("\n  Противник %s: \n", wss.Opponent.CorporationName)
			text = text + d.CombineNamesWithNameAliases(wss.OpponentParticipants)
		}
		d.SendWebhook(text, params.Name, params.ChannelWebhook, wss.Corporation.GetAvatar())
	}
}
func (d *Discord) WhiteStarEnded(wse models.WhiteStarEnded, params *models.ScoreboardParams) {
	text := fmt.Sprintf("%s\nПротивник: %s \n Opponent %d - Our %d \nXPGained %d\n",
		wse.EventType, wse.Opponent.CorporationName, wse.OpponentScore, wse.OurScore, wse.XPGained)
	d.SendWebhook(text, params.Name, params.ChannelWebhook, wse.Corporation.GetAvatar())
}
func oldCombineNames(r string) string {
	switch r {
	case "Mchuleft", "Valenvaryon":
		return "Mchuleft"
	case "Overturned", "Overturned-1.1":
		return "Overturned"
	case "RedArrow", "Light Matter", "Dark Matter", "Drake", "Leviathan":
		return "RedArrow"
	case "falcon_2":
		return "falcon_2(Миша)"
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
	case "Джонни_De", "JonnyDe", "Red-is":
		return "JonnyDe"
	case "N@N", "ChubbChubbs":
		return "ChubbChubbs"
	case "Nixonblade", "TimA", "Ted", "Коньячный ЗАВОД":
		return "Nixonblade"
	case "Widowmaker":
		return "retresh90"
	case "Persil", "Pasis", "ILTS":
		return "Shishu"

	case "Iterius", "Furia":
		return "Iterius"
	case "Gennadiy Reng", "Mr.Reng":
		return "Mr.Reng"
	case "ololoki":
		return "trololo"
	case "Vera", "Prizrak1astu":
		return "Mad Max"
	case "КОРСАРИК", "Танкист39":
		return "KOPCAP"
	case "PovAndy":
		return "vdruzh"
	case "SpecBalrog":
		return "Tauren"
	case "WING☯CHUN":
		return "Jeff1143"

	default:
		return r
	}
}

// send WhiteStarData to statistic bot
func (d *Discord) sendWhiteStarData(data []byte) {
	url := "https://api.tsl.rocks/datajson?token=" + config.Instance.WsToken
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(data))
	if err != nil {
		slog.Error("Request creation failed: " + err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Request failed: " + err.Error())
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		if resp.StatusCode == 409 {
			fmt.Println("already recorded")
			return
		}
		slog.Error("Request failed Status: " + resp.Status)
		return
	}
}

func (d *Discord) CombineNames(input string) string {
	if d.NameAliases == nil {
		nameAliases, err := d.storage.Battles.LoadNameAliases()
		if err != nil {
			d.log.ErrorErr(err)
			return input
		}
		d.NameAliases = nameAliases
	}
	if val, ok := d.NameAliases[input]; ok {
		return val
	}
	return input
}
func (d *Discord) CombineNamesWithNameAliases(in []models.Participant) string {
	var text string
	for i, player := range in {
		combineName := d.CombineNames(player.PlayerName)
		if combineName == player.PlayerName {
			text = text + fmt.Sprintf("%d. %s\n", i+1, player.PlayerName)
		} else {
			text = text + fmt.Sprintf("%d. %s(%s)\n", i+1, combineName, player.PlayerName)
		}
	}
	return text
}
