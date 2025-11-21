package models

//type Button struct {
//	Text string
//	Data string
//}

//type Timer struct {
//	//Id       string `bson:"_id"`
//	Dsmesid  string `bson:"dsmesid"`
//	Dschatid string `bson:"dschatid"`
//	Tgmesid  string `bson:"tgmesid"`
//	Tgchatid string `bson:"tgchatid"`
//	Timed    int    `bson:"timed"`
//}

type Timer struct {
	//Id       string `bson:"_id"`
	Tip    string `bson:"tip"`
	ChatId string `bson:"chatId"`
	MesId  string `bson:"mesId"`
	Timed  int    `bson:"timed"`
}
