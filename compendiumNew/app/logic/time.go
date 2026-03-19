package logic

import (
	"compendium/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

// GeoNames API configuration
// Для использования API необходимо зарегистрироваться на https://www.geonames.org/
// и получить бесплатный username
const (
	geoNamesUsername = "mentalisit" // Замените на ваш username после регистрации
	geoNamesBaseURL  = "http://api.geonames.org"
)

// GeoNamesSearchResult представляет результат поиска города
type GeoNamesSearchResult struct {
	Geonames []struct {
		Name        string `json:"name"`
		Lat         string `json:"lat"`
		Lng         string `json:"lng"`
		CountryName string `json:"countryName"`
	} `json:"geonames"`
}

// GeoNamesTimezoneResult представляет результат запроса часового пояса
type GeoNamesTimezoneResult struct {
	TimezoneId string  `json:"timezoneId"`
	RawOffset  float64 `json:"rawOffset"`
	DstOffset  float64 `json:"dstOffset"`
	GmtOffset  float64 `json:"gmtOffset"`
}

// Константы для регулярных выражений и форматов
const (
	time12HourFormat = "03:04 PM"
	time24HourFormat = "15:04"
	minutesPerHour   = 60
	secondsPerHour   = 3600
)

var (
	// Обновленный regex: поддерживает названия с пробелами и кириллицу
	// Примеры: "New York", "Los Angeles", "Москва", "Тель Авив"
	tzSetRegex   = regexp.MustCompile(`(?i)%tz set\s*(<@\d+>|@\S+)?\s*([+-]?\d+(\.\d+)?|[\p{L}/_]+(?:\s+[\p{L}]+)*)`)
	mentionRegex = regexp.MustCompile(`<@(\d+)>|@(\S+)`)
	timeCommands = []string{"%time", "%время", "%час"}
)

// cityAliases содержит маппинг популярных названий городов на IANA timezone
var cityAliases = map[string]string{
	// США
	"new york":      "America/New_York",
	"newyork":       "America/New_York",
	"ny":            "America/New_York",
	"los angeles":   "America/Los_Angeles",
	"losangeles":    "America/Los_Angeles",
	"la":            "America/Los_Angeles",
	"chicago":       "America/Chicago",
	"denver":        "America/Denver",
	"phoenix":       "America/Phoenix",
	"seattle":       "America/Los_Angeles",
	"miami":         "America/New_York",
	"boston":        "America/New_York",
	"san francisco": "America/Los_Angeles",
	"sanfrancisco":  "America/Los_Angeles",
	"sf":            "America/Los_Angeles",
	"las vegas":     "America/Los_Angeles",
	"lasvegas":      "America/Los_Angeles",
	"houston":       "America/Chicago",
	"dallas":        "America/Chicago",
	"atlanta":       "America/New_York",

	// Европа
	"london":     "Europe/London",
	"paris":      "Europe/Paris",
	"berlin":     "Europe/Berlin",
	"madrid":     "Europe/Madrid",
	"rome":       "Europe/Rome",
	"amsterdam":  "Europe/Amsterdam",
	"brussels":   "Europe/Brussels",
	"vienna":     "Europe/Vienna",
	"warsaw":     "Europe/Warsaw",
	"prague":     "Europe/Prague",
	"budapest":   "Europe/Budapest",
	"athens":     "Europe/Athens",
	"helsinki":   "Europe/Helsinki",
	"stockholm":  "Europe/Stockholm",
	"oslo":       "Europe/Oslo",
	"copenhagen": "Europe/Copenhagen",
	"dublin":     "Europe/Dublin",
	"lisbon":     "Europe/Lisbon",
	"zurich":     "Europe/Zurich",

	// СНГ и Украина
	"moscow": "Europe/Moscow",
	"москва": "Europe/Moscow",
	"kiev":   "Europe/Kiev",
	"kyiv":   "Europe/Kiev",
	"киев":   "Europe/Kiev",
	"київ":   "Europe/Kiev",
	// Украинские города
	"dnipro":           "Europe/Kiev",
	"dnepr":            "Europe/Kiev",
	"днепр":            "Europe/Kiev",
	"дніпро":           "Europe/Kiev",
	"kharkiv":          "Europe/Kiev",
	"kharkov":          "Europe/Kiev",
	"харьков":          "Europe/Kiev",
	"харків":           "Europe/Kiev",
	"odessa":           "Europe/Kiev",
	"odesa":            "Europe/Kiev",
	"одесса":           "Europe/Kiev",
	"одеса":            "Europe/Kiev",
	"lviv":             "Europe/Kiev",
	"lvov":             "Europe/Kiev",
	"львов":            "Europe/Kiev",
	"львів":            "Europe/Kiev",
	"zaporizhzhia":     "Europe/Kiev",
	"zaporozhye":       "Europe/Kiev",
	"запорожье":        "Europe/Kiev",
	"запоріжжя":        "Europe/Kiev",
	"donetsk":          "Europe/Kiev",
	"донецк":           "Europe/Kiev",
	"донецьк":          "Europe/Kiev",
	"mariupol":         "Europe/Kiev",
	"мариуполь":        "Europe/Kiev",
	"маріуполь":        "Europe/Kiev",
	"vinnitsa":         "Europe/Kiev",
	"vinnytsia":        "Europe/Kiev",
	"винница":          "Europe/Kiev",
	"вінниця":          "Europe/Kiev",
	"poltava":          "Europe/Kiev",
	"полтава":          "Europe/Kiev",
	"chernihiv":        "Europe/Kiev",
	"chernigov":        "Europe/Kiev",
	"чернигов":         "Europe/Kiev",
	"чернігів":         "Europe/Kiev",
	"sumy":             "Europe/Kiev",
	"сумы":             "Europe/Kiev",
	"суми":             "Europe/Kiev",
	"mykolaiv":         "Europe/Kiev",
	"nikolaev":         "Europe/Kiev",
	"николаев":         "Europe/Kiev",
	"миколаїв":         "Europe/Kiev",
	"kherson":          "Europe/Kiev",
	"херсон":           "Europe/Kiev",
	"rivne":            "Europe/Kiev",
	"rovno":            "Europe/Kiev",
	"ровно":            "Europe/Kiev",
	"рівне":            "Europe/Kiev",
	"lutsk":            "Europe/Kiev",
	"луцк":             "Europe/Kiev",
	"луцьк":            "Europe/Kiev",
	"uzhgorod":         "Europe/Kiev",
	"ужгород":          "Europe/Kiev",
	"ivano-frankivsk":  "Europe/Kiev",
	"ivano frankivsk":  "Europe/Kiev",
	"ивано-франковск":  "Europe/Kiev",
	"івано-франківськ": "Europe/Kiev",
	"ternopil":         "Europe/Kiev",
	"тернополь":        "Europe/Kiev",
	"тернопіль":        "Europe/Kiev",
	"khmelnytskyi":     "Europe/Kiev",
	"khmelnitsky":      "Europe/Kiev",
	"хмельницкий":      "Europe/Kiev",
	"хмельницький":     "Europe/Kiev",
	"cherkasy":         "Europe/Kiev",
	"cherkassy":        "Europe/Kiev",
	"черкассы":         "Europe/Kiev",
	"черкаси":          "Europe/Kiev",
	"kropyvnytskyi":    "Europe/Kiev",
	"kirovograd":       "Europe/Kiev",
	"кировоград":       "Europe/Kiev",
	"кропивницький":    "Europe/Kiev",
	"zhytomyr":         "Europe/Kiev",
	"житомир":          "Europe/Kiev",
	// Беларусь
	"minsk":   "Europe/Minsk",
	"минск":   "Europe/Minsk",
	"мінск":   "Europe/Minsk",
	"gomel":   "Europe/Minsk",
	"гомель":  "Europe/Minsk",
	"brest":   "Europe/Minsk",
	"брест":   "Europe/Minsk",
	"grodno":  "Europe/Minsk",
	"гродно":  "Europe/Minsk",
	"vitebsk": "Europe/Minsk",
	"витебск": "Europe/Minsk",
	"mogilev": "Europe/Minsk",
	"могилев": "Europe/Minsk",
	// Казахстан
	"almaty":     "Asia/Almaty",
	"алматы":     "Asia/Almaty",
	"astana":     "Asia/Almaty",
	"астана":     "Asia/Almaty",
	"nur-sultan": "Asia/Almaty",
	"нур-султан": "Asia/Almaty",
	"shymkent":   "Asia/Almaty",
	"шымкент":    "Asia/Almaty",
	"karaganda":  "Asia/Almaty",
	"караганда":  "Asia/Almaty",
	"aktobe":     "Asia/Aqtobe",
	"актобе":     "Asia/Aqtobe",
	"atyrau":     "Asia/Atyrau",
	"атырау":     "Asia/Atyrau",
	// Узбекистан
	"tashkent":  "Asia/Tashkent",
	"ташкент":   "Asia/Tashkent",
	"samarkand": "Asia/Samarkand",
	"самарканд": "Asia/Samarkand",
	// Азербайджан
	"baku": "Asia/Baku",
	"баку": "Asia/Baku",
	// Грузия
	"tbilisi": "Asia/Tbilisi",
	"тбилиси": "Asia/Tbilisi",
	// Армения
	"yerevan": "Asia/Yerevan",
	"ереван":  "Asia/Yerevan",
	// Прибалтика
	"riga":    "Europe/Riga",
	"рига":    "Europe/Riga",
	"tallinn": "Europe/Tallinn",
	"таллин":  "Europe/Tallinn",
	"vilnius": "Europe/Vilnius",
	"вильнюс": "Europe/Vilnius",
	// Россия (дополнительные города)
	"saint petersburg": "Europe/Moscow",
	"st petersburg":    "Europe/Moscow",
	"санкт-петербург":  "Europe/Moscow",
	"питер":            "Europe/Moscow",
	"spb":              "Europe/Moscow",
	"novosibirsk":      "Asia/Novosibirsk",
	"новосибирск":      "Asia/Novosibirsk",
	"yekaterinburg":    "Asia/Yekaterinburg",
	"екатеринбург":     "Asia/Yekaterinburg",
	"kazan":            "Europe/Moscow",
	"казань":           "Europe/Moscow",
	"nizhny novgorod":  "Europe/Moscow",
	"нижний новгород":  "Europe/Moscow",
	"samara":           "Europe/Samara",
	"самара":           "Europe/Samara",
	"omsk":             "Asia/Omsk",
	"омск":             "Asia/Omsk",
	"chelyabinsk":      "Asia/Yekaterinburg",
	"челябинск":        "Asia/Yekaterinburg",
	"rostov":           "Europe/Moscow",
	"ростов":           "Europe/Moscow",
	"ufa":              "Asia/Yekaterinburg",
	"уфа":              "Asia/Yekaterinburg",
	"krasnoyarsk":      "Asia/Krasnoyarsk",
	"красноярск":       "Asia/Krasnoyarsk",
	"voronezh":         "Europe/Moscow",
	"воронеж":          "Europe/Moscow",
	"perm":             "Asia/Yekaterinburg",
	"пермь":            "Asia/Yekaterinburg",
	"volgograd":        "Europe/Volgograd",
	"волгоград":        "Europe/Volgograd",
	"krasnodar":        "Europe/Moscow",
	"краснодар":        "Europe/Moscow",
	"vladivostok":      "Asia/Vladivostok",
	"владивосток":      "Asia/Vladivostok",
	"irkutsk":          "Asia/Irkutsk",
	"иркутск":          "Asia/Irkutsk",

	// Азия
	"tokyo":     "Asia/Tokyo",
	"токио":     "Asia/Tokyo",
	"seoul":     "Asia/Seoul",
	"сеул":      "Asia/Seoul",
	"beijing":   "Asia/Shanghai",
	"пекин":     "Asia/Shanghai",
	"shanghai":  "Asia/Shanghai",
	"шанхай":    "Asia/Shanghai",
	"hong kong": "Asia/Hong_Kong",
	"hongkong":  "Asia/Hong_Kong",
	"singapore": "Asia/Singapore",
	"сингапур":  "Asia/Singapore",
	"bangkok":   "Asia/Bangkok",
	"бангкок":   "Asia/Bangkok",
	"dubai":     "Asia/Dubai",
	"дубай":     "Asia/Dubai",
	"mumbai":    "Asia/Kolkata",
	"мумбаи":    "Asia/Kolkata",
	"delhi":     "Asia/Kolkata",
	"дели":      "Asia/Kolkata",
	"jakarta":   "Asia/Jakarta",
	"джакарта":  "Asia/Jakarta",
	"manila":    "Asia/Manila",
	"манила":    "Asia/Manila",
	"istanbul":  "Europe/Istanbul",
	"стамбул":   "Europe/Istanbul",
	"tel aviv":  "Asia/Jerusalem",
	"telaviv":   "Asia/Jerusalem",
	"тель авив": "Asia/Jerusalem",
	"jerusalem": "Asia/Jerusalem",
	"иерусалим": "Asia/Jerusalem",

	// Австралия и Океания
	"sydney":    "Australia/Sydney",
	"сидней":    "Australia/Sydney",
	"melbourne": "Australia/Melbourne",
	"мельбурн":  "Australia/Melbourne",
	"brisbane":  "Australia/Brisbane",
	"брисбен":   "Australia/Brisbane",
	"perth":     "Australia/Perth",
	"перт":      "Australia/Perth",
	"auckland":  "Pacific/Auckland",
	"окленд":    "Pacific/Auckland",

	// Южная Америка
	"sao paulo":    "America/Sao_Paulo",
	"saopaulo":     "America/Sao_Paulo",
	"buenos aires": "America/Argentina/Buenos_Aires",
	"buenosaires":  "America/Argentina/Buenos_Aires",
	"santiago":     "America/Santiago",
	"lima":         "America/Lima",
	"bogota":       "America/Bogota",

	// Африка
	"cairo":        "Africa/Cairo",
	"каир":         "Africa/Cairo",
	"johannesburg": "Africa/Johannesburg",
	"cape town":    "Africa/Johannesburg",
	"capetown":     "Africa/Johannesburg",
	"lagos":        "Africa/Lagos",
	"nairobi":      "Africa/Nairobi",

	// Канада
	"toronto":   "America/Toronto",
	"торонто":   "America/Toronto",
	"vancouver": "America/Vancouver",
	"ванкувер":  "America/Vancouver",
	"montreal":  "America/Montreal",
	"монреаль":  "America/Montreal",
	"winnipeg":  "America/Winnipeg",
	"виннипег":  "America/Winnipeg",
	"calgary":   "America/Edmonton",
	"калгари":   "America/Edmonton",
	"edmonton":  "America/Edmonton",
	"эдмонтон":  "America/Edmonton",
	"ottawa":    "America/Toronto",
	"оттава":    "America/Toronto",
}

// TimezoneInfo содержит информацию о часовом поясе
type TimezoneInfo struct {
	Name       string
	Offset     int    // смещение в минутах (текущее)
	Location   string // название локации для DST (например, "America/New_York")
	IsLocation bool   // true если timezone задан как локация (поддерживает DST)
}

// TzTime обрабатывает все команды, связанные с часовыми поясами
func (c *Hs) TzTime(m models.IncomingMessage) bool {
	switch {
	case c.TzTimeSet(m):
		return true
	case c.TzGet(m):
		return true
	case c.TzGetTime(m):
		return true
	}
	return false
}

// TzTimeSet устанавливает часовой пояс для пользователя
func (c *Hs) TzTimeSet(m models.IncomingMessage) bool {
	matches := tzSetRegex.FindStringSubmatch(m.Text)
	if len(matches) < 3 {
		return false
	}

	optionalParam := matches[1]
	timezone := matches[2]

	tzInfo, err := c.parseTimezone(timezone)
	if err != nil {
		text := fmt.Sprintf(c.getText(m, "I_COULD_NOT_FIND_ANY"), m.MentionName, timezone)
		c.sendChat(m, text)
		return true
	}

	c.setTimezone(tzInfo, optionalParam, m)
	return true
}

// parseTimezone парсит строку часового пояса и возвращает TimezoneInfo
func (c *Hs) parseTimezone(timezone string) (*TimezoneInfo, error) {
	// Проверяем числовое смещение
	if offset, err := strconv.ParseFloat(timezone, 64); err == nil {
		return &TimezoneInfo{
			Name:       formatTimezoneName(offset),
			Offset:     int(offset * minutesPerHour),
			Location:   "",
			IsLocation: false,
		}, nil
	}

	// Нормализуем ввод для поиска в алиасах (lowercase, trim)
	normalizedInput := strings.ToLower(strings.TrimSpace(timezone))

	// Проверяем алиасы городов (например, "New York" -> "America/New_York")
	if ianaTimezone, ok := cityAliases[normalizedInput]; ok {
		timezone = ianaTimezone
	}

	// Пытаемся загрузить местоположение
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		// Если не нашли в IANA базе, пробуем через GeoNames API
		geoTimezone, geoErr := lookupTimezoneByCity(normalizedInput)
		if geoErr != nil {
			return nil, err // Возвращаем оригинальную ошибку
		}
		timezone = geoTimezone
		loc, err = time.LoadLocation(timezone)
		if err != nil {
			return nil, err
		}
	}

	now := time.Now().In(loc)
	_, offsetSeconds := now.Zone()
	offsetHours := float64(offsetSeconds) / secondsPerHour

	// Сохраняем название локации для поддержки DST
	// При отображении времени будем использовать Location для динамического расчета
	return &TimezoneInfo{
		Name:       timezone, // Сохраняем IANA название локации (например, "America/New_York")
		Offset:     int(offsetHours * minutesPerHour),
		Location:   timezone,
		IsLocation: true,
	}, nil
}

