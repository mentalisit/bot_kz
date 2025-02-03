package postgres

import (
	"fmt"
	"github.com/jackc/pgx/v5"
	"rs/models"
)

func (d *Db) UserAccountInsert(u models.UserAccount) error {
	ctx, cancel := d.GetContext()
	defer cancel()
	insert := `INSERT INTO rs_bot.user_account(general_name, user_id_tg, user_id_ds, user_id_game, user_name_active, user_accounts) 
       VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := d.db.Exec(ctx, insert, u.GeneralName, u.TgId, u.DsId, u.GameId, u.ActiveName, u.Accounts)
	if err != nil {
		return err
	}
	return nil
}
func (d *Db) UserAccountUpdate(u models.UserAccount) error {
	ctx, cancel := d.GetContext()
	defer cancel()
	upd := `update rs_bot.user_account set user_id_tg = $1, user_id_ds = $2, user_id_game = $3, user_name_active = $4 where user_accounts = $5`
	_, err := d.db.Exec(ctx, upd, u.TgId, u.DsId, u.GameId, u.ActiveName, u.Accounts)
	if err != nil {
		return err
	}
	return nil
}

func (d *Db) UserAccountGetByInternalUserId(IId string) (*models.UserAccount, error) {
	ctx, cancel := d.GetContext()
	defer cancel()
	var u models.UserAccount
	selectUser := "SELECT * FROM rs_bot.user_account WHERE internal_user_id = $1 "
	err := d.db.QueryRow(ctx, selectUser, IId).Scan(&u.InternalId, &u.GeneralName, &u.TgId, &u.DsId, &u.GameId, &u.ActiveName, &u.Accounts)
	if err != nil {
		return nil, pgx.ErrNoRows
	}
	return &u, nil
}
func (d *Db) UserAccountGetByTgUserId(TgId string) (*models.UserAccount, error) {
	ctx, cancel := d.GetContext()
	defer cancel()
	var u models.UserAccount
	selectUser := "SELECT * FROM rs_bot.user_account WHERE user_id_tg = $1 "
	err := d.db.QueryRow(ctx, selectUser, TgId).Scan(&u.InternalId, &u.GeneralName, &u.TgId, &u.DsId, &u.GameId, &u.ActiveName, &u.Accounts)
	if err != nil {
		return nil, pgx.ErrNoRows
	}
	return &u, nil
}
func (d *Db) UserAccountGetByDsUserId(DsId string) (*models.UserAccount, error) {
	ctx, cancel := d.GetContext()
	defer cancel()
	var u models.UserAccount
	selectUser := "SELECT * FROM rs_bot.user_account WHERE user_id_ds = $1 "
	err := d.db.QueryRow(ctx, selectUser, DsId).Scan(&u.InternalId, &u.GeneralName, &u.TgId, &u.DsId, &u.GameId, &u.ActiveName, &u.Accounts)
	if err != nil {
		return nil, pgx.ErrNoRows
	}
	return &u, nil
}

func (d *Db) UserAccountGetByGeneralName(generalName string) (*models.UserAccount, error) {
	ctx, cancel := d.GetContext()
	defer cancel()
	var u models.UserAccount
	selectUser := "SELECT * FROM rs_bot.user_account WHERE general_name = $1 "
	row := d.db.QueryRow(ctx, selectUser, generalName)
	err := row.Scan(&u.InternalId, &u.GeneralName, &u.TgId, &u.DsId, &u.GameId, &u.ActiveName, &u.Accounts)
	if err != nil {
		return nil, pgx.ErrNoRows
	}
	return &u, nil
}

func (d *Db) UserAccountGetAll() ([]models.UserAccount, error) {
	ctx, cancel := d.GetContext()
	defer cancel()
	var users []models.UserAccount
	selectUsers := `SELECT * FROM rs_bot.user_account`

	results, err := d.db.Query(ctx, selectUsers)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer results.Close()

	for results.Next() {
		var u models.UserAccount
		err = results.Scan(
			&u.InternalId,
			&u.GeneralName,
			&u.TgId,
			&u.DsId,
			&u.GameId,
			&u.ActiveName,
			&u.Accounts,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}

		users = append(users, u)
	}

	// Проверка ошибок после цикла
	if err = results.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке результатов: %w", err)
	}

	return users, nil
}
