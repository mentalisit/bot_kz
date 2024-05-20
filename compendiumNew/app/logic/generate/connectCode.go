package generate

import (
	"compendium/models"
	"math/rand"
	"strings"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GenerateFormattedString(Identity models.Identity) string {
	//rand
	segments := []string{randString(4), randString(4), randString(4)}
	code := strings.Join(segments, "-")
	go saveCode(code, Identity)
	return code
}
func saveCode(code string, Identity models.Identity) {
	s := codeStruct{
		code:     code,
		time:     time.Now(),
		Identity: Identity,
	}
	codes = append(codes, s)
}

type codeStruct struct {
	code     string
	time     time.Time
	Identity models.Identity
}

var codes []codeStruct

func CheckCode(CheckCode string) models.Identity {
	var i models.Identity
	if len(codes) > 0 {
		var newcodes []codeStruct
		for _, code := range codes {
			if time.Since(code.time) <= 5*time.Minute {
				if code.code == CheckCode {
					i = code.Identity
					newcodes = append(newcodes, code) // это временно
				} else {
					newcodes = append(newcodes, code)
				}
			}
		}
		codes = newcodes
	}
	if i.User.Username == "" {
		i = CheckCodeTest(CheckCode)
	}
	return i
}

func CheckCodeTest(CheckCode string) models.Identity {
	if CheckCode == "test-test-test" {
		return models.Identity{
			User: models.User{
				ID:       "111111111",
				Username: "TestUser",
				Alts:     []string{"alt1", "alt2"},
			},
			Guild: models.Guild{
				ID:   "22222222222",
				Name: "TestGuild",
				Type: "tg",
			},
			Token: "gGUBIlUAU1uTKWd8HssP27ojG0DugoAaPslwFGTDSAbEM6UM",
		}
	}
	return models.Identity{}
}
