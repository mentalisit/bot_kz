package postgresV2

import (
	_ "github.com/lib/pq"
	"github.com/mentalisit/conf/config"
	"github.com/mentalisit/conf/logger"

	"github.com/jmoiron/sqlx"
)

type Db struct {
	db  *sqlx.DB
	log *logger.Logger
	Dns string
}

func NewDb(log *logger.Logger, db *sqlx.DB) *Db {
	d := &Db{
		db:  db,
		log: log,
		Dns: config.Instance.GetDNS(),
	}
	//d.InitChatBridge()
	d.InitLinkCodes()

	return d
}

func (d *Db) Shutdown() {
	d.db.Close()
}

//CREATE TABLE my_compendium.study (
//"uid" UUID NOT NULL,
//"name" TEXT NOT NULL,
//"studies" JSONB NOT NULL DEFAULT '[]',
//-- Теперь уникальность проверяется по связке обоих полей
//PRIMARY KEY ("uid", "name")
//);
