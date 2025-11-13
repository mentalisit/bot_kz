package helpers

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"rs/models"
	"strconv"

	"github.com/fogleman/gg"
)

func PicStatistic(name string, statistics []models.Statistic) []byte {
	const imageWidth = 250.0 // Фиксированная ширина
	const fontSize = 14.0
	const lineSpacing = 1.4
	const marginX = 20.0
	const marginY = 20.0
	const rowHeight = fontSize * lineSpacing

	const wEvent = 40.0
	const wLevel = 50.0
	const wRuns = 50.0
	const wPoints = 80.0

	headerHeight := marginY + (rowHeight * 1) + (rowHeight * 1) + (rowHeight * 1)
	headerRowHeight := rowHeight
	dataRowsHeight := float64(len(statistics)) * rowHeight

	// Общая высота: Заголовок + Строка заголовков + Строки данных + Нижний отступ
	totalHeight := headerHeight + headerRowHeight + dataRowsHeight + marginY

	imageHeight := totalHeight

	dc := gg.NewContext(int(imageWidth), int(imageHeight))

	// Основной фон - темный
	darkBackground := color.RGBA{R: 30, G: 30, B: 30, A: 255}
	dc.SetColor(darkBackground)
	dc.Clear()

	fontPath := "docker/compendium/NotoSans.ttc"
	if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
		fmt.Println("Ошибка загрузки шрифта, используется шрифт по умолчанию:", err)
		// Если нужно, можете тут загрузить другой шрифт
	}

	// Определения цветов
	white := color.RGBA{R: 186, G: 186, B: 186, A: 255}
	grayText := color.RGBA{R: 200, G: 200, B: 200, A: 255}
	headerBg := color.RGBA{R: 50, G: 50, B: 50, A: 255}
	rowBg := color.RGBA{R: 40, G: 40, B: 40, A: 255}

	greenHighlight := color.RGBA{R: 30, G: 224, B: 0, A: 255}
	redHighlight := color.RGBA{R: 224, G: 37, B: 9, A: 255}

	// Координаты начала рисования
	x := marginX
	y := marginY

	// 3. Заголовок таблицы - "Статистика игрока" и имя в новой строке
	dc.LoadFontFace(fontPath, fontSize*1.2)
	dc.SetColor(white)
	dc.DrawString("Статистика игрока", x, y)
	y += rowHeight * 1.2

	dc.SetColor(white)
	dc.DrawString(name, x, y)
	y += rowHeight * 1.5

	// Восстановление размера шрифта для таблицы
	dc.LoadFontFace(fontPath, fontSize)

	// 4. Определение максимальных и минимальных значений для Runs и Points
	maxRuns := -1
	minRuns := 999999
	maxPoints := -1
	minPoints := 999999

	for _, s := range statistics {
		if s.Runs > maxRuns {
			maxRuns = s.Runs
		}
		if s.Runs < minRuns {
			minRuns = s.Runs
		}
		if s.Points > maxPoints {
			maxPoints = s.Points
		}
		if s.Points < minPoints {
			minPoints = s.Points
		}
	}

	// Рисование фона для заголовка столбцов
	headerY := y - rowHeight + 5
	dc.SetColor(headerBg)
	dc.DrawRectangle(marginX, headerY, imageWidth-2*marginX, rowHeight)
	dc.Fill()

	// Рисование заголовка столбцов (текст)
	dc.SetColor(white)

	colX := x
	dc.SetColor(color.RGBA{R: 84, G: 103, B: 241, A: 255})
	dc.DrawString("##", colX, y)
	colX += wEvent
	dc.DrawString("DRS", colX, y)
	colX += wLevel
	dc.DrawString("Игры", colX, y)
	colX += wRuns
	dc.DrawString("Очки", colX, y)

	y += rowHeight // Сдвиг к первой строке данных

	// 5. Цикл для вывода данных с выделением Max/Min
	for i, s := range statistics {
		// Опционально: чередование цвета фона строк
		if i%2 == 1 {
			rowY := y - rowHeight + 5
			dc.SetColor(rowBg)
			dc.DrawRectangle(marginX, rowY, imageWidth-2*marginX, rowHeight)
			dc.Fill()
		}

		colX = x

		// 5.1. Рисование столбца "Ивент"
		dc.SetColor(grayText)
		dc.DrawString(strconv.Itoa(s.EventId), colX, y)
		colX += wEvent

		// 5.2. Рисование столбца "Уровень"
		dc.SetColor(grayText)
		dc.DrawString(fmt.Sprintf("%d", s.Level), colX, y)
		colX += wLevel

		// 5.3. Рисование столбца "Игры" с выделением Max/Min
		if s.Runs == maxRuns {
			dc.SetColor(greenHighlight)
		} else if s.Runs == minRuns {
			dc.SetColor(redHighlight)
		} else {
			dc.SetColor(grayText)
		}
		dc.DrawString(fmt.Sprintf("%d", s.Runs), colX, y)
		colX += wRuns

		// 5.4. Рисование столбца "Очки" с выделением Max/Min
		if s.Points == maxPoints {
			dc.SetColor(greenHighlight)
		} else if s.Points == minPoints {
			dc.SetColor(redHighlight)
		} else {
			dc.SetColor(grayText)
		}
		dc.DrawString(fmt.Sprintf("%d", s.Points), colX, y)

		y += rowHeight
	}

	var b bytes.Buffer
	err := png.Encode(&b, dc.Image())
	if err != nil {
		return nil
	}

	return b.Bytes()
}

