package imageGenerator

import (
	"fmt"
	"github.com/fogleman/gg"
	"image/color"
	"log"
)

func GenerateUser(avatarURL, corpAvararUrl, nikName, corporation string, tech map[int][2]int) []byte {
	user = tech
	// Открываем изображение
	im, err := gg.LoadPNG("compendium/original2.png")
	if err != nil {
		fmt.Println(err)
	}

	// Создаем новый контекст для рисования
	dc := gg.NewContextForImage(im)

	// Устанавливаем параметры шрифта
	err = dc.LoadFontFace("compendium/font.ttf", 32)
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
