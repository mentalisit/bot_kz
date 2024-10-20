package helpers

import (
	"fmt"
	"kz_bot/models"
	"kz_bot/pkg/utils"
	"time"
)

type SaveDM struct {
	usersId   []string
	config    models.CorporationConfig
	timestamp int64
}

func (h *Helpers) IfMessageDM(in models.InMessage) (dm bool, conf models.CorporationConfig) {
	if in.Tip == "dsDM" || in.Tip == "tgDM" {
		h.log.Info(fmt.Sprintf("%s @%s: %s", in.Tip, in.Username, in.Mtext))
		dm = true
	}

	conf = h.checkUserid(in)

	return
}
func (h *Helpers) SaveUsersIdQueue(users []string, config models.CorporationConfig) {
	ch := utils.WaitForMessage("SendChannelDelSecond")
	defer close(ch)
	s := SaveDM{
		usersId:   users,
		config:    config,
		timestamp: time.Now().UTC().Unix(),
	}
	h.saveArray = append(h.saveArray, s)
}
func (h *Helpers) checkUserid(in models.InMessage) (conf models.CorporationConfig) {
	var newSaveDM []SaveDM
	if len(h.saveArray) > 0 {
		for _, dm := range h.saveArray {
			if dm.timestamp+600 < time.Now().UTC().Unix() {
				newSaveDM = append(newSaveDM, dm)
				for _, s := range dm.usersId {
					if in.UserId == s {
						conf = dm.config
					}
				}
			}
		}
	}
	h.saveArray = newSaveDM
	return
}
