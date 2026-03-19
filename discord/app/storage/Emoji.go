package storage

import (
	"github.com/mentalisit/restapi/models"
)

type Emoji interface {
	EmojiModuleReadUsers(name, tip string) (models.EmodjiUser, error)
	EmojiUpdate(name, tip, slot, emo string) string
	ModuleUpdate(name, tip, slot, moduleAndLevel string) string
	WeaponUpdate(name, tip, weapon string) string
	EmInsertEmpty(tip, name string)
}
