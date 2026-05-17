package models

import (
	"fmt"
	"time"
)

// CorpInfo представляет информацию о корпорации
type CorpInfo struct {
	ID         int64     `json:"id" db:"id"`
	CorpName   string    `json:"corp_name" db:"corp_name"`
	CorpID     string    `json:"corp_id" db:"corp_id"`
	Level      int       `json:"level" db:"level"`
	XP         int       `json:"xp" db:"xp"`
	Webhook    bool      `json:"webhook" db:"webhook"`
	LastWin    time.Time `json:"last_win" db:"last_win"`
	DateEnded  time.Time `json:"date_ended" db:"date_ended"`
	LastUpdate time.Time `json:"last_update" db:"last_update"`
}

func (c *CorpInfo) Bonus() bool {
	if c != nil {
		// Вычисляем время, до которого действует бонус
		untilTime := c.LastWin.AddDate(0, 0, 7)

		// Используем встроенный метод Before для сравнения объектов времени
		return time.Now().UTC().Before(untilTime)
	}
	return false
}

func (c *CorpInfo) GetBonusText() string {
	untilTime := c.LastWin.AddDate(0, 0, 7)
	p := percent(c.Level)
	timeStr := formatTime(untilTime.Unix())

	if c.Webhook {
		return fmt.Sprintf("%d%% %s %s\n", p, c.CorpName, timeStr)
	}
	return fmt.Sprintf("~%d%% %s %s\n", p, c.CorpName, timeStr)
}

func percent(xp int) int {
	p := 22
	for i := 2; i < xp; i++ {
		p += 2
	}
	return p
}

func formatTime(ut int64) string {
	// Определите целевую дату
	targetDate := time.Unix(ut, 0)

	// Определите оставшееся время
	remainingTime := targetDate.Sub(time.Now().UTC())

	text := ""

	// Получите дни, часы и минуты из оставшегося времени
	days := remainingTime / (24 * time.Hour)
	if days > 0 {
		text += fmt.Sprintf("%dд ", days)
	}
	remainingTime = remainingTime % (24 * time.Hour)
	hours := remainingTime / time.Hour
	if hours > 0 {
		text += fmt.Sprintf("%dч ", hours)
	}
	remainingTime = remainingTime % time.Hour
	minutes := remainingTime / time.Minute
	if days == 0 && minutes > 0 {
		text += fmt.Sprintf("%dм", minutes)
	}
	return text
}

// LevelThreshold описывает минимальный порог опыта для уровня
type LevelThreshold struct {
	Level      int
	XPRequired int
}

// Пороги опыта на основе твоего файла (отсортированы от большего к меньшему)
var corpLevelThresholds = []LevelThreshold{
	{21, 60000}, {20, 50000}, {19, 40000}, {18, 32000},
	{17, 25000}, {16, 20000}, {15, 16000}, {14, 13000},
	{13, 11000}, {12, 9000}, {11, 7000}, {10, 5000},
	{9, 3000}, {8, 2000}, {7, 1000}, {6, 500},
	{5, 250}, {4, 100}, {3, 30}, {2, 1},
	{1, 0},
}

// GetLevelByXP возвращает уровень корпорации на основе текущего XP
func (c *CorpInfo) GetLevelByXP() int {
	for _, threshold := range corpLevelThresholds {
		if c.XP >= threshold.XPRequired {
			return threshold.Level
		}
	}
	return 1 // Начальный уровень по умолчанию
}
