package dictionary

import (
	"encoding/json"
)

// temp function to integrate to exisiing logic
func (dict *Dictionary) setDictionaryRuJson() {

	dictRuJson := getDictionaryRuJson()

	var dictTemp map[string]map[string]string

	err := json.Unmarshal([]byte(dictRuJson), &dictTemp)
	if err != nil {
		dict.log.ErrorErr(err)
	}

	dict.ru = dictTemp["ru"]
}

func getDictionaryRuJson() string {
	return `{"ru":{
"you_in_queue":" ты уже в очереди",
"temp_queue_started":"%s запустил очередь %dlev КЗ",
"rs_queue":"Очередь КЗ",
"min":"мин.",
"forced_start":"принудительный старт",
"call_rs":"Сбор на КЗ",
"rs":"КЗ",
"you_subscribed_to_rs":"ты уже подписан на кз",
"to_add_to_queue_post":"для добавления в очередь напиши",
"you_subscribed_to_rs_ping":"вы подписались на пинг КЗ",
"you_not_subscribed_to_rs_ping":"ты не подписан на пинг КЗ",
"you_unsubscribed_from_rs_ping":"отписался от пинга КЗ",
"you_joined_queue":" присоединился к очереди",
"another_one_needed_to_complete_queue":"нужен еще один для фулки",
"queue_completed":"сформирована",
"go":"В ИГРУ",
"rs_queue_closed":" закрыл очередь КЗ",
"you_out_of_queue":" ты не в очереди",
"left_queue":" покинул очередь",
"was_deleted":"была удалена",
"info_set_emoji":" Для установки эмоджи пиши текст \nЭмоджи пробел (номер ячейки1-4) пробел эмоджи \n пример \nЭмоджи 1 🚀\n  Ваши слоты",
"your_emoji":"Ваши эмоджи\n",
"for_event":"для ивента",
"info_event_started":"Ивент запущен. После каждого похода на КЗ, один из участников КЗ вносит полученные очки в базу командой К (номер катки) (количество набраных очков)",
"event_mode_enabled":"Режим ивента уже активирован.",
"info_event_starting":"Запуск | Оcтановка Ивента доступен Администратору канала..",
"event_stopped":"Ивент остановлен.",
"info_event_not_active":"Ивент и так не активен. Нечего останавливать ",
"rs_data_entered":"данные о КЗ уже внесены ",
"points_added_to_database":"Очки внесены в базу",
"info_points_cannot_be_added":"добавление очков невозможно. Вы не являетесь участником КЗ под номером",
"event_not_started":"Ивент не запущен.",
"event_game":"Ивент игра",
"contributed":"Внесено",
"info_time_almost_up":" время почти вышло...\nДля продления времени ожидания на 30м жми +\nДля выхода из очереди жми -",
"info_cannot_click_plus":"рановато плюсик жмешь, ты в очереди на КЗ",
"info_cannot_click_minus":"рановато минус жмешь, ты в очереди на КЗ",
"budeshEshe":"будешь еще",
"timer_updated":" время обновлено ",
"empty":" пуста ",
"no_active_queues":"Нет активных очередей ",
"info_forced_start_available":"Принудительный старт доступен участникам очереди.",
"was_launched_incomplete":"была запущена не полной",
"info_max_queue_time":"максимальное время в очереди ограничено на 180 минут\n твое время",
"timer_updated":" время обновлено +30",
"scan_db":"Сканирую базу данных",
"no_history":" История не найдена ",
"form_list":"Формирую список ",
"top_participants":"ТОП Участников",
"event":"ивента",
"you_subscribed_to":"Теперь вы подписаны на",
"you_already_subscribed_to":"Вы уже подписаны на",
"error_rights_assign":"ошибка: недостаточно прав для выдачи роли ",
"error_rights_remove":"ошибка: недостаточно прав для снятия роли  ",
"you_not_subscribed_to_role":"Вы не подписаны на роль",
"role_not_exist":"нет такой роли на сервере",
"you_unsubscribed":"Вы отписались от роли",
"wishing_to":"Желающие на",
"to_add_to_queue":"для добавления в очередь",
"to_exit_the_queue":"для выхода из очереди",
"data_updated":"Данные обновлены",
"info_help_text":"Стать в очередь: [4-11]+  или\n [4-11]+[указать время ожидания в минутах]\n(уровень КЗ)+(время ожидания)\n  9+  встать в очередь на КЗ 9ур.\n  9+60  встать на КЗ 9ур, время ожидания не более 60 минут.\nПокинуть очередь: [4-11] -\n 9- выйти из очереди КЗ 9ур.\nПосмотреть список активных очередей: о[4-11]\n  о9 вывод очередь для вашей КЗ\n Получить роль КЗ: + [5-11]\n +9 получить роль КЗ 9ур.\n -9 снять роль \n Для Тёмных красных звезд\n Для старта очереди\n9*\nДля получения роли \n+d9",
"information":"Справка",
"info_bot_delete_msg":"ВНИМАНИЕ БОТ УДАЛЯЕТ СООБЩЕНИЯ \n ОТ ПОЛЬЗОВАТЕЛЕЙ ЧЕРЕЗ 3 МИНУТЫ",
"info_activation_not_required":"Я уже могу работать на вашем канале\nповторная активация не требуется.\nнапиши Справка",
"tranks_for_activation":"Спасибо за активацию.",
"channel_not_connected":"ваш канал и так не подключен к логике бота ",
"you_disabled_bot_functions":"вы отключили мои возможности",
"drs":"ТКЗ",
"queue_drs":"Очередь ТКЗ",
"drs_queue_closed":" закрыл очередь ТКЗ",
"language_switched_to":"Вы переключили меня на Русский язык"}}`
}
