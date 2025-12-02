package serverV2

import (
	"compendium_s/server/ds"
	"time"

	"github.com/mentalisit/logger"
)

type Roles struct {
	mapaRoles map[string]map[string][]string
	timestamp map[string]int64
	log       *logger.Logger
	ds        *ds.Client
}

func NewRoles(log *logger.Logger) *Roles {
	r := &Roles{
		mapaRoles: make(map[string]map[string][]string),
		timestamp: make(map[string]int64),
		log:       log,
		ds:        ds.NewClient(log),
	}
	return r
}
func (r *Roles) LoadGuild(guildId string) {
	if r.timestamp[guildId] == 0 || r.timestamp[guildId]+600 < time.Now().Unix() {

		membersRoles, err := r.ds.GetMembersRoles(guildId)
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