// formatTimezoneName форматирует название часового пояса
func formatTimezoneName(offset float64) string {
	if offset >= 0 {
		return fmt.Sprintf("UTC+%v", offset)
	}
	return fmt.Sprintf("UTC%v", offset)
}

// setTimezone устанавливает часовой пояс для указанного пользователя
func (c *Hs) setTimezone(tzInfo *TimezoneInfo, mentionName string, m models.IncomingMessage) {
	// Обновляем timezone для текущего аккаунта, если есть
	c.updateCurrentAccountTimezone(m, tzInfo)

	if mentionName == "" {
		c.setTimezoneForSelf(tzInfo, m)
	} else {
		c.setTimezoneForMention(tzInfo, mentionName, m)
	}
}

// updateCurrentAccountTimezone обновляет timezone для текущего мульти-аккаунта
func (c *Hs) updateCurrentAccountTimezone(m models.IncomingMessage, tzInfo *TimezoneInfo) {
	if m.MAcc == nil {
		return
	}

	member, err := c.db.V2.CorpMemberByUId(m.MAcc.UUID)
	if err != nil || member == nil {
		return
	}

	member.TimeZona = tzInfo.Name
	member.ZonaOffset = tzInfo.Offset
	if err := c.db.V2.CorpMemberUpdate(*member); err != nil {
		c.log.ErrorErr(err)
	}
}

