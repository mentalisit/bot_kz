package helpers

import (
	"fmt"
	"rs/models"
	"rs/storage"
	"strings"

	"github.com/mentalisit/logger"
)

const ds = "ds"
const tg = "tg"

type Helpers struct {
	log     *logger.Logger
	storage *storage.Storage
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

	newName := s.Name
	if ForType == ds {
		newName = s.Mention
	} else {
		newName = s.Name
	}

	if mention != nil {
		if mention[0] {
			if s.Mention == "@" {
				newName = fmt.Sprintf("[%s](tg://user?id=%s)", s.Name, s.UserId)
			} else {
				newName = s.Mention
			}
		}
	}
	if s.Wamesid != "" {
		newName = fmt.Sprintf("%s [%s]", newName, s.Wamesid)
	}

	multiAccount, _ := h.storage.Postgres.FindMultiAccountByUserId(s.UserId)
	if multiAccount == nil {

		t := h.storage.Emoji.EmojiModuleReadUsers(name, ForType)

		if len(t.Name) > 0 {
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
	if name == s.Name {
		name = multiAccount.Nickname
	}
	module := h.storage.Postgres.ModuleReadUUID(multiAccount.UUID, name)
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
			newName = fmt.Sprintf("%s (%d/%d/%d)", newName, module.Gen, module.Enr, module.Rse)
		}
	}

	emoji := h.storage.Postgres.EmojiReadUUID(multiAccount.UUID, ForType)
	if emoji != nil && emoji.Tip == ForType {
		newName = fmt.Sprintf("%s %s %s %s %s", newName, emoji.Em1, emoji.Em2, emoji.Em3, emoji.Em4)
	}

	return newName
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
func (h *Helpers) ReadNameModules(in models.InMessage, name string) {
	multiAccount, _ := h.storage.Postgres.FindMultiAccountByUserId(in.UserId)
	if multiAccount != nil {
		if name == "" {
			if multiAccount.Nickname != "" {
				name = multiAccount.Nickname
			} else {
				name = in.Username
			}
		}
		genesis, enrich, rsextender := h.Get2TechDataUserId(name, in.UserId, multiAccount)
		if genesis != 0 || enrich != 0 || rsextender != 0 {
			return
		}
	}
	genesis, enrich, rsextender := 0, 0, 0
	if in.Ds.Guildid != "" {
		genesis, enrich, rsextender = Get2TechDataUserId(name, in.UserId, in.Ds.Guildid)
	} else {
		split := strings.Split(in.Config.TgChannel, "/")
		genesis, enrich, rsextender = Get2TechDataUserId(name, in.UserId, split[0])
	}
	if genesis != 0 && enrich != 0 && rsextender != 0 {
		genesis, enrich, rsextender = Get3TechDataUserId(name, in.UserId)
		if genesis != 0 || enrich != 0 || rsextender != 0 {
			fmt.Printf("use Get3TechDataUserId for %s %s\n", name, in.UserId)
		}
	}

	fmt.Printf("genesis %d enrich %d rsextender %d for:%s\n", genesis, enrich, rsextender, name)
	if multiAccount == nil {
		if name == "" {
			name = in.Username
		}
		multiAccount, _ = h.storage.Postgres.CreateMultiAccountWithPlatform(in.UserId, name, in.Tip, name)
	}

	module := h.storage.Postgres.ModuleReadUUID(multiAccount.UUID, name)
	mod := models.Module{
		Uid:  multiAccount.UUID,
		Name: name,
		Gen:  genesis,
		Enr:  enrich,
		Rse:  rsextender,
	}
	if module == nil {
		h.storage.Postgres.ModuleInsertUUID(mod)
	} else if mod != *module {
		h.storage.Postgres.ModuleUpdateUUID(mod)
	}

	moveUpdateEmoji := func(name, tip string) {
		t := h.storage.Emoji.EmojiModuleReadUsers(name, tip)
		emojiReadUUID := h.storage.Postgres.EmojiReadUUID(multiAccount.UUID, tip)
		if emojiReadUUID == nil {
			h.storage.Postgres.EmojiInsertEmptyUUID(multiAccount.UUID, tip)
		}
		if t.Em1 != "" {
			h.storage.Postgres.EmojiUpdateUUID(multiAccount.UUID, tip, "1", t.Em1)
		}
		if t.Em2 != "" {
			h.storage.Postgres.EmojiUpdateUUID(multiAccount.UUID, tip, "2", t.Em2)
		}
		if t.Em3 != "" {
			h.storage.Postgres.EmojiUpdateUUID(multiAccount.UUID, tip, "3", t.Em3)
		}
		if t.Em4 != "" {
			h.storage.Postgres.EmojiUpdateUUID(multiAccount.UUID, tip, "4", t.Em4)
		}
	}

	if in.IfDiscord() {
		moveUpdateEmoji(name, ds)
	}

	if in.IfTelegram() {
		moveUpdateEmoji(name, tg)
	}
}
