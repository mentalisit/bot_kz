package postgres

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetOldCompendium(guildid, userID string) (tech []byte, TzName string, TzOffset int) {
	apiKey := ""
	if guildid == "716771579278917702" {
		apiKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiI1ODI4ODIxMzc4NDIxMjI3NzMiLCJndWlsZElkIjoiNzE2NzcxNTc5Mjc4OTE3NzAyIiwiaWF0IjoxNzA2MjM3MzY0LCJleHAiOjE3Mzc3OTQ5NjQsInN1YiI6ImFwaSJ9.Wsf-2U8GDGaCNpxafRIUABIKO3zLyYKvPYWzxtbK-LE"
	}
	if guildid == "398761209022644224" {
		apiKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiI5NTIwMDAzOTIxNTYyODY5NzYiLCJndWlsZElkIjoiMzk4NzYxMjA5MDIyNjQ0MjI0IiwiaWF0IjoxNzE3MjU2NDk1LCJleHAiOjE3NDg4MTQwOTUsInN1YiI6ImFwaSJ9.EMULGfwaCLupVeBPOsSrUyBxISqXYZK4_nGHgmM96Xg"
	}
	if guildid == "632245873769971732" {
		apiKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiI3Nzc4ODE5MzYyMDY3NTc5MjkiLCJndWlsZElkIjoiNjMyMjQ1ODczNzY5OTcxNzMyIiwiaWF0IjoxNzE4MjA5MDUwLCJleHAiOjE3NDk3NjY2NTAsInN1YiI6ImFwaSJ9.Q64sbMk9-VEzTIKXWFCTabxTk_y860bQKecyFFjTuT4"
	}
	if guildid == "656495834195558402" {
		apiKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiI1ODI4ODIxMzc4NDIxMjI3NzMiLCJndWlsZElkIjoiNjU2NDk1ODM0MTk1NTU4NDAyIiwiaWF0IjoxNzIzOTExMzQwLCJleHAiOjE3NTU0Njg5NDAsInN1YiI6ImFwaSJ9.8HWmQsRCbrLYAeseYMKV_-VEk-2vuJUMDxxWwnJTWgE"
	}

	if apiKey == "" {
		return nil, "", 0
	}

	// Формирование URL-адреса
	url := fmt.Sprintf("https://bot.hs-compendium.com/compendium/api/tech?token=%s&userid=%s", apiKey, userID)

	// Выполнение GET-запроса
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return nil, "", 0
	}
	defer response.Body.Close()

	// Проверка кода ответа
	if response.StatusCode != http.StatusOK {
		fmt.Printf("Неправильный статус код: %d\n", response.StatusCode)
		return nil, "", 0
	}

	// Чтение тела ответа
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return nil, "", 0
	}

	// Декодирование JSON-данных в структуру TechnicalData
	var technicalData TechnicalData
	err = json.Unmarshal(body, &technicalData)
	if err != nil {
		fmt.Println("Ошибка при декодировании JSON:", err)
		return nil, "", 0
	}

	NewMapModule := make(map[int]TechLevel)
	for _, item := range technicalData.Array {
		NewMapModule[ModuleMap[item.Type]] = TechLevel{
			Level: int(item.Level),
			Ts:    item.Ts,
		}
	}
	bytes, _ := json.Marshal(NewMapModule)
	return bytes, technicalData.TzName, technicalData.TzOffset
}

type TechLevel struct {
	Ts    int64 `json:"ts"`
	Level int   `json:"level"`
}
type TechnicalData struct {
	TokenExpires int64           `json:"tokenExpires"`
	TzName       string          `json:"tz_name"`
	TzOffset     int             `json:"tz_offset"`
	Map          map[string]Item `json:"map"`
	Array        []Item          `json:"array"`
}

type Item struct {
	Type  string `json:"type,omitempty"`
	Level int64  `json:"level"`
	Ts    int64  `json:"ts"`
	Ws    int64  `json:"ws"`
}

var ModuleMap = map[string]int{
	"bs":              101,
	"miner":           102,
	"transp":          103,
	"battery":         202,
	"laser":           203,
	"mass":            204,
	"dual":            205,
	"barrage":         206,
	"dart":            207,
	"chainray":        208,
	"rocketlauncher":  209,
	"pulse":           210,
	"alpha":           301,
	"delta":           302,
	"passive":         303,
	"omega":           304,
	"mirror":          305,
	"blast":           306,
	"area":            307,
	"motionshield":    308,
	"cargobay":        401,
	"computer":        402,
	"rush":            404,
	"tradeburst":      405,
	"shipdrone":       406,
	"dispatch":        411,
	"relicdrone":      412,
	"remoterepair":    413,
	"cargorocket":     414,
	"miningboost":     501,
	"enrich":          503,
	"remote":          504,
	"hydroupload":     505,
	"crunch":          507,
	"genesis":         508,
	"hydrorocket":     510,
	"hydroreplicator": 511,
	"artifactboost":   512,
	"blastdrone":      513,
	"emp":             601,
	"teleport":        602,
	"rsextender":      603,
	"stealth":         608,
	"fortify":         609,
	"destiny":         614,
	"barrier":         615,
	"vengeance":       616,
	"deltarocket":     617,
	"leap":            618,
	"bond":            619,
	"omegarocket":     621,
	"suspend":         622,
	"remotebomb":      623,
	"laserturret":     624,
	"solitude":        625,
	"damageamplifier": 626,
	"rs":              701,
	"shipmentrelay":   702,
	"corplevel":       801,
	"decoydrone":      901,
	"repairdrone":     902,
	"rocketdrone":     904,
	"chainrayturret":  905,
	"deltadrone":      906,
	"dronesquad":      907,
}