// setTimezoneForSelf устанавливает timezone для самого пользователя
func (c *Hs) setTimezoneForSelf(tzInfo *TimezoneInfo, m models.IncomingMessage) {
	var err error

	if m.MAcc != nil {
		member, err := c.DbV2.CorpMemberByUId(m.MAcc.UUID)
		if member != nil && err == nil {
			member.TimeZona = tzInfo.Name
			member.ZonaOffset = tzInfo.Offset
			err = c.db.V2.CorpMemberUpdate(*member)
		}
	} else {
		err = c.corpMember.CorpMemberTZUpdate(m.NameId, m.MGuild.GuildId(), tzInfo.Name, tzInfo.Offset)
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.createNewMemberWithTimezone(m.NameId, tzInfo, m)
		} else {
			c.log.ErrorErr(err)
			return
		}
	}

	text := fmt.Sprintf(c.getText(m, "TIMEZONA_SET"), m.MentionName, m.Name, tzInfo.Name)
	c.sendChat(m, text)
}

// setTimezoneForMention устанавливает timezone для упомянутого пользователя
func (c *Hs) setTimezoneForMention(tzInfo *TimezoneInfo, mentionName string, m models.IncomingMessage) {
	matches := mentionRegex.FindStringSubmatch(mentionName)
	if len(matches) < 3 {
		return
	}

	var targetUserID string
	var displayName string

	if matches[1] != "" {
		// Discord mention: <@123456>
		targetUserID = matches[1]
		displayName = mentionName
		c.updateMentionedAccountTimezone(m, matches[1], tzInfo)
	} else if matches[2] != "" {
		// Telegram mention: @username
		user, err := c.users.UsersGetByUserName(matches[2])
		if err != nil {
			c.log.ErrorErr(err)
			return
		}
		targetUserID = user.ID
		displayName = matches[2]
		c.updateTelegramAccountTimezone(m, matches[2], tzInfo)
	}

	if targetUserID == "" {
		return
	}

	err := c.corpMember.CorpMemberTZUpdate(targetUserID, m.MGuild.GuildId(), tzInfo.Name, tzInfo.Offset)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.createNewMemberWithTimezone(targetUserID, tzInfo, m)
		} else {
			c.log.ErrorErr(err)
			return
		}
	}

	text := fmt.Sprintf(c.getText(m, "TIMEZONA_SET"), m.MentionName, displayName, tzInfo.Name)
	c.sendChat(m, text)
}

