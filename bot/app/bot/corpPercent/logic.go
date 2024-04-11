package corpPercent

import (
	"fmt"
	"github.com/mentalisit/logger"
	"kz_bot/models"
	"kz_bot/storage"
	"time"
)

var log *logger.Logger

func GetHadesStorage(l *logger.Logger, st *storage.Storage) {
	log = l
	keys := getKeyAll()
	listHCorp := make(map[string]models.LevelCorp)

	all, err := st.LevelCorp.ReadCorpLevelAll()
	if err != nil {
		log.ErrorErr(err)
		return
	}

	for _, corp := range all {
		listHCorp[corp.HCorp] = corp
	}

	for _, key := range keys {
		data := getKey(key)

		if listHCorp[data.Corporation1Name].HCorp != "" {
			c := listHCorp[data.Corporation1Name]
			c.EndDate = data.DateEnded
			c.Percent = c.Level - 1
			fmt.Println(c)
			st.LevelCorp.InsertUpdateCorpLevel(c)
		}

		time.Sleep(1000)
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
