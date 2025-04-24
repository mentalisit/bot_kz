package dbpostgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mentalisit/logger"
	"os"
	"sync"
	"time"
	"ws/models"
)

type Db struct {
	log   *logger.Logger
	pool  *pgxpool.Pool
	cache map[string]models.CorporationsData
	//cacheCorp map[string]models.Corporation
	mu sync.RWMutex
}

func NewDb(log *logger.Logger, pass string) *Db {
	dns := fmt.Sprintf("postgres://postgres:%s@postgres:5432/postgres", pass)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dns)
	if err != nil {
		log.ErrorErr(err)
		os.Exit(1)
	}

	d := &Db{
		log:   log,
		pool:  pool,
		cache: make(map[string]models.CorporationsData),
		//cacheCorp: make(map[string]models.Corporation),
	}

	d.createTable()
	d.LoadAllData() // Загружаем данные в кэш

	return d
}

func (d *Db) getContext() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
func (d *Db) createTable() {
	ctx, cancel := d.getContext()
	defer cancel()
	d.pool.Exec(ctx, "CREATE SCHEMA IF NOT EXISTS ws")
	// Создание таблиц

	ctx, cancel = d.getContext()
	defer cancel()

	query := `
	CREATE TABLE IF NOT EXISTS ws.corporations (
		id TEXT PRIMARY KEY,
		data JSONB NOT NULL
	);`
	_, err := d.pool.Exec(ctx, query)
	if err != nil {
		d.log.ErrorErr(err)
	}

	ctx, cancel = d.getContext()
	defer cancel()

	query = `
	CREATE TABLE IF NOT EXISTS ws.corps (
		id TEXT PRIMARY KEY,
		data JSONB NOT NULL
	);`
	_, err = d.pool.Exec(ctx, query)
	if err != nil {
		d.log.ErrorErr(err)
	}

	_, err = d.pool.Exec(ctx,
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
