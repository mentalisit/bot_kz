package NewLogic

import (
	"compendium/config"
	"compendium/logic/generate"
	"compendium/models"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (c *HsLogic) GenerateIdentity(m models.IncomingMessage) models.IdentityV2 {
	var i models.IdentityV2
	i.GuildId = m.Guild.GId.String()
	i.MultiAccount = m.Acc
	i.Token, _ = JWTGenerateToken(i.MultiAccount.UUID)
	cm, err := c.db.DBv2.CorpMemberByUId(i.MultiAccount.UUID)
	if err != nil || cm == nil {
		cm = &models.MultiAccountCorpMember{
			Uid:        i.MultiAccount.UUID,
			GuildIds:   []uuid.UUID{m.Guild.GId},
			TimeZona:   "",
			ZonaOffset: 0,
			AfkFor:     "",
		}
		err = c.db.DBv2.CorpMemberInsert(*cm)
		if err != nil {
			c.log.ErrorErr(err)
		}
	}

	// проверить содержит ли GuildIds текущую гильдию, если нет - добавить и обновить
	contains := false
	for _, gid := range cm.GuildIds {
		if gid == m.Guild.GId {
			contains = true
			break
		}
	}
	if !contains {
		cm.GuildIds = append(cm.GuildIds, m.Guild.GId)
		err = c.db.DBv2.CorpMemberUpdate(*cm)
		if err != nil {
			c.log.ErrorErr(err)
		}
	}

	return i
}

func JWTGenerateToken(uid uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"uuid": uid,
		"exp":  time.Now().AddDate(1, 0, 0).Unix(), // токен на год
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(config.Instance.Postgress.Password))
	if err != nil {
		return "", err
	}

	// добавляем префикс
	return "my_compendium_" + signedToken, nil
}

func (c *HsLogic) generateCodeAndSave(Identity models.IdentityV2) string {
	segments := []string{generate.RandString(4), generate.RandString(4), generate.RandString(4)}

	m := models.CodeV2{
		Code:      strings.Join(segments, "-"),
		Timestamp: time.Now().Unix(),
		Identity:  Identity,
	}

	go func() {
		err := c.db.DBv2.CodeInsert(m)
		if err != nil {
			c.log.ErrorErr(err)
			c.log.InfoStruct("CodeInsert", m)
		}
	}()

	return m.Code
}
