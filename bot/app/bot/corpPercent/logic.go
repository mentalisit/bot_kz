package corpPercent

import (
	"encoding/json"
	"fmt"
	"github.com/mentalisit/logger"
	"kz_bot/clients"
	"kz_bot/models"
	"kz_bot/storage"
	"os"
	"time"
)

type Percent struct {
	log     *logger.Logger
	storage *storage.Storage
	clients *clients.Clients
}

func NewPercent(log *logger.Logger, storage *storage.Storage, clients *clients.Clients) *Percent {
	return &Percent{log: log, storage: storage, clients: clients}
}

func (b *Percent) GetHadesStorage() {
	contentNew := b.getKeyAll()
	listHCorp := make(map[string]models.LevelCorp)

	all, err := b.storage.LevelCorp.ReadCorpLevelAll()
	if err != nil {
		b.log.ErrorErr(err)
		return
	}

	for _, corp := range all {
		listHCorp[corp.HCorp] = corp
	}

	var corpsdata []CorpsData

	for _, cont := range contentNew {
		data := b.getKey(cont.Key)
		corpsdata = append(corpsdata, CorpsData{
			Corp1Name:  data.Corporation1Name,
			Corp2Name:  data.Corporation2Name,
			Corp1Score: data.Corporation1Score,
			Corp2Score: data.Corporation2Score,
			DateEnded:  data.DateEnded,
		})

		if listHCorp[data.Corporation1Name].HCorp != "" {
			c := listHCorp[data.Corporation1Name]
			c.EndDate = data.DateEnded
			c.Percent = c.Level - 1
			fmt.Println(c)
			b.storage.LevelCorp.InsertUpdateCorpLevel(c)
		}

		time.Sleep(1000)
	}

	marshal, err := json.Marshal(corpsdata)
	if err != nil {
		return
	}
	err = os.WriteFile("ws/ws.json", marshal, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}

//func loadListCorp() []Corp {
//	var ListCorp []Corp
//	ListCorp = append(ListCorp, Corp{"Черный Легион", 1})
//	ListCorp = append(ListCorp, Corp{"украина№1", 1})
//	ListCorp = append(ListCorp, Corp{"русь ", 17})
//	ListCorp = append(ListCorp, Corp{"СССР", 1})
//	ListCorp = append(ListCorp, Corp{"Слава Украине!", 1})
//	ListCorp = append(ListCorp, Corp{"UKR Spase", 1})
//	ListCorp = append(ListCorp, Corp{"Октябристы", 1})
//	ListCorp = append(ListCorp, Corp{"UAGC", 1})
//	ListCorp = append(ListCorp, Corp{"Повстанцы Хаоса", 1})
//	ListCorp = append(ListCorp, Corp{"IX Легион", 1})
//	ListCorp = append(ListCorp, Corp{"-=Содружество=-", 1})
//	ListCorp = append(ListCorp, Corp{"Гарри Поттер", 1})
//	ListCorp = append(ListCorp, Corp{"Свободный флот", 1})
//
//	//ListCorp = append(ListCorp, Corp{"Омикрон Альфа",1})145
//	//ListCorp = append(ListCorp, "ВКС")209
//	//ListCorp = append(ListCorp, "Феникс")144
//	//ListCorp = append(ListCorp, "TFMC")
//	//ListCorp = append(ListCorp, "Северный флот")70
//	//ListCorp = append(ListCorp, "СвятыеНегодники").
//	//ListCorp = append(ListCorp, "ГОРИЗОНТ").
//	//ListCorp = append(ListCorp, "Торговая Федерация").
//	//ListCorp = append(ListCorp, "Zeta")
//
//	return ListCorp
//}

//"Sich.ua"
//"DarkMoon"
//"Odessa"

func (b *Percent) SendPercent(Config models.CorporationConfig) {
	currentCorp, err := b.storage.LevelCorp.ReadCorpLevel(Config.CorpName)
	if err != nil {
		b.log.ErrorErr(err)
		return
	}
	untilTime := currentCorp.EndDate.AddDate(0, 0, 7).Unix()
	if time.Now().UTC().Unix() < untilTime {
		return
	}

	all, err := b.storage.LevelCorp.ReadCorpLevelAll()
	if err != nil {
		b.log.ErrorErr(err)
		return
	}
	var preperingText []string
	for _, corp := range all {
		if corp.HCorp != "" && corp.Percent != 0 {
			untilTime = corp.EndDate.AddDate(0, 0, 7).Unix()
			if time.Now().UTC().Unix() < untilTime {
				preperingText = append(preperingText,
					fmt.Sprintf("%d%% %s %+v\n", percent(corp.Level), corp.HCorp, formatTime(untilTime)))
			}
		}
	}
	sortText := sortByFirstTwoDigits(preperingText)

	text := ""

	for _, s := range sortText {
		text += s
	}

	if Config.DsChannel != "" {
		go b.clients.Ds.SendChannelDelSecond(Config.DsChannel, text, 180)
	}
	if Config.TgChannel != "" {
		go b.clients.Tg.SendChannelDelSecond(Config.TgChannel, text, 180)
	}
}

func (b *Percent) GetTextPercent(Config models.CorporationConfig, dark bool) string {
	currentCorp, err := b.storage.LevelCorp.ReadCorpLevel(Config.CorpName)
	if err != nil {
		b.log.ErrorErr(err)
		return ""
	}
	untilTime := currentCorp.EndDate.AddDate(0, 0, 7).Unix()
	if time.Now().UTC().Unix() < untilTime {
		per := percent(currentCorp.Level)
		if dark {
			return fmt.Sprintf(" %d%%", per+100)
		} else {
			return fmt.Sprintf(" %d%%", per)
		}
	}
	return ""
}