// updateMentionedAccountTimezone обновляет timezone для упомянутого Discord аккаунта
func (c *Hs) updateMentionedAccountTimezone(m models.IncomingMessage, discordID string, tzInfo *TimezoneInfo) {
	if m.MAcc == nil || discordID != m.MAcc.DiscordID {
		return
	}

	corpMember, err := c.db.V2.CorpMemberByUId(m.MAcc.UUID)
	if err != nil || corpMember == nil {
		return
	}

	corpMember.TimeZona = tzInfo.Name
	corpMember.ZonaOffset = tzInfo.Offset
	if err := c.db.V2.CorpMemberUpdate(*corpMember); err != nil {
		c.log.ErrorErr(err)
	}
}

// updateTelegramAccountTimezone обновляет timezone для упомянутого Telegram аккаунта
func (c *Hs) updateTelegramAccountTimezone(m models.IncomingMessage, telegramUsername string, tzInfo *TimezoneInfo) {
	if m.MAcc == nil || telegramUsername != m.MAcc.TelegramUsername {
		return
	}

	corpMember, err := c.db.V2.CorpMemberByUId(m.MAcc.UUID)
	if err != nil || corpMember == nil {
		return
	}

	corpMember.TimeZona = tzInfo.Name
	corpMember.ZonaOffset = tzInfo.Offset
	if err := c.db.V2.CorpMemberUpdate(*corpMember); err != nil {
		c.log.ErrorErr(err)
	}
}

