package webServer

import (
	ds "rs/clients/DsApi"
	"time"

	"github.com/mentalisit/logger"
)

type RolesHelper struct {
	mapaRoles map[string]map[string][]string
	timestamp map[string]int64
	log       *logger.Logger
	ds        *ds.Client
}

func NewRolesHelper(log *logger.Logger, dsClient *ds.Client) *RolesHelper {
	r := &RolesHelper{
		mapaRoles: make(map[string]map[string][]string),
		timestamp: make(map[string]int64),
		log:       log,
		ds:        dsClient,
	}
	return r
}

func (r *RolesHelper) LoadGuild(guildId string) {
	if r.timestamp[guildId] == 0 || r.timestamp[guildId]+600 < time.Now().Unix() {

		membersRoles, err := r.ds.GetMembersRoles(guildId)
		if err != nil {
			r.log.ErrorErr(err)
			return
		}
		if len(membersRoles) == 0 {
			r.log.InfoStruct("len(membersRoles)==0  ", guildId)
		}

		m := make(map[string][]string)
		for _, role := range membersRoles {
			m[role.Userid] = role.RolesId
		}
		r.timestamp[guildId] = time.Now().Unix()
		r.mapaRoles[guildId] = m

	}

}

func (r *RolesHelper) CheckRoleDs(guildId, memderId, roleid string) bool {
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
