package imageGenerator

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/fogleman/gg"
)

func GenerateUser(avatarURL, corpAvararUrl, nikName, corporation string, tech map[int][2]int) []byte {
	user = tech
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(20) + 1

	im, err := gg.LoadPNG(fmt.Sprintf("docker/compendium/template/%d.png", randomNumber))
	if err != nil {
		fmt.Println(err)
	}

	// Создаем новый контекст для рисования
	dc := gg.NewContextForImage(im)

	// Устанавливаем параметры шрифта
	err = dc.LoadFontFace("docker/compendium/font.ttf", 32)
	if err != nil {
		fmt.Println(err)
	}

	// Устанавливаем цвет текста
	dc.SetColor(color.White)
	dc.DrawStringAnchored(nikName, 450, 100, 0.5, 0.5)
	dc.DrawStringAnchored(corporation, 450, 150, 0.5, 0.5)

	// Рисуем текст на изображении 1
	addModulesLevel(dc)
	if avatarURL != "" {
		addAvatars(dc, avatarURL, 125, 125)
	}
	if corpAvararUrl != "" {
		addAvatars(dc, corpAvararUrl, 760, 125)
	}

	reader, err := imageToBytes(dc.Image())
	if err != nil {
		log.Fatal(err)
	}
	return reader
}