// createNewMemberWithTimezone создает нового члена корпорации с timezone
func (c *Hs) createNewMemberWithTimezone(userID string, tzInfo *TimezoneInfo, m models.IncomingMessage) {
	if m.MAcc != nil {
		return
	}

	cm := models.CorpMember{
		Name:       m.Name,
		UserId:     userID,
		GuildId:    m.MGuild.GuildId(),
		Tech:       models.TechLevels{},
		AvatarUrl:  m.Avatar,
		TimeZone:   tzInfo.Name,
		ZoneOffset: tzInfo.Offset,
		MGuild:     m.MGuild,
		MAcc:       m.MAcc,
	}
	_ = c.corpMember.CorpMemberInsert(cm)

	u := models.User{
		ID:        userID,
		Username:  m.Name,
		AvatarURL: m.Avatar,
		Alts:      []string{},
	}
	_ = c.users.UsersInsert(u)
}

// TzGet возвращает текущий часовой пояс пользователя
func (c *Hs) TzGet(m models.IncomingMessage) bool {
	if !strings.HasPrefix(m.Text, "%tz get") {
		return false
	}

	timezone := c.getUserTimezone(m)
	text := fmt.Sprintf(c.getText(m, "TIMEZONA_IS_CURRENTLY"), m.MentionName, m.Name, timezone)
	c.sendChat(m, text)
	return true
}

