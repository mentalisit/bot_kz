package dictionary

import (
	"encoding/json"
)

// temp function to integrate to exisiing logic
func (dict *Dictionary) setDictionaryUaJson() {

	//dictUaJson := getDictionaryUaJson()

	//var dictTemp map[string]map[string]string

	//err := json.Unmarshal([]byte(dictUaJson), &dictTemp)
	err := json.Unmarshal([]byte(getDictionaryUaJson()), &dict.dictionary)
	if err != nil {
		dict.log.ErrorErr(err)
	}

	//dict.ua = dictTemp["ua"]
}

func getDictionaryUaJson() string {
	return `{"ua":{
"you_in_queue":" ти вже у черзi",
"temp_queue_started":"%s запустив чергу %slvl КЗ",
"rs_queue":"Черга чз",
"min":"хв.",
"forced_start":"примусовий старт",
"call_rs":"Збір на чз",
"rs":"чз",
"you_subscribed_to_rs":"ти вже підписаний на чз",
"to_add_to_queue_post":"ви підписалися на пінг чз",
"you_subscribed_to_rs_ping":"для додавання до черги напиши",
"you_not_subscribed_to_rs_ping":"ти не підписаний на пінг чз",
"you_unsubscribed_from_rs_ping":"відписався від пінгу чз",
"you_joined_queue":" приєднався до черги",
"another_one_needed_to_complete_queue":"потрібен ще один для повної черги",
"queue_completed":"сформована",
"go":"В ГРУ",
"rs_queue_closed":" закрив чергу чз",
"you_out_of_queue":" ти не в черзі",
"left_queue":" залишив чергу",
"was_deleted":"була видалена",
"info_set_emoji":" Для встановлення емоджі пиши текст \nЕмоджі пробіл (номер ячейки1-4) пробіл емоджі \n приклад \nЕмоджі 1 🚀\n Ваші слоти",
"your_emoji":"Ваші емоджі\n",
"for_event":"для івента",
"info_event_started":"Івента запущено. Після кожного походу на ЧЗ, один із учасників ЧЗ вносить отримані очки в базу командою К (номер катки) (кількість набраних очок)",
"event_mode_enabled":"Режим івента вже активовано.",
"info_event_starting":"Запуск Зупинка Івента доступна Адміністратору каналу.",
"event_stopped":"Івента зупинено.",
"info_event_not_active":"Івент і так не є активним. Нема чого зупиняти",
"rs_data_entered":"дані про чз вже внесено",
"points_added_to_database":"Бали внесені до бази",
"info_points_cannot_be_added":"додавання балів неможливе. Ви не є учасником КЗ під номером",
"event_not_started":"Івента не запущено.",
"event_game":"івент гра",
"contributed":"внесено",
"info_time_almost_up":" час майже вийшов...\nДля продовження часу очікування на 30хв. тисни +\nДля виходу з черги тисни -",
"info_cannot_click_plus":"рано плюсик тиснеш, ти в черзі на чз",
"info_cannot_click_minus":"рано мінус тиснеш, ти в черзі на чз",
"you_will_still":"будешь еще",
"timer_updated":" час оновлено ",
"empty":" порожня ",
"no_active_queues":"Немає активних черг ",
"info_forced_start_available":"Примусовий старт доступний учасникам черги.",
"was_launched_incomplete":"була запущена не повною",
"info_max_queue_time":"максимальний час у черзі обмежений на 180 хвилин\n  твій час",
"timer_updated":" час оновлено +30хв",
"scan_db":"Сканую базу даних",
"no_history":"Історія не знайдена",
"form_list":"Формую список ",
"top_participants":"ТОП учасників",
"event":"івента",
"you_subscribed_to":"Тепер ви підписані на",
"you_already_subscribed_to":"Ви вже підписані на",
"error_rights_assign":"помилка: недостатньо прав для видачі ролі ",
"error_rights_remove":"помилка: недостатньо прав для зняття ролі ",
"you_not_subscribed_to_role":"Ви не підписані на роль",
"role_not_exist":"немає такої ролі на сервері",
"you_unsubscribed":"Ви відмовилися від ролі",
"wishing_to":"Бажаючі на",
"to_add_to_queue":"для додавання до черги",
"to_exit_the_queue":"для виходу з черги",
"data_updated":"Дані оновлені",
"info_help_text":"Стати у чергу: [4-11]+ або\n  [4-11]+[вказати час очікування у хвилинах]\n (Рівень чз) + (час очікування)\n  9+ стати в чергу на ЧЗ 9рівня.\n  9+60 стати на ЧЗ 9рівня, час очікування не більше 60 хвилин.\n Залишити чергу: [4-11] -\n  9 - вийти з черги ЧЗ 9рівня.\nПереглянути список активних черг: ч[4-11]\n ч9 демонстрація черги для вашої Чз\nОтримати роль чз: + [5-11]\n  +9 Отримати роль ЧЗ 9рівня.\n  -9 зняти роль\n Для темних червоних зірок\n  Для старту черги\n9*\nДля отримання ролі\n+d9",
"information":"Довідка",
"info_bot_delete_msg":"УВАГА БОТ ВИДАЛЯЄ ПОВІДОМЛЕННЯ\n  ВІД КОРИСТУВАЧІВ ЧЕРЕЗ 3 ХВИЛИНИ",
"info_activation_not_required":"Я вже можу працювати на вашому каналі\nповторна активація не потрібна.\nнапиши Довідка",
"tranks_for_activation":"Дякую за активацію.",
"channel_not_connected":"ваш канал і так не підключений до логіки бота ",
"you_disabled_bot_functions":"ви відключили мої можливості",
"drs":"тчз",
"queue_drs":"Черга тчз",
"drs_queue_closed":"  закрив чергу тчз",
"language_switched_to":"Ви переключили мене на українську мову",
"select_module_level":"Вибрано модуль: %s, рівень: %d",
"delete_module_level":"Видалений модуль: %s, рівень: %d",
"install_weapon":"Встановлено зброю: %s",
"temp1_queue":"Черга чз%s (%d)\n1️⃣ %s - %sхв. (%d)\n\n%s++ - примусовий старт"}}`
}
