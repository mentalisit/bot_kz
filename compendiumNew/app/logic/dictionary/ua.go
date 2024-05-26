package dictionary

func getDictionaryUkJson() []byte {
	return []byte(`{"uk":{
"HELP_TEXT_DS": "на даний момент доступні команди:\n '**%допомога**' для отримання поточної довідки \n'**%підключити**' для підключення програми\n '**%т і**' для отримання зображення з вашими модулями \n '**%т @ім'я і**' для отримання зображення з модулями іншого гравця\n'**%т ім'я і**' для отримання зображення з модулями альта\n '**%alts add NameAlt**' для створення альта для технологій\n'**%alts del NameAlt**' для видалення альта\n",
"HELP_TEXT_TG": "на даний момент доступні команди:\n '%допомога' для отримання поточної довідки \n'%підключити' для підключення програми\n '%т і' для отримання зображення з вашими модулями \n '%т @ім'я і' для отримання зображення з модулями іншого гравця\n'%т ім'я і' для отримання зображення з модулями альта\n '%alts add NameAlt' для створення альта для технологій\n'%alts del NameAlt' для видалення альта\n'%role create RoleName' створення ролі для телеграм\n'%role delete Rolename' видалення ролі для телеграм\n'%role s RoleName' для підписки на роль\n'%role u RoleName' для видалення підписки на роль\n",
"CODE_FOR_CONNECT":"Код для підключення до сервера %s.",
"ERROR_SEND":"%s будь ласка, відправте мені команду старт в особистих повідомленнях, я як бот не можу перший відправити вам особисте повідомлення. І потім повторіть команду.",
"INSTRUCTIONS_SEND":"%s, Інструкцію надіслали вам у ЛП.",
"PLEASE_PASTE_CODE":"Будь ласка, вставте код у програму \n %s \n або просто перейдіть за посиланням для автоматичної авторизації \n %s",
"DATA_NOT_FOUND":"дані не знайдені",
"ALREADY_EXISTS":"вже існує",
"ALTO_ADDED":"альт %s додано",
"LIST_ALTS":"Список ваших альтов %+v",
"ALTO_REMOVED":"альт %s видалено",
"NO_ALTOS_FOUND":"альтів не знайдено",
"SCHEDULED_RETURNS":"%s, заплановані повернення БЗ",
"NO_SHIP_ARE_SCHEDULED":"%s, повернення кораблів не заплановано.",
"WILL_BE_ABLE_TO_RETURN":"%s %s зможе повернутися на Білу Зірку через 15 хвилин.",
"IS_NOW_ABLE_TO_RETURN":"%s %s тепер може повернутися до Білої зірки",
"IS_DUE_TO_RETURN":"%s %s %s має повернутися о %s (%s)",
"TIME_HAS_ALREADY_PASSED":"час вже пройшов",
"H_M_S":"%dг %dхв %dс",
"CODE_OUTDATED":"код застарів"
}}`)
}