// getUserTimezone получает timezone пользователя
func (c *Hs) getUserTimezone(m models.IncomingMessage) string {
	// Сначала проверяем мульти-аккаунт
	if m.MAcc != nil {
		corpMember, err := c.db.V2.CorpMemberByUId(m.MAcc.UUID)
		if err == nil && corpMember != nil {
			return corpMember.TimeZona
		}
	}

	// Затем проверяем обычных членов
	members, err := c.corpMember.CorpMembersRead(m.MGuild.GuildId())
	if err != nil {
		c.log.ErrorErr(err)
		return "Not set"
	}

	for _, member := range members {
		if member.UserId == m.NameId {
			return member.TimeZone
		}
	}

	return "Not set"
}

// TzGetTime отображает локальное время для всех членов корпорации
func (c *Hs) TzGetTime(m models.IncomingMessage) bool {
	if !isTimeCommand(m.Text) {
		return false
	}

	members := c.collectAllMembers(m)
	data := c.buildTimeTable(members)

	text := fmt.Sprintf(c.getText(m, "LOCAL_TIME_FOR_EVERYONE"), m.MentionName)
	c.sendFormatedText(m, text, data)
	c.sendChat(m, c.getText(m, "UNLISTED_MEMBERS"))
	return true
}

// isTimeCommand проверяет, является ли текст командой времени
func isTimeCommand(text string) bool {
	for _, cmd := range timeCommands {
		if strings.HasPrefix(text, cmd) {
			return true
		}
	}
	return false
}

