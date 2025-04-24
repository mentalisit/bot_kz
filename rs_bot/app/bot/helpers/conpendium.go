package helpers

//type TechnicalData struct {
//	Map   map[string]Item `json:"map"`
//	Array []Item          `json:"array"`
//}
//
//type Item struct {
//	Type  string `json:"type"`
//	Level int    `json:"level"`
//	Ws    int    `json:"ws"`
//}

//func GetTechDataUserId(userID, guildid string) (genesis, enrich, rsextender int) {
//	apiKey := ""
//	if guildid == "716771579278917702" {
//		apiKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiI1ODI4ODIxMzc4NDIxMjI3NzMiLCJndWlsZElkIjoiNzE2NzcxNTc5Mjc4OTE3NzAyIiwiaWF0IjoxNzA2MjM3MzY0LCJleHAiOjE3Mzc3OTQ5NjQsInN1YiI6ImFwaSJ9.Wsf-2U8GDGaCNpxafRIUABIKO3zLyYKvPYWzxtbK-LE"
//	}
//	if guildid == "398761209022644224" {
//		apiKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiI5NTIwMDAzOTIxNTYyODY5NzYiLCJndWlsZElkIjoiMzk4NzYxMjA5MDIyNjQ0MjI0IiwiaWF0IjoxNzE3MjU2NDk1LCJleHAiOjE3NDg4MTQwOTUsInN1YiI6ImFwaSJ9.EMULGfwaCLupVeBPOsSrUyBxISqXYZK4_nGHgmM96Xg"
//	}
//	if guildid == "632245873769971732" {
//		apiKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiI3Nzc4ODE5MzYyMDY3NTc5MjkiLCJndWlsZElkIjoiNjMyMjQ1ODczNzY5OTcxNzMyIiwiaWF0IjoxNzE4MjA5MDUwLCJleHAiOjE3NDk3NjY2NTAsInN1YiI6ImFwaSJ9.Q64sbMk9-VEzTIKXWFCTabxTk_y860bQKecyFFjTuT4"
//	}
//	if guildid == "656495834195558402" {
//		apiKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiI1ODI4ODIxMzc4NDIxMjI3NzMiLCJndWlsZElkIjoiNjU2NDk1ODM0MTk1NTU4NDAyIiwiaWF0IjoxNzIzOTExMzQwLCJleHAiOjE3NTU0Njg5NDAsInN1YiI6ImFwaSJ9.8HWmQsRCbrLYAeseYMKV_-VEk-2vuJUMDxxWwnJTWgE"
//	}
//
//	if apiKey == "" {
//		return 0, 0, 0
//	}
//
//	// Формирование URL-адреса
//	url := fmt.Sprintf("https://bot.hs-compendium.com/compendium/api/tech?token=%s&userid=%s", apiKey, userID)
//
//	// Выполнение GET-запроса
//	response, err := http.Get(url)
//	if err != nil {
//		fmt.Println("Ошибка при выполнении запроса:", err)
//		return
//	}
//	defer response.Body.Close()
//
//	// Проверка кода ответа
//	if response.StatusCode != http.StatusOK {
//		fmt.Printf("Неправильный статус код: %d\n", response.StatusCode)
//		return
//	}
//
//	// Чтение тела ответа
//	body, err := ioutil.ReadAll(response.Body)
//	if err != nil {
//		fmt.Println("Ошибка при чтении ответа:", err)
//		return
//	}
//
//	// Декодирование JSON-данных в структуру TechnicalData
//	var technicalData TechnicalData
//	err = json.Unmarshal(body, &technicalData)
//	if err != nil {
//		fmt.Println("Ошибка при декодировании JSON:", err)
//		return
//	}
//	genesis = technicalData.Map["genesis"].Level
//	enrich = technicalData.Map["enrich"].Level
//	rsextender = technicalData.Map["rsextender"].Level
//	return
//}
