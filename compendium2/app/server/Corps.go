package server

import (
	"compendium/models"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
)

func (s *Server) getCorps(limit, page string) *models.Corp {
	file, err := os.ReadFile("ws/corps.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	var cors []models.Corporation
	err = json.Unmarshal(file, &cors)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil
	}
	sort.Slice(cors, func(i, j int) bool {
		return cors[i].Name < cors[j].Name
	})

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
	w := models.Corp{
		MaxPage: int(math.Ceil(float64(len(cors)) / float64(limitInt))),
		Matches: selectedMatches,
	}

	return &w
}
