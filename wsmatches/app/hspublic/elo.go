package hspublic

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"ws/models"
)

// Имитация базы данных рейтингов
var ratings = map[string]float64{}

func EloLogic(match []models.Match, corps []models.Corporation) {
	for _, c := range corps {
		ratings[c.Id] = 1200
	}

	processMatches(match)

	var elo []models.Corporation
	for _, c := range corps {
		el := c
		el.Elo = int(ratings[c.Id])
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

func calculateEloRating(ratingA, ratingB, actualScoreA, actualScoreB float64, kFactor float64) (float64, float64) {
	var scoreA, scoreB float64
	if actualScoreA > actualScoreB {
		// Победа игрока A
		scoreA = 1
		scoreB = 0
	} else if actualScoreA < actualScoreB {
		// Победа игрока B
		scoreA = 0
		scoreB = 1
	} else {
		// Ничья
		scoreA = 0.5
		scoreB = 0.5
	}

	expectedScoreA := 1 / (1 + math.Pow(10, float64(ratingB-ratingA)/400))
	expectedScoreB := 1 / (1 + math.Pow(10, float64(ratingA-ratingB)/400))

	newRatingA := ratingA + kFactor*(scoreA-expectedScoreA)
	newRatingB := ratingB + kFactor*(scoreB-expectedScoreB)

	return newRatingA, newRatingB
}
func processMatches(matches []models.Match) {
	for _, match := range matches {
		corp1Rating := ratings[match.Corporation1Id]
		corp2Rating := ratings[match.Corporation2Id]

		newCorp1Rating, newCorp2Rating := calculateEloRating(corp1Rating, corp2Rating, float64(match.Corporation1Score), float64(match.Corporation2Score), 30)
		ratings[match.Corporation1Id] = newCorp1Rating
		ratings[match.Corporation2Id] = newCorp2Rating
	}
}
