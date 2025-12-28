package imageGenerator

import (
	"bytes"
	"compendium/models"
	"fmt"
	"image"
	"image/gif"
	_ "image/jpeg"
	"image/png"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

var user models.TechLevels

func GetLevel(i int) string {
	if user[i].Level > 0 {
		return strconv.Itoa(user[i].Level)
	}
	return ""
}
func addModulesLevel(dc *gg.Context) {
	y := float64(305)
	x := float64(100)
	x2 := x + 90
	x3 := x2 + 90
	x4 := x3 + 90
	x5 := x4 + 185
	x6 := x5 + 90
	x7 := x6 + 90
	x8 := x7 + 90
	//transport1
	dc.DrawStringAnchored(GetLevel(401), x, y, 0, 0.5) //90
	dc.DrawStringAnchored(GetLevel(402), x2, y, 0, 0.5)
	dc.DrawStringAnchored(GetLevel(413), x3, y, 0, 0.5)
	dc.DrawStringAnchored(GetLevel(404), x4, y, 0, 0.5)
	//mainer1
	dc.DrawStringAnchored(GetLevel(501), x5, y, 0, 0.5)
	dc.DrawStringAnchored(GetLevel(511), x6, y, 0, 0.5) ////
	dc.DrawStringAnchored(GetLevel(512), x7, y, 0, 0.5)
	dc.DrawStringAnchored(GetLevel(504), x8, y, 0, 0.5)
	///
	y += 70
	//transport2
	dc.DrawStringAnchored(GetLevel(608), x, y, 0, 0.5) //90
	dc.DrawStringAnchored(GetLevel(405), x2, y, 0, 0.5)
	dc.DrawStringAnchored(GetLevel(406), x3, y, 0, 0.5)
	dc.DrawStringAnchored(GetLevel(603), x4, y, 0, 0.5)
	//mainer2
	dc.DrawStringAnchored(GetLevel(508), x5, y, 0, 0.5)
	dc.DrawStringAnchored(GetLevel(503), x6, y, 0, 0.5)
	dc.DrawStringAnchored(GetLevel(507), x7, y, 0, 0.5)
	dc.DrawStringAnchored(GetLevel(505), x8, y, 0, 0.5)
	///
	y += 70
	//transport3
	dc.DrawStringAnchored(GetLevel(412), x, y, 0, 0.5) //90
	dc.DrawStringAnchored(GetLevel(411), x2, y, 0, 0.5)
	dc.DrawStringAnchored(GetLevel(414), x3, y, 0, 0.5)
	//mainer3
	dc.DrawStringAnchored(GetLevel(510), x5, y, 0, 0.5) ////
	dc.DrawStringAnchored(GetLevel(513), x6, y, 0, 0.5) ////

	////
	y += 110
	//weapon1
	dc.DrawStringAnchored(GetLevel(203), x, y, 0, 0.5) //90
	dc.DrawStringAnchored(GetLevel(204), x2, y, 0, 0.5)
	dc.DrawStringAnchored(GetLevel(202), x3, y, 0, 0.5)
	dc.DrawStringAnchored(GetLevel(205), x4, y, 0, 0.5)
	//Shield1
	dc.DrawStringAnchored(GetLevel(301), x5, y, 0, 0.5)
	dc.DrawStringAnchored(GetLevel(302), x6, y, 0, 0.5)
	dc.DrawStringAnchored(GetLevel(303), x7, y, 0, 0.5)
	dc.DrawStringAnchored(GetLevel(304), x8, y, 0, 0.5)
	///
	y += 70
	//weapon2
	dc.DrawStringAnchored(GetLevel(206), x, y, 0, 0.5) //90
	dc.DrawStringAnchored(GetLevel(207), x2, y, 0, 0.5)
	dc.DrawStringAnchored(GetLevel(208), x3, y, 0, 0.5)
	dc.DrawStringAnchored(GetLevel(209), x4, y, 0, 0.5)
	//Shield2
	dc.DrawStringAnchored(GetLevel(306), x5, y, 0, 0.5)
	dc.DrawStringAnchored(GetLevel(305), x6, y, 0, 0.5)
	dc.DrawStringAnchored(GetLevel(307), x7, y, 0, 0.5)
	dc.DrawStringAnchored(GetLevel(308), x8, y, 0, 0.5)
	y += 70
	//weapon3
	dc.DrawStringAnchored(GetLevel(210), x, y, 0, 0.5)
	////
	y += 110
	//support1
	dc.DrawStringAnchored(GetLevel(601), x, y, 0, 0.5)  //90
	dc.DrawStringAnchored(GetLevel(625), x2, y, 0, 0.5) /////
	dc.DrawStringAnchored(GetLevel(609), x3, y, 0, 0.5)
	dc.DrawStringAnchored(GetLevel(602), x4, y, 0, 0.5)
	//drone1
	dc.DrawStringAnchored(GetLevel(901), x5, y, 0, 0.5)
	dc.DrawStringAnchored(GetLevel(902), x6, y, 0, 0.5)
	dc.DrawStringAnchored(GetLevel(904), x7, y, 0, 0.5)
	dc.DrawStringAnchored(GetLevel(624), x8, y, 0, 0.5)
	///
	y += 70
	//support2
	dc.DrawStringAnchored(GetLevel(626), x, y, 0, 0.5) //90
	dc.DrawStringAnchored(GetLevel(614), x2, y, 0, 0.5)
	dc.DrawStringAnchored(GetLevel(615), x3, y, 0, 0.5)
	dc.DrawStringAnchored(GetLevel(616), x4, y, 0, 0.5)
	//drone2
	dc.DrawStringAnchored(GetLevel(908), x5, y, 0, 0.5)
	dc.DrawStringAnchored(GetLevel(905), x6, y, 0, 0.5) ////
	dc.DrawStringAnchored(GetLevel(906), x7, y, 0, 0.5)
	dc.DrawStringAnchored(GetLevel(907), x8, y, 0, 0.5)
	///
	y += 70
	//support3
	dc.DrawStringAnchored(GetLevel(617), x, y, 0, 0.5) //90
	dc.DrawStringAnchored(GetLevel(618), x2, y, 0, 0.5)
	dc.DrawStringAnchored(GetLevel(619), x3, y, 0, 0.5)
	dc.DrawStringAnchored(GetLevel(622), x4, y, 0, 0.5) ///////?
	y += 70
	//support4
	dc.DrawStringAnchored(GetLevel(621), x, y, 0, 0.5) //90
	dc.DrawStringAnchored(GetLevel(623), x2, y, 0, 0.5)
	//level
	dc.DrawStringAnchored(GetLevel(101), x5, y, 0, 0.5)
	dc.DrawStringAnchored(GetLevel(103), x6, y, 0, 0.5)
	dc.DrawStringAnchored(GetLevel(102), x7, y, 0, 0.5)
	dc.DrawStringAnchored(GetLevel(701), x8, y, 0, 0.5)
}

func addAvatars(dc *gg.Context, avatarURL string, centerX, centerY int) {
	var img image.Image
	var imgType string
	//var err error

	// Если URL начинается с определенного префикса, открываем изображение локально
	if strings.HasPrefix(avatarURL, "https://compendiumnew.mentalisit.myds.me/compendium/avatars") {
		// Преобразуем URL в локальный путь
		localPath := strings.Replace(avatarURL, "https://compendiumnew.mentalisit.myds.me/", "", 1)
		localPath = "docker/" + localPath

		// Открываем локальный файл
		file, err := os.Open(localPath)
		if err != nil {
			fmt.Println("Error opening local avatar:", err)
			return
		}
		defer file.Close()

		// Определяем тип изображения
		img, imgType, err = image.Decode(file)
		if err != nil {
			fmt.Println("Error decoding local avatar:", err)
			return
		}
	} else {
		// Загружаем изображение по URL
		response, err := http.Get(avatarURL)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer response.Body.Close()

		// Определяем тип изображения
		img, imgType, err = image.Decode(response.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// Если это GIF, обрабатываем его иначе
	if imgType == "gif" {
		// Открываем изображение повторно для корректной обработки GIF
		response, err := http.Get(avatarURL)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer response.Body.Close()

		gifImg, err := gif.DecodeAll(response.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Используем первый кадр GIF для наложения
		img = gifImg.Image[0]
	}

	// Изменяем размер изображения аватара
	img = resize.Resize(185, 185, img, resize.Lanczos3)

	// Определяем радиус круга
	radius := 270 / 3

	// Рисуем круг с изображением
	dc.ResetClip()
	dc.DrawCircle(float64(centerX), float64(centerY), float64(radius))
	dc.Clip()
	dc.DrawImageAnchored(img, centerX, centerY, 0.5, 0.5)
}
func imageToBytes(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
