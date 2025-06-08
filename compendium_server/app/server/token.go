package server

import (
	"compendium_s/config"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

func GenerateToken() string {
	// Вычисляем необходимый размер байт для указанной длины токена
	tokenBytes := make([]byte, 174)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return ""
	}

	// Кодируем байты в строку base64
	token := base64.URLEncoding.EncodeToString(tokenBytes)
	return token
}

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

func GetTokenData(tokenString string) (uid, gid uuid.UUID, err error) {
	mapClaims, err := parseToken(tokenString)
	if err != nil {
		return
	}
	uid, _ = uuid.Parse(mapClaims["uuid"].(string))
	gid, _ = uuid.Parse(mapClaims["gid"].(string))

	return uid, gid, nil
}

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

func GetTokenUserData(tokenString string) (userid string, gid uuid.UUID, err error) {
	mapClaims, err := parseToken(tokenString)
	if err != nil {
		return
	}
	userid = mapClaims["userId"].(string)
	gid, _ = uuid.Parse(mapClaims["multiGuildId"].(string))
	return userid, gid, nil
}
