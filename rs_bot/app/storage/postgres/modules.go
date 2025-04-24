package postgres

import (
	"github.com/google/uuid"
	"rs/models"
	"rs/pkg/utils"
)

func (d *Db) ModuleReadUUID(uid uuid.UUID, name string) *models.Module {
	ctx, cancel := d.getContext()
	defer cancel()
	module := "SELECT * FROM rs_bot.module WHERE uid = $1 AND name = $2"
	results, err := d.db.Query(ctx, module, uid, name)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
	}
	var t models.Module
	for results.Next() {
		err = results.Scan(&t.Uid, &t.Name, &t.Gen, &t.Enr, &t.Rse)
		if err != nil {
			d.log.ErrorErr(err)
		}
	}
	return &t
}

// fmt.Sprintf("Модуль %s обновлен\n%s", slot, moduleAndLevel)
func (d *Db) ModuleUpdateUUID(m models.Module) {
	ctx, cancel := d.getContext()
	defer cancel()
	ch := utils.WaitForMessage("ModuleUpdate")
	defer close(ch)
	sqlUpd := `update rs_bot.module set gen = $1, enr = $2, rse = $3 where uid = $4 AND name = $5`
	_, err := d.db.Exec(ctx, sqlUpd, m.Gen, m.Enr, m.Rse, m.Uid, m.Name)
	if err != nil {
		d.log.ErrorErr(err)
	}
}

func (d *Db) ModuleInsertUUID(m models.Module) {
	ctx, cancel := d.getContext()
	defer cancel()
	insert := `INSERT INTO rs_bot.module(uid,name,gen,enr,rse) VALUES ($1,$2,$3,$4,$5)`
	_, err := d.db.Exec(ctx, insert, m.Uid, m.Name, m.Gen, m.Enr, m.Rse)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
