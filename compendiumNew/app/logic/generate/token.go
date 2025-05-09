package generate

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateToken(size int) string {
	// Вычисляем необходимый размер байт для указанной длины токена
	tokenBytes := make([]byte, size)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return ""
	}

	// Кодируем байты в строку base64
	token := base64.URLEncoding.EncodeToString(tokenBytes)
	return token
}
