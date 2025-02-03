package models

import "strconv"

type UserAccount struct {
	InternalId  int
	GeneralName string
	TgId        string
	DsId        string
	GameId      []string
	ActiveName  string
	Accounts    []string
}

// GetAlt функция получения по номеру
// 0 возвращает GeneralName
// 1-5, возвращает альтов 0-4
// если i больше чем количество аккаунтов возвращаем GeneralName
func (u *UserAccount) GetAlt(i int) string {
	if i == 0 || i >= len(u.Accounts) {
		return u.GeneralName
	}
	return u.Accounts[i-1]
}
func (u *UserAccount) ContainsGameId(i int64) bool {
	s := strconv.FormatInt(i, 10)
	if len(u.GameId) == 0 {
		return false
	}
	for _, ss := range u.GameId {
		if ss == s {
			return true
		}
	}
	return false
}
