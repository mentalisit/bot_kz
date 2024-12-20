package helpers

import (
	"fmt"
	"github.com/mentalisit/logger"
	"regexp"
	"rs/models"
	"rs/storage"
	"strconv"
)

const ds = "ds"
const tg = "tg"

type Helpers struct {
	log       *logger.Logger
	storage   *storage.Storage
	saveArray []SaveDM
	//Gemini    *Gemini
}

func NewHelpers(log *logger.Logger, storage *storage.Storage) *Helpers {
	return &Helpers{
		log:     log,
		storage: storage,
		//Gemini:  NewGemini(log),
	}
}
func (h *Helpers) GetQueueDiscord(n map[string]string, u models.Users) map[string]string {
	if u.User1.Name != "" {
		n["name1"] = fmt.Sprintf("%s  🕒  %d  (%d)", h.emReadName(u.User1, ds), u.User1.Timedown, u.User1.Numkzn)
	}
	if u.User2 != nil && u.User2.Name != "" {
		n["name2"] = fmt.Sprintf("%s  🕒  %d  (%d)", h.emReadName(*u.User2, ds), u.User2.Timedown, u.User2.Numkzn)
	}
	if u.User3 != nil && u.User3.Name != "" {
		n["name3"] = fmt.Sprintf("%s  🕒  %d  (%d)", h.emReadName(*u.User3, ds), u.User3.Timedown, u.User3.Numkzn)
	}
	textcount := ""
	if u.User3 != nil && u.User3.Name != "" {
		textcount = fmt.Sprintf("\n1️⃣ %s \n2️⃣ %s \n3️⃣ %s \n\n",
			n["name1"], n["name2"], n["name3"])
	} else if u.User2 != nil && u.User2.Name != "" {
		textcount = fmt.Sprintf("\n1️⃣ %s \n2️⃣ %s \n\n",
			n["name1"], n["name2"])
	} else if u.User1.Name != "" {
		textcount = fmt.Sprintf("\n1️⃣ %s \n\n",
			n["name1"])
	}
	n["textcount"] = textcount
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
	if u.User2 != nil && u.User2.Name != "" {
		users += fmt.Sprintf("2️⃣ %s - %d%s (%d) \n",
			h.emReadName(*u.User2, tg), u.User2.Timedown, n["min"], u.User2.Numkzn)
	}
	if u.User3 != nil && u.User3.Name != "" {
		users += fmt.Sprintf("3️⃣ %s - %d%s (%d) \n",
			h.emReadName(*u.User3, tg), u.User3.Timedown, n["min"], u.User3.Numkzn)
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
	t := h.storage.Emoji.EmojiModuleReadUsers(name, tip)
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
			name = in.Username
		}
		t := h.storage.Emoji.EmojiModuleReadUsers(name, ds)

		genesis1, enrich1, rsextender1 := 0, 0, 0
		if name == in.Username {
			genesis1, enrich1, rsextender1 = GetTechDataUserId(in.UserId, in.Config.Guildid)
		}
		genesis2, enrich2, rsextender2 := Get2TechDataUserId(name, in.UserId, in.Ds.Guildid)

		genesis := max(genesis1, genesis2)
		if genesis == 0 {
			genesis = extractNumbers(t.Module2)
		}

		enrich := max(enrich1, enrich2)
		if enrich == 0 {
			enrich = extractNumbers(t.Module3)
		}

		rsextender := max(rsextender1, rsextender2)
		if rsextender == 0 {
			rsextender = extractNumbers(t.Module1)
		}

		fmt.Printf("genesis %d enrich %d rsextender %d for:%s\n", genesis, enrich, rsextender, name)

		if t.Name == "" {
			h.storage.Emoji.EmInsertEmpty("ds", name)
		}
		one := fmt.Sprintf("<:rse:1199068829511335946> %d ", rsextender)
		two := fmt.Sprintf("<:genesis:1199068748280242237> %d ", genesis)
		three := fmt.Sprintf("<:enrich:1199068793633251338> %d ", enrich)
		if rsextender != 0 && t.Module1 != one {
			h.storage.Emoji.ModuleUpdate(name, "ds", "1", one)
		}
		if genesis != 0 && t.Module2 != two {
			h.storage.Emoji.ModuleUpdate(name, "ds", "2", two)
		}
		if enrich != 0 && t.Module3 != three {
			h.storage.Emoji.ModuleUpdate(name, "ds", "3", three)
		}
	}
}
func (h *Helpers) UpdateCompendiumModules(in models.InMessage) string {
	genesis, enrich, rsextender := GetTechDataUserId(in.UserId, in.Config.Guildid)
	if genesis+enrich+rsextender == 0 {
		genesis, enrich, rsextender = Get2TechDataUserId(in.Username, in.UserId, in.Ds.Guildid)
		if genesis+enrich+rsextender == 0 {
			return "модули не найдены "
		}
	}

	one := fmt.Sprintf("<:rse:1199068829511335946> %d ", rsextender)
	two := fmt.Sprintf("<:genesis:1199068748280242237> %d ", genesis)
	three := fmt.Sprintf("<:enrich:1199068793633251338> %d ", enrich)
	if rsextender != 0 {
		h.storage.Emoji.ModuleUpdate(in.Username, "ds", "1", one)
	}
	if genesis != 0 {
		h.storage.Emoji.ModuleUpdate(in.Username, "ds", "2", two)
	}
	if enrich != 0 {
		h.storage.Emoji.ModuleUpdate(in.Username, "ds", "3", three)
	}
	return " загружено из компендиум бота " + one + two + three
}

func (h *Helpers) NameMention(u models.Users, tip string) (n1, n2, n3, n4 string) {
	if u.User1.Tip == tip {
		n1 = h.emReadMention(u.User1, tip)
	} else {
		n1 = u.User1.Name
	}
	if u.User2 != nil {
		if u.User2.Tip == tip {
			n2 = h.emReadMention(*u.User2, tip)
		} else {
			n2 = u.User2.Name
		}
	}
	if u.User3 != nil {
		if u.User3.Tip == tip {
			n3 = h.emReadMention(*u.User3, tip)
		} else {
			n3 = u.User3.Name
		}
	}
	if u.User4 != nil {
		if u.User4.Tip == tip {
			n4 = h.emReadMention(*u.User4, tip)
		} else {
			n4 = u.User4.Name
		}
	}

	return
}
func (h *Helpers) emReadMention(u models.Sborkz, tip string) string { // склеиваем имя и эмоджи
	t := models.EmodjiUser{}
	if u.Wamesid != "" {
		t = h.storage.Emoji.EmojiModuleReadUsers(u.Wamesid, tip)
	} else {
		t = h.storage.Emoji.EmojiModuleReadUsers(u.Name, tip)
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
