package server

import (
	"compendium_s/config"
	"fmt"
	"time"

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

func GetTokenData(tokenString string) (uid, gid uuid.UUID, err error) {
	// Восстановление после паники
	defer func() {
		if r := recover(); r != nil {
			// Присваиваем ошибку в именованный возврат
			err = fmt.Errorf("panic recovered in GetTokenData: %v", r)

			// Опционально: логируем стек вызовов для отладки
			// log.Printf("Stack trace: %s", debug.Stack())
		}
	}()

	mapClaims, err := parseToken(tokenString)
	if err != nil {
		return
	}

	printClaims := func() {
		for s, a := range mapClaims {
			fmt.Printf("mapClaims: %s: %+v\n", s, a)
		}
	}

	// Эти строки могут вызвать панику (Type Assertion Panic)
	uid, err = uuid.Parse(mapClaims["uuid"].(string))
	if err != nil {
		printClaims()
		return
	}

	gid, err = uuid.Parse(mapClaims["gid"].(string))
	if err != nil {
		printClaims()
		return
	}

	return uid, gid, nil
}

func JWTGenerateToken(uid uuid.UUID, gid uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"uuid": uid,
		"gid":  gid,
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

func JWTGenerateTokenV2(uuid, gid uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"uuid": uuid,
		"gid":  gid,
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
