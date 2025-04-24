package generate

import (
	"compendium/config"
	"github.com/google/uuid"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func JWTGenerateToken(uid uuid.UUID, GId uuid.UUID, NickName string) (string, error) {
	claims := jwt.MapClaims{
		"uuid": uid,
		"gid":  GId,
		"nick": NickName,
		"exp":  time.Now().AddDate(1, 0, 0).Unix(), // токен на год
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(config.Instance.Postgress.Password))
	if err != nil {
		return "", err
	}

	// добавляем префикс
	return "Multi_" + signedToken, nil
}
