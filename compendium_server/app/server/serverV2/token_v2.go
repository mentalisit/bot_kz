package serverV2

import (
	"compendium_s/config"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func parseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.Instance.Postgress.Password), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("cannot parse claims")
	}

	return claims, nil
}

func GetTokenData(tokenString string) (uid uuid.UUID, err error) {
	mapClaims, err := parseToken(tokenString)
	if err != nil {
		return
	}
	uid, _ = uuid.Parse(mapClaims["uuid"].(string))

	return uid, nil
}
