package helpers

import (
	"fmt"
	"github.com/mentalisit/logger"
	"regexp"
	"rs/models"
	"rs/storage"
	"strconv"
	"strings"
)

const ds = "ds"
const tg = "tg"

type Helpers struct {
	log     *logger.Logger
	storage *storage.Storage
	//saveArray []SaveDM
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
func (h *Helpers) emReadName(s models.Sborkz, ForType string, mention ...bool) string { // склеиваем имя и эмоджи
	name := s.Name
	if s.Wamesid != "" {
		name = s.Wamesid
	}
	//multiAccount, _ := h.storage.Postgres.FindMultiAccountByUserId(s.UserId)
	//if multiAccount != nil {
	//	emojiReadUUID := h.storage.Postgres.EmojiReadUUID(multiAccount.UUID, ForType)
	//}

	newName := s.Name
	if ForType == ds {
		newName = s.Mention
	} else {
		newName = s.Name
	}

	t := h.storage.Emoji.EmojiModuleReadUsers(name, ForType)

	if mention != nil {
		if mention[0] {
			newName = s.Mention
		}
	}

	if len(t.Name) > 0 {
		////nickName
		//if t.Weapon != "" && s.Tip == "tg" {
		//	newName = fmt.Sprintf("%s [%s]", newName, t.Weapon)
		//}
		//Alt
		if s.Wamesid != "" {
			newName = fmt.Sprintf("%s [%s]", newName, s.Wamesid)
		}
		if t.Module1 != "" {
			if ForType == ds {
				newName = fmt.Sprintf("%s %s %s %s %s", newName, t.Module1, t.Module2, t.Module3, t.Weapon)
			} else if ForType == tg {
				newName = fmt.Sprintf("%s (%s/%s/%s)", newName, t.Module1, t.Module2, t.Module3)
			}
		}
		newName = fmt.Sprintf("%s %s%s%s%s", newName, t.Em1, t.Em2, t.Em3, t.Em4)
	}
	return newName
}

func (h *Helpers) ReadNameModules(in models.InMessage, name string) {
	if name == "" {
		multiAccount, _ := h.storage.Postgres.FindMultiAccountByUserId(in.UserId)
		if multiAccount != nil {
			h.ReadNameModulesUUID(in, name)
			//return
		}
		name = in.Username
	}
	var tds, ttg models.EmodjiUser
	var DsGenesis, DsEnrich, DsRsExtender int
	var TgGenesis, TgEnrich, TgRsExtender int
	var genesisA, enrichA, rseA int
	if in.IfDiscord() {
		genesis1, enrich1, rsextender1 := 0, 0, 0
		//if name == in.Username {
		//	genesis1, enrich1, rsextender1 = GetTechDataUserId(in.UserId, in.Config.Guildid)
		//}
		genesis2, enrich2, rsextender2 := Get2TechDataUserId(name, in.UserId, in.Ds.Guildid)

		DsGenesis = max(genesis1, genesis2)
		DsEnrich = max(enrich1, enrich2)
		DsRsExtender = max(rsextender1, rsextender2)
	}
	if in.IfTelegram() {
		split := strings.Split(in.Config.TgChannel, "/")
		TgGenesis, TgEnrich, TgRsExtender = Get2TechDataUserId(name, in.UserId, split[0])
	}
	if in.Tip == tg {
		genesisA, enrichA, rseA = Get3TechDataUserId(name, in.UserId)
		if genesisA != 0 || enrichA != 0 || rseA != 0 {
			fmt.Printf("use Get3TechDataUserId for %s %s\n", name, in.UserId)
		}
	}

	genesis := max(DsGenesis, TgGenesis, genesisA)
	enrich := max(DsEnrich, TgEnrich, enrichA)
	rsextender := max(DsRsExtender, TgRsExtender, rseA)

	fmt.Printf("genesis %d enrich %d rsextender %d for:%s\n", genesis, enrich, rsextender, name)

	if in.IfDiscord() {
		tds = h.storage.Emoji.EmojiModuleReadUsers(name, ds)
		if tds.Name == "" {
			h.storage.Emoji.EmInsertEmpty(ds, name)
		}

		one := fmt.Sprintf("<:genesis:1199068748280242237> %d ", genesis)
		two := fmt.Sprintf("<:enrich:1199068793633251338> %d ", enrich)
		three := fmt.Sprintf("<:rse:1199068829511335946> %d ", rsextender)
		if genesis != 0 && tds.Module1 != one {
			h.storage.Emoji.ModuleUpdate(name, "ds", "1", one)
		}
		if enrich != 0 && tds.Module2 != two {
			h.storage.Emoji.ModuleUpdate(name, "ds", "2", two)
		}
		if rsextender != 0 && tds.Module3 != three {
			h.storage.Emoji.ModuleUpdate(name, "ds", "3", three)
		}
	}

	if in.IfTelegram() {
		ttg = h.storage.Emoji.EmojiModuleReadUsers(name, tg)
		if ttg.Name == "" {
			h.storage.Emoji.EmInsertEmpty(tg, name)
		}

		gen := strconv.Itoa(genesis)
		enr := strconv.Itoa(enrich)
		rse := strconv.Itoa(rsextender)
		if genesis != 0 && ttg.Module1 != gen {
			h.storage.Emoji.ModuleUpdate(name, "tg", "1", gen)
		}
		if enrich != 0 && ttg.Module2 != enr {
			h.storage.Emoji.ModuleUpdate(name, "tg", "2", enr)
		}
		if rsextender != 0 && ttg.Module1 != rse {
			h.storage.Emoji.ModuleUpdate(name, "tg", "3", rse)
		}
	}
}

func (h *Helpers) NameMention(u models.Users, tip string) (n1, n2, n3, n4 string) {
	if u.User1.Tip == tip {
		n1 = h.emReadName(u.User1, tip, true)
	} else {
		n1 = h.emReadName(u.User1, tip)
	}
	if u.User2 != nil {
		if u.User2.Tip == tip {
			n2 = h.emReadName(*u.User2, tip, true)
		} else {
			n2 = h.emReadName(*u.User2, tip)
		}
	}
	if u.User3 != nil {
		if u.User3.Tip == tip {
			n3 = h.emReadName(*u.User3, tip, true)
		} else {
			n3 = h.emReadName(*u.User3, tip)
		}
	}
	if u.User4 != nil {
		if u.User4.Tip == tip {
			n4 = h.emReadName(*u.User4, tip, true)
		} else {
			n4 = h.emReadName(*u.User4, tip)
		}
	}

	return
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