func BattleStatsImage(corporation string, stats []*models.BattleStats) []byte {
	const imageWidth = 400.0
	const fontSize = 12.0
	const lineSpacing = 1.4
	const marginX = 15.0
	const marginY = 15.0
	const rowHeight = fontSize * lineSpacing

	// Ширина колонок
	const wName = 120.0
	const wLevel = 40.0
	const wPointsSum = 60.0
	const wRecords = 40.0
	const wAvgPoints = 60.0
	const wQuality = 50.0

	// Высота заголовка: отступ + заголовок корпорации + отступ + заголовки колонок
	headerHeight := marginY + (rowHeight * 1.5) + rowHeight + marginY/2
	dataRowsHeight := float64(len(stats)) * rowHeight
	footerHeight := rowHeight + marginY/2 // ← ДОБАВЛЕНО: высота для подписи

	// Общая высота: Заголовок + Строки данных + Подпись + Нижний отступ
	totalHeight := headerHeight + dataRowsHeight + footerHeight + marginY // ← ИСПРАВЛЕНО: добавлен footerHeight

	imageHeight := totalHeight

	dc := gg.NewContext(int(imageWidth), int(imageHeight))

	// Основной фон - темный
	darkBackground := color.RGBA{R: 30, G: 30, B: 30, A: 255}
	dc.SetColor(darkBackground)
	dc.Clear()

	fontPath := "docker/compendium/NotoSans.ttc"
	if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
		fmt.Println("Ошибка загрузки шрифта, используется шрифт по умолчанию:", err)
	}

	// Определения цветов
	white := color.RGBA{R: 186, G: 186, B: 186, A: 255}
	grayText := color.RGBA{R: 200, G: 200, B: 200, A: 255}
	headerBg := color.RGBA{R: 50, G: 50, B: 50, A: 255}
	rowBg := color.RGBA{R: 40, G: 40, B: 40, A: 255}

	greenHighlight := color.RGBA{R: 30, G: 224, B: 0, A: 255}
	yellowHighlight := color.RGBA{R: 255, G: 215, B: 0, A: 255}
	blueHighlight := color.RGBA{R: 84, G: 103, B: 241, A: 255}

	// Координаты начала рисования
	x := marginX
	y := marginY

	// Заголовок - название корпорации
	dc.LoadFontFace(fontPath, fontSize*1.3)
	dc.SetColor(blueHighlight)
	dc.DrawString("Статистика корпорации", x, y)
	y += rowHeight * 1.2

	dc.LoadFontFace(fontPath, fontSize*1.1)
	dc.SetColor(white)
	dc.DrawString(corporation, x, y)
	y += rowHeight * 1.3

	// Восстановление размера шрифта для таблицы
	dc.LoadFontFace(fontPath, fontSize)

	// Определение максимальных значений для выделения
	maxQuality := 0.0
	maxAvgPoints := 0.0
	maxPointsSum := 0

	for _, s := range stats {
		if s.Quality > maxQuality {
			maxQuality = s.Quality
		}
		if s.AveragePoints > maxAvgPoints {
			maxAvgPoints = s.AveragePoints
		}
		if s.PointsSum > maxPointsSum {
			maxPointsSum = s.PointsSum
		}
	}

	// Рисование фона для заголовка столбцов
	headerY := y - rowHeight/2
	dc.SetColor(headerBg)
	dc.DrawRectangle(marginX, headerY, imageWidth-2*marginX, rowHeight)
	dc.Fill()

	// Рисование заголовков столбцов
	dc.SetColor(white)
	colX := x

	dc.DrawString("Игрок", colX, y)
	colX += wName

	dc.DrawString("Ур.", colX, y)
	colX += wLevel

	dc.DrawString("Сумма", colX, y)
	colX += wPointsSum

	dc.DrawString("Игры", colX, y)
	colX += wRecords

	dc.DrawString("Среднее", colX, y)
	colX += wAvgPoints

	dc.DrawString("Кач-во", colX, y)

	y += rowHeight

	// Цикл для вывода данных с выделением
	for i, s := range stats {
		// Чередование цвета фона строк
		if i%2 == 0 {
			rowY := y - rowHeight + 4
			dc.SetColor(rowBg)
			dc.DrawRectangle(marginX, rowY, imageWidth-2*marginX, rowHeight)
			dc.Fill()
		}

		colX = x

		// Имя игрока
		dc.SetColor(grayText)
		name := s.Name
		if len(name) > 15 {
			name = name[:15] + "..."
		}
		dc.DrawString(name, colX, y)
		colX += wName

		// Уровень
		dc.SetColor(grayText)
		dc.DrawString(s.Level, colX, y)
		colX += wLevel

		// Сумма очков
		if s.PointsSum == maxPointsSum {
			dc.SetColor(greenHighlight)
		} else {
			dc.SetColor(grayText)
		}
		dc.DrawString(fmt.Sprintf("%d", s.PointsSum), colX, y)
		colX += wPointsSum

		// Количество записей
		dc.SetColor(grayText)
		dc.DrawString(fmt.Sprintf("%d", s.RecordsCount), colX, y)
		colX += wRecords

		// Средние очки
		if s.AveragePoints == maxAvgPoints {
			dc.SetColor(yellowHighlight)
		} else {
			dc.SetColor(grayText)
		}
		dc.DrawString(fmt.Sprintf("%.0f", s.AveragePoints), colX, y)
		colX += wAvgPoints

		// Качество
		if s.Quality == maxQuality {
			dc.SetColor(greenHighlight)
		} else if s.Quality >= 1.0 {
			dc.SetColor(yellowHighlight)
		} else {
			dc.SetColor(grayText)
		}
		dc.DrawString(fmt.Sprintf("%.1f", s.Quality), colX, y)

		y += rowHeight
	}

	// Добавляем подпись внизу (теперь есть место)
	footerY := y + marginY // ← ИСПРАВЛЕНО: увеличили отступ
	dc.LoadFontFace(fontPath, fontSize*0.8)
	dc.SetColor(color.RGBA{R: 150, G: 150, B: 150, A: 255})
	dc.DrawString("Качество = Средние очки / Среднее по уровню", marginX, footerY)

	var b bytes.Buffer
	err := png.Encode(&b, dc.Image())
	if err != nil {
		return nil
	}

	return b.Bytes()
}
