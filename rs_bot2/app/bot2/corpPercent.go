package bot2

import (
	"errors"
	"fmt"
	"rs/models"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

func (b *Bot) SendOtherCorporationsPercent(Config models.CorporationConfigV2) {
	fmt.Println("SendOtherCorporationsPercent")
	bonuses := make(map[string]string)
	needSend := make(map[string]*models.Info)

	for ch, info := range Config.Channels {
		if info.Game != nil && info.Game.GameCorporationId != "" {
			currentCorp, err := b.storage.ReadCorpInfoByCorpID(info.Game.GameCorporationId)
			if err != nil {
				b.log.ErrorErr(err)
			}
			if currentCorp != nil && currentCorp.Bonus() {
				bonuses[info.Game.GameCorporationId] = currentCorp.GetBonusText()
				continue
			}
			needSend[ch] = info
		}
	}

	if len(needSend) != 0 {
		if Config.Bonuses != nil && len(Config.Bonuses) > 0 {
			for _, bonus := range Config.Bonuses {
				ci, err := b.storage.ReadCorpInfoByCorpID(bonus.GameCorporationId)
				if err != nil {
					b.log.ErrorErr(err)
				}
				if ci != nil {
					bonuses[ci.CorpID] = ci.GetBonusText()
				}
			}
		}
		var preparingText []string
		for _, persentText := range bonuses {
			preparingText = append(preparingText, persentText)
		}
		sortText := sortByFirstTwoDigits(preparingText)

		text := ""
		for _, s := range sortText {
			text += s
		}
		for ch, info := range needSend {
			b.sendTypeMessenger(info.TypeMessenger, ch, text, 180)
		}
	}
}

func (b *Bot) GetTextPercent(g *models.GameSettings, dark bool) string {
	if g == nil || g.GameCorporationId == "" {
		return ""
	}
	currentCorp, err := b.storage.ReadCorpInfoByCorpID(g.GameCorporationId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			twoYearsAgo := time.Now().AddDate(-2, 0, 0) // Отнимаем 2 года

			_, err = b.storage.CreateCorpInfo(models.CorpInfo{
				CorpName:   g.GameCorporation,
				CorpID:     g.GameCorporationId,
				Level:      g.GameLevel,
				XP:         g.GameXP,
				LastWin:    twoYearsAgo,
				DateEnded:  twoYearsAgo,
				LastUpdate: time.Now().UTC(),
			})
			if err != nil {
				b.log.ErrorErr(err)
			}
			return ""
		}
		b.log.ErrorErr(err)
		return ""
	} else if currentCorp == nil {
		return ""
	}
	return currentCorp.GetBonusText()
}

func sortByFirstTwoDigits(input []string) []string {
	// Вспомогательная функция для извлечения числа из начала строки
	extractNumber := func(s string) int {
		if s == "" {
			return 0
		}

		// Убираем тильду и пробелы в начале, если они есть
		s = strings.TrimLeft(s, "~ ")

		// Ищем, где заканчиваются цифры
		var numStr string
		for _, r := range s {
			if r >= '0' && r <= '9' {
				numStr += string(r)
			} else {
				// Как только встретили не цифру (например, '%'), прерываем цикл
				break
			}
		}

		// Преобразуем накопленную строку в число
		num, _ := strconv.Atoi(numStr)
		return num

	}

	// Используем sort.Slice напрямую, кастомный тип тут не обязателен
	sort.Slice(input, func(i, j int) bool {
		return extractNumber(input[i]) > extractNumber(input[j])
	})

	return input
}

func (b *Bot) GetTextQueueComplite(c *models.CorporationConfigV2, dark bool) string {
	if c == nil {
		return ""
	}
	var text string
	for _, info := range c.Channels {
		if info.Game != nil && info.Game.GameCorporationId != "" {
			ci, _ := b.storage.ReadCorpInfoByCorpID(info.Game.GameCorporationId)
			if ci.Bonus() {
				text = fmt.Sprintf("%s\n", ci.GetBonusText())
			}
		}
	}
	return text
}
