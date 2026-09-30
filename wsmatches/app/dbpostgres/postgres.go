package dbpostgres

import (
	"sync"
	"ws/models"

	"github.com/jmoiron/sqlx"
	"github.com/mentalisit/conf/logger"

	_ "github.com/lib/pq"
)

type Db struct {
	log   *logger.Logger
	pool  *sqlx.DB
	cache map[string]models.CorporationsData
	mu    sync.RWMutex
}

func NewDb(log *logger.Logger, db *sqlx.DB) *Db {

	d := &Db{
		log:   log,
		pool:  db,
		cache: make(map[string]models.CorporationsData),
	}

	d.createTable()
	d.LoadAllData() // Загружаем данные в кэш

	return d
}

func (d *Db) createTable() {
	d.pool.Exec("CREATE SCHEMA IF NOT EXISTS ws")
	// Создание таблиц

	query := `
	CREATE TABLE IF NOT EXISTS ws.corporations (
		id TEXT PRIMARY KEY,
		data JSONB NOT NULL
	);`
	_, err := d.pool.Exec(query)
	if err != nil {
		d.log.ErrorErr(err)
	}

	query = `
	CREATE TABLE IF NOT EXISTS ws.corps (
		id TEXT PRIMARY KEY,
		data JSONB NOT NULL
	);`
	_, err = d.pool.Exec(query)
	if err != nil {
		d.log.ErrorErr(err)
	}

	_, err = d.pool.Exec(
		`CREATE TABLE IF NOT EXISTS ws.corpsLevel (
            corpName       text,
            level     	   integer,
            endDate        text,
            hCorp    	   text,
            percent    	   integer,
            last_update    text,
            relic          integer
        );
    `)
	if err != nil {
		d.log.ErrorErr(err)
		return
	}
}
