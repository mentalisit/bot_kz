package helpers

import (
	"fmt"
	"rs/models"
)

func (h *Helpers) NameMention(u *[]models.QueueActive, tip string) (n1, n2, n3, n4 string) {
	if len(*u) != 0 {
		for i, q := range *u {
			var name string
			if q.Data.Tip == tip {
				name = h.emReadName(q, tip, true)
			} else {
				name = h.emReadName(q, tip)
			}
			switch i {
			case 0:
				n1 = name
			case 1:
				n2 = name
			case 3:
				n3 = name
			case 4:
				n4 = name
			}
		}
	}
	return
}

func (h *Helpers) EmReadName(s *models.QueueActive, ForType string, mention ...bool) string {
	return h.emReadName(*s, ForType, mention...)
}

func (h *Helpers) emReadName(s models.QueueActive, ForType string, mention ...bool) string { // склеиваем имя и эмоджи
	name := s.Data.Name
	if s.Data.Alt != "" {
		name = s.Data.Alt
	}

	newName := s.Data.Name
	if ForType == ds {
		newName = s.Data.Mention
	} else {
		newName = s.Data.Name
	}

	if mention != nil {
		if mention[0] {
			if s.Data.Mention == "@" {
				newName = fmt.Sprintf("[%s](tg://user?id=%s)", s.Data.Name, s.Data.UserID)
			} else {
				newName = s.Data.Mention
			}
		}
	}
	if s.Data.Alt != "" {
		newName = fmt.Sprintf("%s [%s]", newName, s.Data.Alt)
	}

	multiAccount, _ := h.storage.FindMultiAccountByUserId(s.Data.UserID)
	if multiAccount == nil {
		return newName
	}
	if name == s.Data.Name {
		name = multiAccount.Nickname
	}
	module := h.storage.ModuleCompendiumGet(multiAccount.UUID, name)
	if module != nil {
		if ForType == ds {
			gen := ""
			if module.Gen != 0 {
				gen = fmt.Sprintf("<:genesis:1199068748280242237> %d ", module.Gen)
			}

			enr := ""
			if module.Enr != 0 {
				enr = fmt.Sprintf("<:enrich:1199068793633251338> %d ", module.Enr)
			}

			rse := ""
			if module.Rse != 0 {
				rse = fmt.Sprintf("<:rse:1199068829511335946> %d ", module.Rse)
			}

			newName = fmt.Sprintf("%s %s %s %s", newName, gen, enr, rse)
		} else if ForType == tg {
			if module.Rse != 0 || module.Gen != 0 || module.Enr != 0 {
				newName = fmt.Sprintf("%s (%d/%d/%d)", newName, module.Gen, module.Enr, module.Rse)
			}
		}
	}

	emoji := h.storage.EmojiReadUUID(multiAccount.UUID, ForType)
	if emoji != nil && emoji.Tip == ForType {
		newName = fmt.Sprintf("%s %s %s %s %s", newName, emoji.Em1, emoji.Em2, emoji.Em3, emoji.Em4)
	}

	return newName
}
