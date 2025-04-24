package postgres

import (
	"fmt"
	"github.com/google/uuid"
	"rs/models"
)

func (d *Db) EmojiReadUUID(uid uuid.UUID, tip string) models.Emoji {
	ctx, cancel := d.getContext()
	defer cancel()
	Emoji := "SELECT * FROM rs_bot.emoji WHERE uid = $1 AND tip = $2"
	results, err := d.db.Query(ctx, Emoji, uid, tip)
	defer results.Close()
	if err != nil {
		d.log.ErrorErr(err)
	}
	var t models.Emoji
	for results.Next() {
		err = results.Scan(&t.Uid, &t.Tip, &t.Em1, &t.Em2, &t.Em3, &t.Em4)
		if err != nil {
			d.log.ErrorErr(err)
		}
	}
	return t
}
func (d *Db) EmojiUpdateUUID(uid uuid.UUID, tip, slot, emo string) string {
	ctx, cancel := d.getContext()
	defer cancel()
	sqlUpd := fmt.Sprintf("update rs_bot.emoji set em%s = $1 where uid = $2 AND tip = $3", slot)
	_, err := d.db.Exec(ctx, sqlUpd, emo, uid, tip)
	if err != nil {
		d.log.ErrorErr(err)
	}
	return fmt.Sprintf("Слот %s обновлен\n%s", slot, emo)
}
func (d *Db) EmojiInsertEmptyUUID(uid uuid.UUID, tip string) {
	ctx, cancel := d.getContext()
	defer cancel()
	insert := `INSERT INTO rs_bot.emoji(uid,tip,em1,em2,em3,em4) VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := d.db.Exec(ctx, insert, uid, tip, "", "", "", "")
	if err != nil {
		d.log.ErrorErr(err)
	}
}
