package storage

import "context"

type Count interface {
	СountName(ctx context.Context, userid, lvlkz, corpName string) (int, error)
	CountQueue(ctx context.Context, lvlkz, CorpName string) (int, error)
	CountNumberNameActive1(ctx context.Context, lvlkz, CorpName, userid string) (int, error)
	CountNameQueue(ctx context.Context, userid string) (countNames int)
	CountNameQueueCorp(ctx context.Context, userid, corp string) (countNames int)
	ReadTop5Level(corpname string) []string
}
