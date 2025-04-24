package multi

import (
	"compendium/models"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"strings"
	"time"
)

func (d *Db) GenerateLinkCode(uid uuid.UUID) (*models.AccountLinkCode, error) {
	ctx, cancel := d.getContext()
	defer cancel()

	// Генерация кода
	code := strings.ToUpper(uuid.New().String()[:6])

	// Время истечения 10 минут в UTC
	expiresAt := time.Now().UTC().Add(10 * time.Minute)

	// Время создания кода в UTC
	createdAt := time.Now().UTC()

	// Вставка кода в базу данных
	_, err := d.db.Exec(ctx, `
		INSERT INTO compendium.accounts_link_codes (code, uuid, expires_at, created_at)
		VALUES ($1, $2, $3, $4)
	`, code, uid, expiresAt, createdAt)

	if err != nil {
		d.log.ErrorErr(err)
		return nil, err
	}

	// Возврат данных о созданном коде
	return &models.AccountLinkCode{
		Code:      code,
		UUID:      uid,
		ExpiresAt: expiresAt,
		CreatedAt: createdAt,
	}, nil
}

func (d *Db) ValidateLinkCode(code string) (*uuid.UUID, error) {
	ctx, cancel := d.getContext()
	defer cancel()

	// Проверка кода в базе данных
	var uid uuid.UUID
	var expiresAt time.Time

	// Запрос на получение UUID и времени истечения для кода
	row := d.db.QueryRow(ctx, `
        SELECT uuid, expires_at
        FROM compendium.accounts_link_codes
        WHERE code = $1
    `, code)

	err := row.Scan(&uid, &expiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("invalid or expired code")
		}
		return nil, err
	}

	// Проверка, не истек ли срок действия кода
	if time.Now().UTC().After(expiresAt) {
		return nil, fmt.Errorf("code has expired")
	}

	// Возврат UUID, если код валидный
	return &uid, nil
}

func (d *Db) DeleteLinkCodesByUUID(uid uuid.UUID) error {
	ctx, cancel := d.getContext()
	defer cancel()

	_, err := d.db.Exec(ctx, `
		DELETE FROM compendium.accounts_link_codes
		WHERE uuid = $1
	`, uid)

	if err != nil {
		return err
	}

	return nil
}
