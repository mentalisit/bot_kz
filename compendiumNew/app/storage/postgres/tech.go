package postgres

import (
	"compendium/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"io/ioutil"
	"net/http"
)

func (d *Db) TechInsert(username, userid, guildid string, tech []byte) error {
	var count int
	sel := "SELECT count(*) as count FROM hs_compendium.tech WHERE guildid = $1 AND userid = $2 AND username = $3"
	err := d.db.QueryRow(context.Background(), sel, guildid, userid, username).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		if guildid == "716771579278917702" {
			apiKey := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiI1ODI4ODIxMzc4NDIxMjI3NzMiLCJndWlsZElkIjoiNzE2NzcxNTc5Mjc4OTE3NzAyIiwiaWF0IjoxNzA2MjM3MzY0LCJleHAiOjE3Mzc3OTQ5NjQsInN1YiI6ImFwaSJ9.Wsf-2U8GDGaCNpxafRIUABIKO3zLyYKvPYWzxtbK-LE"
			getOldCompendium := GetOldCompendium(apiKey, userid)
			if getOldCompendium != nil {
				tech = getOldCompendium
			}
		}
		if guildid == "398761209022644224" {
			apiKey := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiI5NTIwMDAzOTIxNTYyODY5NzYiLCJndWlsZElkIjoiMzk4NzYxMjA5MDIyNjQ0MjI0IiwiaWF0IjoxNzE3MjU2NDk1LCJleHAiOjE3NDg4MTQwOTUsInN1YiI6ImFwaSJ9.EMULGfwaCLupVeBPOsSrUyBxISqXYZK4_nGHgmM96Xg"
			getOldCompendium := GetOldCompendium(apiKey, userid)
			if getOldCompendium != nil {
				tech = getOldCompendium
			}
		}
		insert := `INSERT INTO hs_compendium.tech(username, userid, guildid, tech) VALUES ($1,$2,$3,$4)`
		_, err = d.db.Exec(context.Background(), insert, username, userid, guildid, tech)
		if err != nil {
			return err
		}
	}
	return nil
}
func (d *Db) TechGet(username, userid, guildid string) ([]byte, error) {
	var tech []byte
	sel := "SELECT tech FROM hs_compendium.tech WHERE userid = $1 AND guildid = $2 AND username = $3"
	err := d.db.QueryRow(context.Background(), sel, userid, guildid, username).Scan(&tech)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}
	return tech, nil
}
func (d *Db) TechGetName(username, guildid string) ([]byte, string, error) {
	var tech []byte
	var userid string
	sel := "SELECT userid,tech FROM hs_compendium.tech WHERE guildid = $1 AND username = $2"
	err := d.db.QueryRow(context.Background(), sel, guildid, username).Scan(&userid, &tech)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, "", err
		}
	}
	return tech, userid, nil
}
func (d *Db) TechGetAll(cm models.CorpMember) ([]models.CorpMember, error) {
	var acm []models.CorpMember
	sel := "SELECT username,tech FROM hs_compendium.tech WHERE userid = $1 AND guildid = $2"
	q, err := d.db.Query(context.Background(), sel, cm.UserId, cm.GuildId)
	defer q.Close()
	if err != nil {
		return nil, err
	}

	for q.Next() {
		var ncm models.CorpMember
		ncm = cm
		var tech []byte
		err = q.Scan(&ncm.Name, &tech)
		if err != nil {
			return nil, err
		}

		var techl models.TechLevels
		err = json.Unmarshal(tech, &techl)
		if err != nil {
			return nil, err
		}
		if len(techl) > 0 {
			ncm.Tech = make(map[int][2]int)
			for i, level := range techl {
				ncm.Tech[i] = [2]int{level.Level, int(level.Ts)}
			}
		}
		acm = append(acm, ncm)
	}
	if err = q.Err(); err != nil { // Проверка ошибок после завершения итерации
		return nil, err
	}
	return acm, nil
}

func (d *Db) TechUpdate(username, userid, guildid string, tech []byte) error {
	upd := `update hs_compendium.tech set tech = $1 where username = $2 and userid = $3 and guildid = $4`
	updresult, err := d.db.Exec(context.Background(), upd, tech, username, userid, guildid)
	if err != nil {
		return err
	}
	if updresult.RowsAffected() == 0 {
		err = d.TechInsert(username, userid, guildid, tech)
		if err != nil {
			d.log.ErrorErr(err)
			return err
		}
	}
	return nil
}
func (d *Db) TechDelete(username, userid, guildid string) error {
	var count int
	sel := "SELECT count(*) as count FROM hs_compendium.tech WHERE guildid = $1 AND userid = $2 AND username = $3"
	err := d.db.QueryRow(context.Background(), sel, guildid, userid, username).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		del := "delete from hs_compendium.tech where username = $1 and userid = $2 and guildid = $3"
		_, err = d.db.Exec(context.Background(), del, username, userid, guildid)
		if err != nil {
			return err
		}
	}
	return nil
}
func (d *Db) Unsubscribe(ctx context.Context, name, lvlkz string, TgChannel string, tipPing int) {
	del := "delete from kzbot.subscribe where name = $1 AND lvlkz = $2 AND chatid = $3 AND tip = $4"
	_, err := d.db.Exec(ctx, del, name, lvlkz, TgChannel, tipPing)
	if err != nil {
		d.log.ErrorErr(err)
	}
}
func (d *Db) TechGetCount(userid, guildid string) (int, error) {
	var count int
	sel := "SELECT count(*) as count FROM hs_compendium.tech WHERE guildid = $1 AND userid = $2"
	err := d.db.QueryRow(context.Background(), sel, guildid, userid).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
func GetOldCompendium(apiKey, userID string) []byte {
	// Формирование URL-адреса
	url := fmt.Sprintf("https://bot.hs-compendium.com/compendium/api/tech?token=%s&userid=%s", apiKey, userID)

	// Выполнение GET-запроса
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return nil
	}
	defer response.Body.Close()

	// Проверка кода ответа
	if response.StatusCode != http.StatusOK {
		fmt.Printf("Неправильный статус код: %d\n", response.StatusCode)
		return nil
	}

	// Чтение тела ответа
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return nil
	}

	// Декодирование JSON-данных в структуру TechnicalData
	var technicalData TechnicalData
	err = json.Unmarshal(body, &technicalData)
	if err != nil {
		fmt.Println("Ошибка при декодировании JSON:", err)
		return nil
	}
	NewMapModule := make(map[int]models.TechLevel)
	for _, item := range technicalData.Array {
		NewMapModule[ModuleMap[item.Type]] = models.TechLevel{
			Level: int(item.Level),
			Ts:    item.Ts,
		}
	}
	bytes, _ := json.Marshal(NewMapModule)
	return bytes
}

type TechnicalData struct {
	TokenExpires int64           `json:"tokenExpires"`
	TzName       string          `json:"tz_name"`
	TzOffset     int64           `json:"tz_offset"`
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
