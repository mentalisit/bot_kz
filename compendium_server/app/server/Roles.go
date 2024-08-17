package server

import (
	"bytes"
	"compendium_s/models"
	"encoding/json"
	"fmt"
	"github.com/mentalisit/logger"
	"net/http"
	"time"
)

type Roles struct {
	mapaRoles map[string]map[string][]string
	timestamp map[string]int64
	log       *logger.Logger
}

func NewRoles(log *logger.Logger) *Roles {
	r := &Roles{
		mapaRoles: make(map[string]map[string][]string),
		timestamp: make(map[string]int64),
		log:       log,
	}
	return r
}
func (r *Roles) LoadGuild(guildId string) {
	if r.timestamp[guildId] == 0 || r.timestamp[guildId]+600 < time.Now().Unix() {

		membersRoles, err := GetMembersRoles(guildId)
		if err != nil {
			r.log.ErrorErr(err)
			return
		}

		m := make(map[string][]string)
		for _, role := range membersRoles {
			m[role.Userid] = role.RolesId
		}
		r.timestamp[guildId] = time.Now().Unix()
		r.mapaRoles[guildId] = m

	}

}
func (r *Roles) CheckRoleDs(guildId, memderId, roleid string) bool {
	if roleid == "" {
		return true
	}
	if len(r.mapaRoles[guildId]) > 0 && len(r.mapaRoles[guildId][memderId]) > 0 {
		for _, role := range r.mapaRoles[guildId][memderId] {
			if role == roleid {
				return true
			}
		}
	}
	return false
}

func GetMembersRoles(guildId string) ([]models.DsMembersRoles, error) {
	data, err := json.Marshal(guildId)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post("http://"+apiname+"/discord/GetMembersRoles", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
		//resp, err = http.Post("http://192.168.100.155:802/data", "application/json", bytes.NewBuffer(data))
		return nil, err
	}

	var guildMembersRoles []models.DsMembersRoles

	err = json.NewDecoder(resp.Body).Decode(&guildMembersRoles)
	if err != nil {
		return nil, err
	}

	return guildMembersRoles, nil
}
