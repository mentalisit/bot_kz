package dictionary

func (dict *Dictionary) setDictionaryRu() {

	var dictLg = make(map[string]string)

	dictLg["1"] = "один"
	dictLg["2"] = "два"
	dictLg["you_in_queue"] = " ты уже в очереди"
	dictLg["temp_queue_started"] = "%s запустил очередь %dlev КЗ"
	dictLg["rs_queue"] = "Очередь кз"
	dictLg["min."] = "мин."
	dictLg["prinuditelniStart"] = "принудительный старт"
	dictLg["SborNaKz"] = "Сбор на кз"
	dictLg["kz"] = "кз"
	dictLg["tiUjePodpisanNaKz"] = "ты уже подписан на кз"
	dictLg["dlyaDobavleniyaVochered"] = "для добавления в очередь напиши"
	dictLg["viPodpisalisNaPing"] = "вы подписались на пинг кз"
	dictLg["tiNePodpisanNaPingKz"] = "ты не подписан на пинг кз"
	dictLg["otpisalsyaOtPingaKz"] = "отписался от пинга кз"
	dictLg["prisoedenilsyKocheredi"] = " присоединился к очереди"
	dictLg["nujenEsheOdinDlyFulki"] = "нужен еще один для фулки"
	dictLg["sformirovana"] = "сформирована"
	dictLg["Vigru"] = "В ИГРУ"
	dictLg["zakrilOcheredKz"] = " закрыл очередь кз"
	dictLg["tiNeVOcheredi"] = " ты не в очереди"
	dictLg["pokinulOchered"] = " покинул очередь"
	dictLg["bilaUdalena"] = "была удалена"
	dictLg["dly ustanovki"] = "	Для установки эмоджи пиши текст \nЭмоджи пробел (номер ячейки1-4) пробел эмоджи \n пример \nЭмоджи 1 🚀\n	Ваши слоты"
	dictLg["vashiEmodji"] = "Ваши эмоджи\n"
	dictLg["dly iventa"] = "для ивента"
	dictLg["iventZapushen"] = "Ивент запущен. После каждого похода на КЗ, один из участников КЗ вносит полученные очки в базу командой К (номер катки) (количество набраных очков)"
	dictLg["rejimIventaUje"] = "Режим ивента уже активирован."
	dictLg["zapuskIostanovka"] = "Запуск | Оcтановка Ивента доступен Администратору канала."
	dictLg["IventOstanovlen"] = "Ивент остановлен."
	dictLg["iventItakAktiven"] = "Ивент и так не активен. Нечего останавливать "
	dictLg["dannieKzUjeVneseni"] = "данные о кз уже внесены "
	dictLg["ochki vnesen"] = "Очки внесены в базу"
	dictLg["dobavlenieOchkovNevozmojno"] = "добавление очков невозможно. Вы не являетесь участником КЗ под номером"
	dictLg["iventNeZapushen"] = "Ивент не запущен."
	dictLg["iventIgra"] = "Ивент игра"
	dictLg["vneseno"] = "Внесено"
	dictLg["VremyaPochtiVishlo"] = " время почти вышло...\nДля продления времени ожидания на 30м жми +\nДля выхода из очереди жми -"
	dictLg["ranovatoPlysik"] = "рановато плюсик жмешь, ты в очереди на кз"
	dictLg["budeshEshe"] = "будешь еще"
	dictLg["vremyaObnovleno"] = " время обновлено "
	dictLg["ranovatoMinus"] = "рановато минус жмешь, ты в очереди на кз"
	dictLg["pusta"] = " пуста " //очередь кз пуста
	dictLg["netAktivnuh"] = "Нет активных очередей "
	dictLg["prinuditelniStartDostupen"] = "Принудительный старт доступен участникам очереди."
	dictLg["bilaZapushenaNe"] = "была запущена не полной"
	dictLg["maksimalnoeVremya"] = "максимальное время в очереди ограничено на 180 минут\n твое время"
	dictLg["vremyaObnovleno"] = " время обновлено +30"
	dictLg["ScanDB"] = "Сканирую базу данных"
	dictLg["noHistory"] = " История не найдена "
	dictLg["formlist"] = "Формирую список "
	dictLg["topUchastnikov"] = "ТОП Участников"
	dictLg["iventa"] = "ивента"
	dictLg["teperViPodpisani"] = "Теперь вы подписаны на"
	dictLg["ViUjePodpisan"] = "Вы уже подписаны на"
	dictLg["oshibkaNedostatochno"] = "ошибка: недостаточно прав для выдачи роли "
	dictLg["viNePodpisani"] = "Вы не подписаны на роль"
	dictLg["netTakoiRoli"] = "нет такой роли на сервере"
	dictLg["ViOtpisalis"] = "Вы отписались от роли"
	dictLg["OshibkaNedostatochnadlyaS"] = "ошибка: недостаточно прав для снятия роли  "
	dictLg["jelaushieNa"] = "Желающие на"
	dictLg["DlyaDobavleniya"] = "для добавления в очередь"
	dictLg["DlyaVihodaIz"] = "для выхода из очереди"
	dictLg["DannieObnovleni"] = "Данные обновлены"
	dictLg["hhelpText"] = "Стать в очередь: [4-11]+  или\n " +
		"[4-11]+[указать время ожидания в минутах]\n" +
		"(уровень кз)+(время ожидания)\n" +
		" 9+  встать в очередь на КЗ 9ур.\n" +
		" 9+60  встать на КЗ 9ур, время ожидания не более 60 минут.\n" +
		"Покинуть очередь: [4-11] -\n 9- выйти из очереди КЗ 9ур.\n" +
		"Посмотреть список активных очередей: о[4-11]\n" +
		" о9 вывод очередь для вашей Кз\n" +
		"Получить роль кз: + [5-11]\n +9 получить роль КЗ 9ур.\n -9 снять роль \n" +
		"Для Тёмных красных звезд\n Для старта очереди\n9*\nДля получения роли \n+d9"
	dictLg["spravka"] = "Справка"
	dictLg["botUdalyaet"] = "ВНИМАНИЕ БОТ УДАЛЯЕТ СООБЩЕНИЯ \n ОТ ПОЛЬЗОВАТЕЛЕЙ ЧЕРЕЗ 3 МИНУТЫ"
	dictLg["accessAlready"] = "Я уже могу работать на вашем канале\nповторная активация не требуется.\nнапиши Справка"
	dictLg["accessTY"] = "Спасибо за активацию."
	dictLg["accessYourChannel"] = "ваш канал и так не подключен к логике бота "
	dictLg["YouDisabledMyFeatures"] = "вы отключили мои возможности"
	dictLg["dkz"] = "ткз"
	dictLg["ocheredTKz"] = "Очередь ткз"
	dictLg["zakrilOcheredTKz"] = " закрыл очередь ткз"
	dictLg["vashLanguage"] = "Вы переключили меня на Русский язык"

	dict.ru = dictLg
}
