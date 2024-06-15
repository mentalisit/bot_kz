package dictionary

func getDictionaryRuJson() []byte {
	return []byte(`{"ru":{
"HELP_TEXT_DS": "на текуший момент доступны команды:\n '**%помощь**' для получения текущей справки \n'**%подключить**' для подключения приложения\n '**%т и**' для получения изображения с вашими модулями\n '**%т @имя и**' для получения изображения с модулями другого игрока\n'**%т имя и**' для получения изображения с модулями альта\n '**%alts add NameAlt**' для создания альта для технологий\n'**%alts del NameAlt**' для удаления альта\n'**%ник игровоеИмя**' установка игрового имени \n'**%time**' отображение текущего времени \n",
"HELP_TEXT_TG": "на текуший момент доступны команды:\n '%помощь' для получения текущей справки \n'%подключить' для подключения приложения\n '%т и' для получения изображения с вашими модулями\n '%т @имя и' для получения изображения с модулями другого игрока\n'%т имя и' для получения изображения с модулями альта\n '%alts add NameAlt' для создания альта для технологий\n'%alts del NameAlt' для удаления альта\n'%role create RoleName' создание роли для телеграм\n'%role delete Rolename' удаление роли для телеграм\n'%role s RoleName' для подписки на роль\n'%role u RoleName' для удаления подписки на роль\n'%ник игровоеИмя' установка игрового имени \n'%time' отображение текущего времени \n",
"CODE_FOR_CONNECT":"Код для подключения приложения к серверу %s.",
"ERROR_SEND":"%s пожалуйста отправьте мне команду старт в личных сообщениях, я как бот не могу первый отправить вам личное сообщение. И после повторите команду.",
"INSTRUCTIONS_SEND":"%s, Инструкцию отправили вам в ЛС.",
"PLEASE_PASTE_CODE":"Пожалуйста, вставьте код в приложение \n %s \n или просто перейдите по ссылке для автоматической авторизации \n %s",
"DATA_NOT_FOUND":"данные не найдены",
"ALREADY_EXISTS":"уже существует",
"ALTO_ADDED":"альт %s добавлен",
"LIST_ALTS":"Список ваших альтов %+v",
"ALTO_REMOVED":"альт %s удален",
"NO_ALTOS_FOUND":"альты не найдены",
"SCHEDULED_RETURNS":"%s, запланированные возвраты БЗ",
"NO_SHIP_ARE_SCHEDULED":"%s, возвращения кораблей не запланировано.",
"WILL_BE_ABLE_TO_RETURN":"%s %s сможет вернуться на Белую Звезду через 15 минут.",
"IS_NOW_ABLE_TO_RETURN":"%s %s теперь может вернуться на Белую Звезду.",
"IS_DUE_TO_RETURN":"%s %s %s должен вернуться в %s (%s)",
"TIME_HAS_ALREADY_PASSED":"время уже прошло",
"H_M_S":"%dч %dм %dс",
"CODE_OUTDATED":"код устарел",
"I_COULD_NOT_FIND_ANY":"%s, мне не удалось найти часовые пояса, соответствующие '%s'",
"TIMEZONA_SET":"%s, часовой пояс для %s установлен на %s",
"TIMEZONA_IS_CURRENTLY":"%s, часовой пояс для %s в настоящее время установлен на '%s'",
"LOCAL_TIME_FOR_EVERYONE":"%s Местное время для всех:",
"UNLISTED_MEMBERS":"У участников, не включенных в список, нет настройки часового пояса. Чтобы установить его, они могут использовать команду %tz set +3.",
"YOU_ARE_NOT_FOUND":"%s вас нет в базе данных, отправьте %%connect",
"GAME_NAME_SET":"%s, игровое имя установлено на '%s'",
"HELP_NICKNAME":"%s, Команда %%ник используется для установки имени,\n' %%ник имя ', если имя не содержит пробелов \nили \n' %%ник \"мое имя\" ', если имя состоит из нескольких слов \n nпример\n' %%ник Вася '\n' %%ник \"Вася Иванов\" '",
"SECRET_LINK":"Вот ваша постоянная [секретная ссылка](%s) для сервера %s, не передавайте её никому."
}}`)
}
