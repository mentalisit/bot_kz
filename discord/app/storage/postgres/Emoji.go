package postgres

import (
	"discord/models"
	"fmt"
)

func (d *Db) EmojiModuleReadUsers(name, tip string) (models.EmodjiUser, error) {
	ctx, cancel := d.getContext()
	defer cancel()
	selec := "SELECT * FROM kzbot.users WHERE name = $1 AND tip = $2"
	results, err := d.db.Query(ctx, selec, name, tip)
	defer results.Close()
	if err != nil {
		return models.EmodjiUser{}, err
	}
	var t models.EmodjiUser
	for results.Next() {
		err = results.Scan(&t.Id, &t.Tip, &t.Name, &t.Em1, &t.Em2, &t.Em3, &t.Em4, &t.Module1, &t.Module2, &t.Module3, &t.Weapon)
		if err != nil {
			return models.EmodjiUser{}, err
		}
	}
	return t, nil
}
func (d *Db) EmojiUpdate(name, tip, slot, emo string) string {
	ctx, cancel := d.getContext()
	defer cancel()
	sqlUpd := fmt.Sprintf(`update kzbot.users set em%s = $1 where name = $2 AND tip = $3`, slot)
	_, err := d.db.Exec(ctx, sqlUpd, emo, name, tip)
	if err != nil {
		d.log.ErrorErr(err)
	}
	return fmt.Sprintf("Слот %s обновлен\n%s", slot, emo)
}
func (d *Db) ModuleUpdate(name, tip, slot, moduleAndLevel string) string {
	ctx, cancel := d.getContext()
	defer cancel()
	sqlUpd := fmt.Sprintf(`update kzbot.users set module%s = $1 where name = $2 AND tip = $3`, slot)
	_, err := d.db.Exec(ctx, sqlUpd, moduleAndLevel, name, tip)
	if err != nil {
		d.log.ErrorErr(err)
	}
	return fmt.Sprintf("Модуль %s обновлен\n%s", slot, moduleAndLevel)
}
func (d *Db) WeaponUpdate(name, tip, weapon string) string {
	ctx, cancel := d.getContext()
	defer cancel()
	sqlUpd := `update kzbot.users set weapon = $1 where name = $2 AND tip = $3`
	_, err := d.db.Exec(ctx, sqlUpd, weapon, name, tip)
	if err != nil {
		d.log.ErrorErr(err)
	}
	return fmt.Sprintf("Оружие обновлено\n%s", weapon)
}
func (d *Db) EmInsertEmpty(tip, name string) {
	ctx, cancel := d.getContext()
	defer cancel()
	insert := `INSERT INTO kzbot.users(tip,name,em1,em2,em3,em4,module1,module2,module3,weapon) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := d.db.Exec(ctx, insert, tip, name, "", "", "", "", "", "", "", "")
	if err != nil {
		d.log.ErrorErr(err)
	}
}
