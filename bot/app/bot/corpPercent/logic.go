package corpPercent

import (
	"fmt"
	"github.com/mentalisit/logger"
	"time"
)

var log *logger.Logger

func GetHadesStorage(l *logger.Logger) {
	log = l
	keys := getKeyAll()

	for _, key := range keys {
		data := getKey(key)

		fmt.Println("\"" + data.Corporation1Name + "\"")

		//fmt.Printf("%d-%d %s-%s\n", data.Corporation1Score, data.Corporation2Score, data.Corporation1Name, data.Corporation2Name)

		//for _, corp := range list {
		//	if data.Corporation1Name == corp.Name {
		//		fmt.Println("\"" + data.Corporation1Name + "\"")
		//	}
		//}

		time.Sleep(1000)
	}
	//fmt.Println("keys:", keys)
}

//type Corp struct {
//	Name  string
//	Level int
//}
//
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
//	//ListCorp = append(ListCorp, "СвятыеНегодники")
//	//ListCorp = append(ListCorp, "ГОРИЗОНТ")
//	//ListCorp = append(ListCorp, "Торговая Федерация")
//	//ListCorp = append(ListCorp, "Zeta")
//
//	return ListCorp
//}

//"Sich.ua"
//"DarkMoon"
//"Odessa"
