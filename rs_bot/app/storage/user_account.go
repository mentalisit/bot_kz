package storage

import "rs/models"

type UserAccount interface {
	UserAccountInsert(u models.UserAccount) error
	UserAccountUpdate(u models.UserAccount) error
	UserAccountGetByInternalUserId(IId string) (*models.UserAccount, error)
	UserAccountGetByTgUserId(TgId string) (*models.UserAccount, error)
	UserAccountGetByDsUserId(DsId string) (*models.UserAccount, error)
	UserAccountGetAll() ([]models.UserAccount, error)
	//FakeUserGetAll() ([]models.PlayerStats, error)
	//FakeUserInsert(userName string, points, level int) error
}
