package helpers

import (
	"bytes"
	"fmt"
	"github.com/fogleman/gg"
	"image/color"
	"image/png"
	"os"
	"rs/models"
	"sort"
	"strconv"
	"strings"
)

func (h *Helpers) CreateScoreboard(filename string, corpName string, eventId int) string {
	all, err := h.storage.Battles.BattlesGetAll(corpName, eventId)
	if err != nil {
		h.log.ErrorErr(err)
		return ""
	}
	if len(all) == 0 {
		return h.CreateScoreboardTop(filename, corpName)
	}
	if corpName == "rusb" {
		allbest, _ := h.storage.Battles.BattlesGetAll("best", eventId)
		aa := make(map[string]models.PlayerStats)
		for _, stats := range all {
			aa[stats.Player] = stats
		}
		for _, stats := range allbest {
			if existing, ok := aa[stats.Player]; ok {
				// Если уже есть — складываем нужные поля
				existing.Points += stats.Points
				existing.Runs += stats.Runs
				aa[stats.Player] = existing // обновляем обратно
			} else {
				// Иначе просто записываем
				aa[stats.Player] = stats
			}
		}
		all = []models.PlayerStats{}
		for _, stats := range aa {
			all = append(all, stats)
		}
	}
	sort.Slice(all, func(i, j int) bool { return all[i].Points > all[j].Points })
	var data []models.EntryScoreboard
	for _, stats := range all {
		data = append(data, models.EntryScoreboard{
			DisplayName: stats.Player,
			RsLevel:     stats.Level,
			StarsCount:  stats.Runs,
			Score:       stats.Points,
		})
	}
	if len(data) > 60 {
		data = data[:59]
	}
	folder := "docker/scoreboard/" + filename

	err = CreateScoreboardImage(data, folder)
	if err != nil {
		h.log.ErrorErr(err)
		return ""
	}

	return folder
}

// Функция для создания изображения с табло
func CreateScoreboardImage(data []models.EntryScoreboard, filePath string) (err error) {
	// Цвета и настройки
	var (
		colorBg       = color.RGBA{R: 18, G: 18, B: 18, A: 255}
		colorMainText = color.RGBA{R: 255, G: 183, B: 77, A: 255}
		colorName     = color.RGBA{R: 187, G: 134, B: 252, A: 255}
		colorScore    = color.RGBA{R: 52, G: 168, B: 83, A: 255}
		colorRs       = color.RGBA{R: 255, G: 75, B: 75, A: 255}

		normalFontSize = 18.0
		rsFontSize     = 20.0

		width  = 900
		height = 930

		offsets = map[string]float64{
			"rank":    0,
			"name":    40,
			"rs":      220,
			"rsCount": 280,
			"score":   340,
		}

		headerY    = 70.0
		lineHeight = 28.0

		fontPath = "docker/compendium/NotoSans.ttc" // Путь к шрифту
	)

	formatNumber := func(n int) string {
		s := fmt.Sprintf("%d", n)
		var sb strings.Builder
		length := len(s)
		for i, c := range s {
			if i > 0 && (length-i)%3 == 0 {
				sb.WriteRune(' ')
			}
			sb.WriteRune(c)
		}
		return sb.String()
	}

	if len(data) <= 30 {
		height = 110 + len(data)*int(lineHeight)
		width /= 2 // Если игроков ≤ 30, ширина в два раза меньше
	}

	leftColX := 20.0
	rightColX := leftColX + float64(width)/2 // Уменьшено расстояние между колонками

	dc := gg.NewContext(width, height)
	dc.SetColor(colorBg)
	dc.Clear()

	if err = dc.LoadFontFace(fontPath, normalFontSize); err != nil {
		return err
	}

	dc.SetColor(colorMainText)
	dc.DrawStringAnchored("Топ участников", 50, 35, 0, 0)

	// Отрисовка заголовков столбцов
	dc.SetColor(colorMainText)
	dc.DrawString("№", leftColX+offsets["rank"], headerY)
	dc.DrawString("Имя", leftColX+offsets["name"], headerY)
	dc.DrawString("Ур", leftColX+offsets["rs"], headerY)
	dc.DrawString("Кол", leftColX+offsets["rsCount"], headerY)
	dc.DrawString("Счёт", leftColX+offsets["score"], headerY)
	if len(data) > 30 {
		dc.DrawString("№", rightColX+offsets["rank"], headerY)
		dc.DrawString("Имя", rightColX+offsets["name"], headerY)
		dc.DrawString("Ур", rightColX+offsets["rs"], headerY)
		dc.DrawString("Кол", rightColX+offsets["rsCount"], headerY)
		dc.DrawString("Счёт", rightColX+offsets["score"], headerY)
	}

	var y float64
	colX := leftColX
	var total int
	for i, entry := range data {
		total += entry.Score
		rowIndex := i
		if len(data) > 30 && i >= 30 {
			colX = rightColX
			rowIndex -= 30
			dc.SetColor(colorName)
			dc.DrawString("|", colX-20, y)
		}

		y = headerY + float64(rowIndex+1)*lineHeight

		rank := strconv.Itoa(1 + i)
		name := entry.DisplayName
		if name == "" {
			name = "Unknown"
		}

		rsString := strconv.Itoa(entry.RsLevel)
		rsharp := strconv.Itoa(entry.StarsCount)
		score := formatNumber(entry.Score)

		dc.SetColor(colorMainText)
		dc.DrawString(rank, colX+offsets["rank"], y)

		dc.SetColor(colorName)
		dc.DrawString(name, colX+offsets["name"], y)

		dc.SetColor(colorRs)
		_ = dc.LoadFontFace(fontPath, rsFontSize)
		dc.DrawString(rsString, colX+offsets["rs"], y)

		dc.SetColor(colorMainText)
		_ = dc.LoadFontFace(fontPath, normalFontSize)
		dc.DrawString(rsharp, colX+offsets["rsCount"], y)

		dc.SetColor(colorScore)
		dc.DrawString(score, colX+offsets["score"], y)
	}

	dc.SetColor(colorName)
	dc.DrawString("Всего:  ", colX+240, y+lineHeight)
	dc.SetColor(colorScore)
	dc.DrawString(formatNumber(total), colX+300, y+lineHeight)

	var buf bytes.Buffer
	err = png.Encode(&buf, dc.Image())
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, buf.Bytes(), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (h *Helpers) CreateScoreboardTop(filename string, corpName string) string {
	all, err := h.storage.Battles.BattlesTopGetAll(corpName)
	if err != nil {
		h.log.ErrorErr(err)
		return ""
	}
	if len(all) == 0 {
		return ""
	}
	sort.Slice(all, func(i, j int) bool { return all[i].Count > all[j].Count })
	var data []models.EntryScoreboard
	for _, stats := range all {
		data = append(data, models.EntryScoreboard{
			DisplayName: stats.Name,
			RsLevel:     stats.Level,
			StarsCount:  stats.Count,
			Score:       0,
		})
	}
	if len(data) > 60 {
		data = data[:59]
	}
	folder := "docker/scoreboard/" + filename

	err = CreateScoreboardImage(data, folder)
	if err != nil {
		h.log.ErrorErr(err)
		return ""
	}

	return folder
}
