package server

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"
)

type CorporationsData struct {
	Corp1Name  string    `json:"Corp1Name"`
	Corp2Name  string    `json:"Corp2Name"`
	Corp1Score int       `json:"Corp1Score"`
	Corp2Score int       `json:"Corp2Score"`
	DateEnded  time.Time `json:"DateEnded"`
}
type ws struct {
	MaxPage int                `json:"MaxPage"`
	Matches []CorporationsData `json:"matches"`
}

func (s *Server) getMatches(limit, page string) *ws {
	file, err := os.ReadFile("ws/ws.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	var cors []CorporationsData
	err = json.Unmarshal(file, &cors)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil
	}

	// Обработка параметра limit
	limitInt := len(cors) // По умолчанию возвращаем все записи
	if limit != "" {
		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			fmt.Println("Error converting limit to integer:", err)
			return nil
		}
	}

	// Обработка параметра page
	pageInt := 1 // По умолчанию используется первая страница
	if page != "" {
		pageInt, err = strconv.Atoi(page)
		if err != nil {
			fmt.Println("Error converting page to integer:", err)
			return nil
		}
	}

	// Расчет начального индекса
	startIndex := (pageInt - 1) * limitInt
	if startIndex > len(cors) || startIndex < 0 {
		fmt.Println("Invalid page number")
		return nil
	}

	// Расчет конечного индекса
	endIndex := startIndex + limitInt
	if endIndex > len(cors) {
		endIndex = len(cors)
	}

	// Создание подмножества корпораций
	selectedMatches := cors[startIndex:endIndex]

	// Создание и возвращение структуры ws с данными о пагинации
	w := ws{
		MaxPage: int(math.Ceil(float64(len(cors)) / float64(limitInt))),
		Matches: selectedMatches,
	}

	return &w
}
