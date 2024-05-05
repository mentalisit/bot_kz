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
var elo map[string]int

func EloLogic(match []models.Match, corps []models.Corporation) {
	for _, c := range corps {
		ratings[c.Id] = 1200
	}

	processMatches(match)
	elo = make(map[string]int)

	var eloa []models.Corporation
	for _, c := range corps {
		el := c
		elocurrent := int(math.Ceil(ratings[c.Id]))
		el.Elo = elocurrent
		eloa = append(eloa, el)
		elo[c.Id] = elocurrent
	}

	mElo, _ := json.Marshal(eloa)
	path := "ws/corps.json"
	err := os.WriteFile(path, mElo, 0644)
	if err != nil {
		return
	}
	fmt.Println("файл сохранен " + path)
}

func calculateEloRating(ratingA, ratingB float64, actualScoreA, actualScoreB int, kFactor float64) (float64, float64) {
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

	expectedScoreA := 1 / (1 + math.Pow(10, (ratingB-ratingA)/400))
	expectedScoreB := 1 / (1 + math.Pow(10, (ratingA-ratingB)/400))

	newRatingA := ratingA + kFactor*(scoreA-expectedScoreA)
	newRatingB := ratingB + kFactor*(scoreB-expectedScoreB)

	return newRatingA, newRatingB
}
func processMatches(matches []models.Match) {
	for _, match := range matches {
		corp1Rating := ratings[match.Corporation1Id]
		corp2Rating := ratings[match.Corporation2Id]

		newCorp1Rating, newCorp2Rating := calculateEloRating(corp1Rating, corp2Rating, match.Corporation1Score, match.Corporation2Score, 30)
		ratings[match.Corporation1Id] = newCorp1Rating
		ratings[match.Corporation2Id] = newCorp2Rating
	}
}
