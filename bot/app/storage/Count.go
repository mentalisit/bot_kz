package storage

type Count interface {
	СountName(userid, lvlkz, corpName string) (int, error)
	CountQueue(lvlkz, CorpName string) (int, error)
	CountNumberNameActive1(lvlkz, CorpName, userid string) (int, error)
	CountNameQueue(userid string) (countNames int)
	CountNameQueueCorp(userid, corp string) (countNames int)
	ReadTop5Level(corpname string) []string
}
