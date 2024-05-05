package hspublic

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"ws/models"
)

// Имитация базы данных рейтингов
var ratings = map[string]int{}

func EloLogic(match []models.Match, corps []models.Corporation) {
	for _, c := range corps {
		ratings[c.Id] = 1200
	}

	processMatches(match)

	var elo []models.Corporation
	for _, c := range corps {
		el := c
		el.Elo = ratings[c.Id]
		elo = append(elo, el)
	}

	mElo, _ := json.Marshal(elo)
	path := "ws/corps.json"
	err := os.WriteFile(path, mElo, 0644)
	if err != nil {
		return
	}
	fmt.Println("файл сохранен " + path)
}

// Функция для расчета нового рейтинга
func calculateElo(winnerRating, loserRating int) (int, int) {
	const K = 32
	expectedScoreWinner := 1.0 / (1 + math.Pow(10, float64(loserRating-winnerRating)/400))
	expectedScoreLoser := 1 - expectedScoreWinner
	newWinnerRating := winnerRating + int(float64(K)*(1-expectedScoreWinner))
	newLoserRating := loserRating + int(float64(K)*(0-expectedScoreLoser))
	return newWinnerRating, newLoserRating
}

func processMatches(matches []models.Match) {
	for _, match := range matches {
		corp1Rating := ratings[match.Corporation1Id]
		corp2Rating := ratings[match.Corporation2Id]

		if match.Corporation1Score > match.Corporation2Score {
			newCorp1Rating, newCorp2Rating := calculateElo(corp1Rating, corp2Rating)
			ratings[match.Corporation1Id] = newCorp1Rating
			ratings[match.Corporation2Id] = newCorp2Rating
		} else if match.Corporation1Score < match.Corporation2Score {
			newCorp2Rating, newCorp1Rating := calculateElo(corp2Rating, corp1Rating)
			ratings[match.Corporation2Id] = newCorp2Rating
			ratings[match.Corporation1Id] = newCorp1Rating
		}
		// Нет изменений в рейтинге в случае ничьи
	}
}
