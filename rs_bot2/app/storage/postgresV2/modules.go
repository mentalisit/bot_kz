package postgresV2

import (
	"encoding/json"
	"fmt"
	"rs/models"
	"time"

	"github.com/google/uuid"
)

//func (d *Db) ModuleReadUUID(uid uuid.UUID, name string) *models.Module {
//	cancel := d.getContext()
//	defer cancel()
//	module := "SELECT * FROM rs_bot.module WHERE uid = $1 AND name = $2"
//	results, err := d.db.Query(module, uid, name)
//	if err != nil {
//		d.log.ErrorErr(err)
//		return nil
//	}
//	defer results.Close()
//	var t models.Module
//	for results.Next() {
//		err = results.Scan(&t.Uid, &t.Name, &t.Gen, &t.Enr, &t.Rse)
//		if err != nil {
//			d.log.ErrorErr(err)
//		}
//	}
//	if t.Name == "" {
//		return nil
//	}
//	return &t
//}
//
//func (d *Db) ModuleReadAllUUID(uid uuid.UUID) []models.Module {
//	cancel := d.getContext()
//	defer cancel()
//
//	var mm []models.Module
//
//	module := "SELECT * FROM rs_bot.module WHERE uid = $1"
//
//	results, err := d.db.Query(module, uid)
//	if err != nil {
//		d.log.ErrorErr(err)
//		return mm
//	}
//	defer results.Close()
//
//	for results.Next() {
//		var t models.Module
//		err = results.Scan(&t.Uid, &t.Name, &t.Gen, &t.Enr, &t.Rse)
//		if err != nil {
//			d.log.ErrorErr(err)
//			continue
//		}
//		mm = append(mm, t)
//	}
//	return mm
//}
//
//func (d *Db) ModuleUpdateUUID(m models.Module) error {
//	go d.ModuleCompendiumUpdate(m)
//	cancel := d.getContext()
//	defer cancel()
//
//	sqlUpd := `update rs_bot.module set gen = $1, enr = $2, rse = $3 where uid = $4 AND name = $5`
//
//	// 1. Получаем результат выполнения
//	res, err := d.db.Exec(sqlUpd, m.Gen, m.Enr, m.Rse, m.Uid, m.Name)
//	if err != nil {
//		return err
//	}
//
//	// 2. Проверяем, сколько строк было изменено
//	rows := res.RowsAffected()
//	if rows == 0 {
//		// Данные не обновились, так как запись не найдена
//		return fmt.Errorf("module not found or no changes made")
//	}
//
//	return nil
//}
//
//func (d *Db) ModuleInsertUUID(m models.Module) {
//	cancel := d.getContext()
//	defer cancel()
//	insert := `INSERT INTO rs_bot.module(uid,name,gen,enr,rse) VALUES ($1,$2,$3,$4,$5)`
//	_, err := d.db.Exec(insert, m.Uid, m.Name, m.Gen, m.Enr, m.Rse)
//	if err != nil {
//		d.log.ErrorErr(err)
//	}
//}

func (d *Db) ModuleCompendiumGetAll(uid uuid.UUID) []models.Module {

	var techLevels []models.Module

	// Используем оператор -> для доступа к объекту и ->> для текста,
	// затем приводим к INTEGER ::int
	query := `
        SELECT 
            username, 
            COALESCE((tech->'508'->>'level')::int, 0) as gen,
            COALESCE((tech->'603'->>'level')::int, 0) as rse,
            COALESCE((tech->'503'->>'level')::int, 0) as enr
        FROM my_compendium.technologies 
        WHERE uid = $1`

	results, err := d.db.Query(query, uid)
	if err != nil {
		d.log.ErrorErr(err)
		return techLevels
	}
	defer results.Close()

	for results.Next() {
		var m models.Module
		m.Uid = uid

		// Теперь сканируем напрямую в поля структуры
		err = results.Scan(&m.Name, &m.Gen, &m.Rse, &m.Enr)
		if err != nil {
			d.log.ErrorErr(err)
			continue
		}
		techLevels = append(techLevels, m)
	}
	return techLevels
}

func (d *Db) ModuleCompendiumGet(uid uuid.UUID, nickname string) *models.Module {

	query := `
        SELECT 
            username, 
            COALESCE((tech->'508'->>'level')::int, 0) as gen,
            COALESCE((tech->'603'->>'level')::int, 0) as rse,
            COALESCE((tech->'503'->>'level')::int, 0) as enr
        FROM my_compendium.technologies 
        WHERE uid = $1 and username = $2`

	var m models.Module
	m.Uid = uid

	err := d.db.QueryRow(query, uid, nickname).Scan(&m.Name, &m.Gen, &m.Rse, &m.Enr)
	if err != nil {
		return nil
	}

	return &m
}

func (d *Db) ModuleCompendiumInsertUpdate(m models.Module) error {
	fmt.Printf("ModuleCompendiumInsertUpdate %+v\n", m)

	// Формируем "патч" в формате JSON
	patch := map[string]models.TechLevel{
		"508": {Level: m.Gen, Ts: time.Now().UTC().UnixMilli()},
		"603": {Level: m.Rse, Ts: time.Now().UTC().UnixMilli()},
		"503": {Level: m.Enr, Ts: time.Now().UTC().UnixMilli()},
	}

	patchJSON, err := json.Marshal(patch)
	if err != nil {
		return err
	}

	query := `INSERT INTO my_compendium.technologies (uid, username, tech)
		VALUES ($1, $2, $3::jsonb)
		ON CONFLICT (uid, username) 
		DO UPDATE SET tech = technologies.tech || EXCLUDED.tech;`

	_, err = d.db.Exec(query, m.Uid.String(), m.Name, patchJSON)
	return err
}

func (d *Db) ModuleCompendiumUpdateNickName(uid uuid.UUID, oldName, newName string) error {

	upd := `update my_compendium.technologies set username = $1 where username = $2 and uid = $3`

	_, err := d.db.Exec(upd, newName, oldName, uid.String())
	return err
}

func (d *Db) ModuleCompendiumCleanDeletedAlts(uid uuid.UUID, mainNickname string, currentAlts []string) error {

	// Собираем всех, кого НЕЛЬЗЯ удалять (основа + текущие альты)
	keepNames := append(currentAlts, mainNickname)

	// Удаляем из таблицы технологий тех, кого нет в списке keepNames
	query := `
        DELETE FROM my_compendium.technologies 
        WHERE uid = $1 AND username NOT IN (SELECT unnest($2::text[]))`

	_, err := d.db.Exec(query, uid, keepNames)
	return err
}
