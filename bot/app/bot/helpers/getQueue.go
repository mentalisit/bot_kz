package helpers

import (
	"context"
	"fmt"
	"github.com/mentalisit/logger"
	"kz_bot/models"
	"kz_bot/storage"
)

const ds = "ds"
const tg = "tg"

type Helpers struct {
	log     *logger.Logger
	storage *storage.Storage
}

func NewHelpers(log *logger.Logger, storage *storage.Storage) *Helpers {
	return &Helpers{log: log, storage: storage}
}
func (h *Helpers) GetQueueDiscord(n map[string]string, u models.Users) map[string]string {
	if u.User1.Name != "" {
		n["name1"] = fmt.Sprintf("%s  🕒  %d  (%d)", h.emReadName(u.User1.Name, u.User1.Mention, ds), u.User1.Timedown, u.User1.Numkzn)
	}
	if u.User2.Name != "" {
		n["name2"] = fmt.Sprintf("%s  🕒  %d  (%d)", h.emReadName(u.User2.Name, u.User2.Mention, ds), u.User2.Timedown, u.User2.Numkzn)
	}
	if u.User3.Name != "" {
		n["name3"] = fmt.Sprintf("%s  🕒  %d  (%d)", h.emReadName(u.User3.Name, u.User3.Mention, ds), u.User3.Timedown, u.User3.Numkzn)
	}
	return n
}
func (h *Helpers) GetQueueTelegram(n map[string]string, u models.Users) (users string) {
	if n["text1"] != "" {
		users = n["text1"]
	}
	if u.User1.Name != "" {
		users += fmt.Sprintf("1️⃣ %s - %d%s (%d) \n",
			h.emReadName(u.User1.Name, u.User1.Mention, tg), u.User1.Timedown, n["min"], u.User1.Numkzn)
	}
	if u.User2.Name != "" {
		users += fmt.Sprintf("2️⃣ %s - %d%s (%d) \n",
			h.emReadName(u.User2.Name, u.User2.Mention, tg), u.User2.Timedown, n["min"], u.User2.Numkzn)
	}
	if u.User3.Name != "" {
		users += fmt.Sprintf("3️⃣ %s - %d%s (%d) \n",
			h.emReadName(u.User3.Name, u.User3.Mention, tg), u.User3.Timedown, n["min"], u.User3.Numkzn)
	}

	if n["text2"] != "" {
		users += n["text2"]
	}
	return users
}
func (h *Helpers) emReadName(name, nameMention, tip string) string { // склеиваем имя и эмоджи
	t := h.storage.Emoji.EmojiModuleReadUsers(context.Background(), name, tip)
	newName := name
	if tip == ds {
		newName = nameMention
	} else {
		newName = name
	}

	if len(t.Name) > 0 {
		if tip == ds && tip == t.Tip {
			newName = fmt.Sprintf("%s %s %s %s %s %s%s%s%s", nameMention, t.Module1, t.Module2, t.Module3, t.Weapon, t.Em1, t.Em2, t.Em3, t.Em4)
		} else if tip == tg && tip == t.Tip {
			newName = fmt.Sprintf("%s %s%s%s%s", name, t.Em1, t.Em2, t.Em3, t.Em4)
			if t.Weapon != "" {
				newName = fmt.Sprintf("%s [%s] %s%s%s%s", name, t.Weapon, t.Em1, t.Em2, t.Em3, t.Em4)
			}
		}
		//} else if in.Tip == ds && in.Config.Guildid == "716771579278917702" && in.Name == name {
		//	genesis, enrich, rsextender := GetTechDataUserId(in.Ds.Nameid)
		//	b.storage.Emoji.EmInsertEmpty(context.Background(), "ds", name)
		//	one := fmt.Sprintf("<:rse:1199068829511335946> %d ", rsextender)
		//	two := fmt.Sprintf("<:genesis:1199068748280242237> %d ", genesis)
		//	three := fmt.Sprintf("<:enrich:1199068793633251338> %d ", enrich)
		//	newName = fmt.Sprintf("%s ", nameMention)
		//	if rsextender != 0 {
		//		b.storage.Emoji.ModuleUpdate(context.Background(), name, "ds", "1", one)
		//		newName += one
		//	}
		//	if genesis != 0 {
		//		b.storage.Emoji.ModuleUpdate(context.Background(), name, "ds", "2", two)
		//		newName += two
		//	}
		//	if enrich != 0 {
		//		b.storage.Emoji.ModuleUpdate(context.Background(), name, "ds", "3", three)
		//		newName += three
		//	}
		//
	}
	return newName
}
func (h *Helpers) ReadNameModules(in models.InMessage) {
	t := h.storage.Emoji.EmojiModuleReadUsers(context.Background(), in.Name, "ds")
	if len(t.Name) > 0 {
		return
	} else if in.Tip == ds {
		genesis, enrich, rsextender := 0, 0, 0
		if in.Config.Guildid == "716771579278917702" {
			genesis, enrich, rsextender = GetTechDataUserId(in.Ds.Nameid)
		}
		if genesis+enrich+rsextender == 0 {
			genesis, enrich, rsextender = Get2TechDataUserId(in.Name, in.Ds.Nameid, in.Ds.Guildid)
			if genesis+enrich+rsextender == 0 {
				return
			}
		}
		h.storage.Emoji.EmInsertEmpty(context.Background(), "ds", in.Name)
		one := fmt.Sprintf("<:rse:1199068829511335946> %d ", rsextender)
		two := fmt.Sprintf("<:genesis:1199068748280242237> %d ", genesis)
		three := fmt.Sprintf("<:enrich:1199068793633251338> %d ", enrich)
		if rsextender != 0 {
			h.storage.Emoji.ModuleUpdate(context.Background(), in.Name, "ds", "1", one)
		}
		if genesis != 0 {
			h.storage.Emoji.ModuleUpdate(context.Background(), in.Name, "ds", "2", two)
		}
		if enrich != 0 {
			h.storage.Emoji.ModuleUpdate(context.Background(), in.Name, "ds", "3", three)
		}
	}
}
func (h *Helpers) UpdateCompendiumModules(in models.InMessage) string {
	genesis, enrich, rsextender := GetTechDataUserId(in.Ds.Nameid)
	if genesis+enrich+rsextender == 0 {
		genesis, enrich, rsextender = Get2TechDataUserId(in.Name, in.Ds.Nameid, in.Ds.Guildid)
		if genesis+enrich+rsextender == 0 {
			return "модули не найдены "
		}
	}

	one := fmt.Sprintf("<:rse:1199068829511335946> %d ", rsextender)
	two := fmt.Sprintf("<:genesis:1199068748280242237> %d ", genesis)
	three := fmt.Sprintf("<:enrich:1199068793633251338> %d ", enrich)
	if rsextender != 0 {
		h.storage.Emoji.ModuleUpdate(context.Background(), in.Name, "ds", "1", one)
	}
	if genesis != 0 {
		h.storage.Emoji.ModuleUpdate(context.Background(), in.Name, "ds", "2", two)
	}
	if enrich != 0 {
		h.storage.Emoji.ModuleUpdate(context.Background(), in.Name, "ds", "3", three)
	}
	return " загружено из компендиум бота " + one + two + three
}