// collectAllMembers собирает всех членов корпорации из разных источников
func (c *Hs) collectAllMembers(m models.IncomingMessage) map[string]models.CorpMember {
	memberMap := make(map[string]models.CorpMember)

	// Получаем обычных членов
	members, err := c.corpMember.CorpMembersRead(m.MGuild.GuildId())
	if err != nil {
		c.log.ErrorErr(err)
	}

	for _, member := range members {
		c.enrichMemberWithLocalTime(&member)
		memberMap[member.Name] = member
	}

	// Получаем мульти-аккаунты, если есть гильдия
	if m.MGuild != nil {
		c.collectMultiAccountMembers(m, memberMap)
	}

	return memberMap
}

// collectMultiAccountMembers собирает членов из мульти-аккаунтов
func (c *Hs) collectMultiAccountMembers(m models.IncomingMessage, memberMap map[string]models.CorpMember) {
	// Первый источник мульти-аккаунтов
	membersRead, _ := c.DbV2.CorpMembersReadMulti(&m.MGuild.GId)
	for _, member := range membersRead {
		accountUUID, _ := c.DbV2.FindMultiAccountUUID(member.MAcc.UUID)
		if accountUUID == nil {
			continue
		}

		corpMember := models.CorpMember{
			Name:       accountUUID.Nickname,
			UserId:     member.MAcc.UUID.String(),
			AvatarUrl:  accountUUID.AvatarURL,
			TimeZone:   member.TimeZone,
			ZoneOffset: member.ZoneOffset,
			AfkFor:     member.AfkFor,
		}
		c.enrichMemberWithLocalTime(&corpMember)
		memberMap[corpMember.Name] = corpMember
	}

	// Второй источник мульти-аккаунтов
	memberS, _ := c.db.V2.CorpMembersReadMulti(&m.MGuild.GId)
	for _, member := range memberS {
		corpMember := models.CorpMember{
			Name:       member.MAcc.Nickname,
			UserId:     member.MAcc.UUID.String(),
			AvatarUrl:  member.AvatarUrl,
			TimeZone:   member.TimeZone,
			ZoneOffset: member.ZoneOffset,
			AfkFor:     member.AfkFor,
		}
		c.enrichMemberWithLocalTime(&corpMember)
		memberMap[corpMember.Name] = corpMember
	}
}

// enrichMemberWithLocalTime добавляет локальное время к члену корпорации
// Поддерживает DST: если TimeZone содержит название локации (например, "America/New_York"),
// время вычисляется динамически с учетом текущего летнего/зимнего времени
func (c *Hs) enrichMemberWithLocalTime(member *models.CorpMember) {
	if member.TimeZone == "" {
		return
	}
	// Используем функцию с поддержкой DST
	// TimeZone может содержать либо название локации (America/New_York), либо UTC offset (UTC+3)
	member.LocalTime, member.LocalTime24 = getTimeStringsWithDST(member.TimeZone, member.ZoneOffset)
}

