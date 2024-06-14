package helpers

import (
	"context"
	"fmt"
	"github.com/mentalisit/logger"
	"kz_bot/models"
	"kz_bot/storage"
	"regexp"
	"strconv"
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
		n["name1"] = fmt.Sprintf("%s  🕒  %d  (%d)", h.emReadName(u.User1, ds), u.User1.Timedown, u.User1.Numkzn)
	}
	if u.User2.Name != "" {
		n["name2"] = fmt.Sprintf("%s  🕒  %d  (%d)", h.emReadName(u.User2, ds), u.User2.Timedown, u.User2.Numkzn)
	}
	if u.User3.Name != "" {
		n["name3"] = fmt.Sprintf("%s  🕒  %d  (%d)", h.emReadName(u.User3, ds), u.User3.Timedown, u.User3.Numkzn)
	}
	return n
}
func (h *Helpers) GetQueueTelegram(n map[string]string, u models.Users) (users string) {
	if n["text1"] != "" {
		users = n["text1"]
	}
	if u.User1.Name != "" {
		users += fmt.Sprintf("1️⃣ %s - %d%s (%d) \n",
			h.emReadName(u.User1, tg), u.User1.Timedown, n["min"], u.User1.Numkzn)
	}
	if u.User2.Name != "" {
		users += fmt.Sprintf("2️⃣ %s - %d%s (%d) \n",
			h.emReadName(u.User2, tg), u.User2.Timedown, n["min"], u.User2.Numkzn)
	}
	if u.User3.Name != "" {
		users += fmt.Sprintf("3️⃣ %s - %d%s (%d) \n",
			h.emReadName(u.User3, tg), u.User3.Timedown, n["min"], u.User3.Numkzn)
	}

	if n["text2"] != "" {
		users += n["text2"]
	}
	return users
}
func (h *Helpers) emReadName(s models.Sborkz, tip string) string { // склеиваем имя и эмоджи
	name := s.Name
	if s.Wamesid != "" {
		name = s.Wamesid
	}
	t := h.storage.Emoji.EmojiModuleReadUsers(context.Background(), name, tip)
	newName := s.Name
	if tip == ds {
		newName = s.Mention
	} else {
		newName = s.Name
	}

	if len(t.Name) > 0 {
		if tip == ds && tip == t.Tip {
			newName = fmt.Sprintf("%s %s %s %s %s %s%s%s%s", s.Mention, t.Module1, t.Module2, t.Module3, t.Weapon, t.Em1, t.Em2, t.Em3, t.Em4)
			if s.Wamesid != "" {
				newName = fmt.Sprintf("%s [%s]  %s %s %s %s%s%s%s", s.Mention, s.Wamesid, t.Module1, t.Module2, t.Module3, t.Em1, t.Em2, t.Em3, t.Em4)
			}
		} else if tip == tg && tip == t.Tip {
			newName = fmt.Sprintf("%s %s%s%s%s", s.Name, t.Em1, t.Em2, t.Em3, t.Em4)
			if t.Weapon != "" {
				newName = fmt.Sprintf("%s [%s] %s%s%s%s", s.Name, t.Weapon, t.Em1, t.Em2, t.Em3, t.Em4)
			}
			if s.Wamesid != "" {
				newName = fmt.Sprintf("%s [%s] %s%s%s%s", s.Name, s.Wamesid, t.Em1, t.Em2, t.Em3, t.Em4)
			}
		}
	}
	return newName
}
func (h *Helpers) ReadNameModules(in models.InMessage, name string) {
	if in.Tip == ds {
		if name == "" {
			name = in.Name
		}
		t := h.storage.Emoji.EmojiModuleReadUsers(context.Background(), name, ds)

		genesis1, enrich1, rsextender1 := 0, 0, 0
		if name == in.Name {
			genesis1, enrich1, rsextender1 = GetTechDataUserId(in.Ds.Nameid, in.Config.Guildid)
		}
		genesis2, enrich2, rsextender2 := Get2TechDataUserId(name, in.Ds.Nameid, in.Ds.Guildid)

		genesis := max(genesis1, genesis2, extractNumbers(t.Module2))
		enrich := max(enrich1, enrich2, extractNumbers(t.Module3))
		rsextender := max(rsextender1, rsextender2, extractNumbers(t.Module1))
		fmt.Printf("genesis %d enrich %d rsextender %d for:%s\n", genesis, enrich, rsextender, name)

		if t.Name == "" {
			h.storage.Emoji.EmInsertEmpty(context.Background(), "ds", name)
		}
		one := fmt.Sprintf("<:rse:1199068829511335946> %d ", rsextender)
		two := fmt.Sprintf("<:genesis:1199068748280242237> %d ", genesis)
		three := fmt.Sprintf("<:enrich:1199068793633251338> %d ", enrich)
		if rsextender != 0 && t.Module1 != one {
			h.storage.Emoji.ModuleUpdate(context.Background(), name, "ds", "1", one)
		}
		if genesis != 0 && t.Module2 != two {
			h.storage.Emoji.ModuleUpdate(context.Background(), name, "ds", "2", two)
		}
		if enrich != 0 && t.Module3 != three {
			h.storage.Emoji.ModuleUpdate(context.Background(), name, "ds", "3", three)
		}
	}
}
func (h *Helpers) UpdateCompendiumModules(in models.InMessage) string {
	genesis, enrich, rsextender := GetTechDataUserId(in.Ds.Nameid, in.Config.Guildid)
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

func (h *Helpers) NameMention(u models.Users, tip string) (n1, n2, n3, n4 string) {
	if u.User1.Tip == tip {
		n1 = h.emReadMention(u.User1, tip)
	} else {
		n1 = u.User1.Name
	}
	if u.User2.Tip == tip {
		n2 = h.emReadMention(u.User2, tip)
	} else {
		n2 = u.User2.Name
	}
	if u.User3.Tip == tip {
		n3 = h.emReadMention(u.User3, tip)
	} else {
		n3 = u.User3.Name
	}
	if u.User4.Tip == tip {
		n4 = h.emReadMention(u.User4, tip)
	} else {
		n4 = u.User4.Name
	}
	return
}
func (h *Helpers) emReadMention(u models.Sborkz, tip string) string { // склеиваем имя и эмоджи
	t := models.EmodjiUser{}
	if u.Wamesid != "" {
		t = h.storage.Emoji.EmojiModuleReadUsers(context.Background(), u.Wamesid, tip)
	} else {
		t = h.storage.Emoji.EmojiModuleReadUsers(context.Background(), u.Name, tip)
	}

	newName := u.Mention

	if len(t.Name) > 0 {
		if tip == ds && tip == t.Tip {
			newName = fmt.Sprintf("%s %s %s %s %s %s%s%s%s", u.Mention, t.Module1, t.Module2, t.Module3, t.Weapon, t.Em1, t.Em2, t.Em3, t.Em4)
			if u.Wamesid != "" {
				newName = fmt.Sprintf("%s [%s] %s %s %s %s%s%s%s", u.Mention, u.Wamesid, t.Module1, t.Module2, t.Module3, t.Em1, t.Em2, t.Em3, t.Em4)
			}
		} else if tip == tg && tip == t.Tip {
			newName = fmt.Sprintf("%s %s%s%s%s", u.Mention, t.Em1, t.Em2, t.Em3, t.Em4)
			if t.Weapon != "" {
				newName = fmt.Sprintf("%s [%s] %s%s%s%s", u.Mention, t.Weapon, t.Em1, t.Em2, t.Em3, t.Em4)
			}
			if u.Wamesid != "" {
				newName = fmt.Sprintf("%s [%s] %s%s%s%s", u.Mention, u.Wamesid, t.Em1, t.Em2, t.Em3, t.Em4)
			}
		}
	}
	return newName
}
func extractNumbers(input string) int {
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(input, -1)

	if len(matches) == 0 {
		return 0
	}
	level, err := strconv.Atoi(matches[len(matches)-1])
	if err != nil || level > 16 {
		return 0
	}
	return level
}
