package server

import (
	"compendium/models"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
)

func (s *Server) getMatches(limit, page, filter string) *models.Ws {
	file, err := os.ReadFile("ws/ws.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	var cors []models.CorporationsData
	err = json.Unmarshal(file, &cors)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil
	}

	if filter != "" {
		fmt.Println(filter)
		var f models.Corporation
		err = json.Unmarshal([]byte(filter), &f)
		if err != nil {
			s.log.ErrorErr(err)
		} else if f.Id != "" {
			fmt.Println(f)
			cors = FilterCorporation(cors, f)
		}
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
	w := models.Ws{
		MaxPage: int(math.Ceil(float64(len(cors)) / float64(limitInt))),
		Matches: selectedMatches,
	}

	return &w
}

func (s *Server) getMatchesAll(limit, page, filter string) *models.Ws {
	file, err := os.ReadFile("ws/wsAll.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	var cors []models.CorporationsData
	err = json.Unmarshal(file, &cors)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil
	}

	if filter != "" {
		fmt.Println(filter)
		var f models.Corporation
		err = json.Unmarshal([]byte(filter), &f)
		if err != nil {
			s.log.ErrorErr(err)
		} else if f.Id != "" {
			fmt.Println(f)
			cors = FilterCorporation(cors, f)
		}
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
	w := models.Ws{
		MaxPage: int(math.Ceil(float64(len(cors)) / float64(limitInt))),
		Matches: selectedMatches,
	}

	return &w
}

func FilterCorporation(cors []models.CorporationsData, f models.Corporation) []models.CorporationsData {
	var newCors []models.CorporationsData
	for _, cor := range cors {
		if cor.Corporation1Id == f.Id {
			newCors = append(newCors, cor)
		}
		if cor.Corporation2Id == f.Id {
			newCors = append(newCors, cor)
		}
	}
	return newCors
}