// buildTimeTable строит таблицу времени для отображения
func (c *Hs) buildTimeTable(memberMap map[string]models.CorpMember) [][]string {
	data := [][]string{
		{"Local Time", "User", ""},
	}

	for _, member := range memberMap {
		if member.TimeZone != "" {
			newRow := []string{member.LocalTime24, member.Name, ""}
			data = append(data, newRow)
		}
	}

	return data
}

// getTimeStrings возвращает время в 12-часовом и 24-часовом форматах
// Использует фиксированное смещение (для обратной совместимости)
func getTimeStrings(offsetMinutes int) (string, string) {
	now := time.Now().UTC()
	offsetDuration := time.Duration(offsetMinutes) * time.Minute
	timeWithOffset := now.Add(offsetDuration)

	return timeWithOffset.Format(time12HourFormat), timeWithOffset.Format(time24HourFormat)
}

// getTimeStringsWithDST возвращает время с учетом DST (летнего/зимнего времени)
// Если timezone - это название локации (например, "America/New_York"),
// то смещение вычисляется динамически с учетом текущего DST
func getTimeStringsWithDST(timezone string, fallbackOffsetMinutes int) (string, string) {
	now := time.Now()

	// Пытаемся загрузить локацию для динамического расчета DST
	if timezone != "" {
		loc, err := time.LoadLocation(timezone)
		if err == nil {
			timeInLocation := now.In(loc)
			return timeInLocation.Format(time12HourFormat), timeInLocation.Format(time24HourFormat)
		}
	}

	// Fallback на фиксированное смещение
	return getTimeStrings(fallbackOffsetMinutes)
}

// getCurrentOffset возвращает текущее смещение для локации с учетом DST
func getCurrentOffset(timezone string) (int, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return 0, err
	}

	now := time.Now().In(loc)
	_, offsetSeconds := now.Zone()
	return offsetSeconds / 60, nil // возвращаем в минутах
}

// lookupTimezoneByCity ищет часовой пояс по названию города через GeoNames API
// Поддерживает любые города мира, включая "Днепр", "Харьков" и т.д.
func lookupTimezoneByCity(cityName string) (string, error) {
	// Шаг 1: Поиск координат города
	lat, lng, err := searchCityCoordinates(cityName)
	if err != nil {
		return "", err
	}

	// Шаг 2: Получение часового пояса по координатам
	return getTimezoneByCoordinates(lat, lng)
}

// searchCityCoordinates ищет координаты города через GeoNames Search API
func searchCityCoordinates(cityName string) (string, string, error) {
	// Формируем URL для поиска
	searchURL := fmt.Sprintf("%s/searchJSON?q=%s&maxRows=1&username=%s&featureClass=P",
		geoNamesBaseURL,
		url.QueryEscape(cityName),
		geoNamesUsername,
	)

	// Выполняем HTTP запрос
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(searchURL)
	if err != nil {
		return "", "", fmt.Errorf("failed to search city: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("GeoNames API returned status: %d", resp.StatusCode)
	}

	// Парсим ответ
	var result GeoNamesSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", fmt.Errorf("failed to decode search response: %w", err)
	}

	if len(result.Geonames) == 0 {
		return "", "", fmt.Errorf("city not found: %s", cityName)
	}

	return result.Geonames[0].Lat, result.Geonames[0].Lng, nil
}

// getTimezoneByCoordinates получает IANA timezone по координатам через GeoNames Timezone API
func getTimezoneByCoordinates(lat, lng string) (string, error) {
	// Формируем URL для получения timezone
	tzURL := fmt.Sprintf("%s/timezoneJSON?lat=%s&lng=%s&username=%s",
		geoNamesBaseURL,
		lat,
		lng,
		geoNamesUsername,
	)

	// Выполняем HTTP запрос
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(tzURL)
	if err != nil {
		return "", fmt.Errorf("failed to get timezone: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GeoNames Timezone API returned status: %d", resp.StatusCode)
	}

	// Парсим ответ
	var result GeoNamesTimezoneResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode timezone response: %w", err)
	}

	if result.TimezoneId == "" {
		return "", fmt.Errorf("timezone not found for coordinates: %s, %s", lat, lng)
	}

	return result.TimezoneId, nil
}
