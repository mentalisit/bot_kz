package server

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"ws/models"
)

//func (s *Server) getMatches(limit, page, filter string) *models.Ws {
//	file, err := os.ReadFile("ws/ws.json")
//	if err != nil {
//		fmt.Println("Error reading file:", err)
//		return nil
//	}
//
//	var cors []models.CorporationsData
//	err = json.Unmarshal(file, &cors)
//	if err != nil {
//		fmt.Println("Error unmarshalling JSON:", err)
//		return nil
//	}
//
//	if filter != "" {
//		fmt.Println(filter)
//		var f models.Corporation
//		err = json.Unmarshal([]byte(filter), &f)
//		if err != nil {
//			s.log.ErrorErr(err)
//		} else if f.Id != "" {
//			fmt.Println(f)
//			cors = FilterCorporation(cors, f)
//		}
//	}
//
//	// Обработка параметра limit
//	limitInt := len(cors) // По умолчанию возвращаем все записи
//	if limit != "" {
//		limitInt, err = strconv.Atoi(limit)
//		if err != nil {
//			fmt.Println("Error converting limit to integer:", err)
//			return nil
//		}
//	}
//
//	// Обработка параметра page
//	pageInt := 1 // По умолчанию используется первая страница
//	if page != "" {
//		pageInt, err = strconv.Atoi(page)
//		if err != nil {
//			fmt.Println("Error converting page to integer:", err)
//			return nil
//		}
//	}
//
//	// Расчет начального индекса
//	startIndex := (pageInt - 1) * limitInt
//	if startIndex > len(cors) || startIndex < 0 {
//		fmt.Println("Invalid page number")
//		return nil
//	}
//
//	// Расчет конечного индекса
//	endIndex := startIndex + limitInt
//	if endIndex > len(cors) {
//		endIndex = len(cors)
//	}
//
//	// Создание подмножества корпораций
//	selectedMatches := cors[startIndex:endIndex]
//
//	// Создание и возвращение структуры ws с данными о пагинации
//	w := models.Ws{
//		MaxPage: int(math.Ceil(float64(len(cors)) / float64(limitInt))),
//		Matches: selectedMatches,
//	}
//
//	return &w
//}

func (s *Srv) getMatchesAll(limit, page, filter string) *models.Ws {
	file, err := os.ReadFile("ws/wsAll.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	var cors []models.Match
	err = json.Unmarshal(file, &cors)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil
	}

	if filter != "" {
		var f models.FilterCorps
		err = json.Unmarshal([]byte(filter), &f)
		if err != nil {
			fmt.Printf("Filter bad: %+v\n", filter)
			s.log.ErrorErr(err)
		} else if f.Corp[0].Id != "" {
			cors = FilterCorporation(cors, f.Corp)
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

func FilterCorporation(cors []models.Match, f []models.Corporation) []models.Match {
	fmap := make(map[string]models.Corporation)
	for _, corporation := range f {
		fmap[corporation.Id] = corporation
	}

	var newCors []models.Match
	var filterName []string

	for _, cor := range cors {
		if fmap[cor.Corporation1Id].Id == cor.Corporation1Id {
			newCors = append(newCors, cor)
			filterName = append(filterName, cor.Corporation1Name)
		}
		if fmap[cor.Corporation2Id].Id == cor.Corporation2Id {
			newCors = append(newCors, cor)
			filterName = append(filterName, cor.Corporation2Name)
		}
	}

	fmt.Printf("FilterCorporation: %+v\n", RemoveDuplicates(filterName))

	return newCors
}

func RemoveDuplicates[T comparable](a []T) []T {
	result := make([]T, 0, len(a))
	temp := map[T]struct{}{}
	for _, item := range a {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
