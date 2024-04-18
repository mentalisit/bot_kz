package dictionary

import (
	"encoding/json"
)

// temp function to integrate to exisiing logic
func (dict *Dictionary) setDictionaryEnJson() {

	//dictEnJson := getDictionaryEnJson()

	//var dictTemp map[string]map[string]string

	//err := json.Unmarshal([]byte(dictEnJson), &dictTemp)
	err := json.Unmarshal([]byte(getDictionaryEnJson()), &dict.dictionary)
	if err != nil {
		dict.log.ErrorErr(err)
	}

	//dict.en = dictTemp["en"]
}

func getDictionaryEnJson() string {
	return `{"en":{
"you_in_queue":"You're already in the queue",
"temp_queue_started":"%s got in queue %slvl RS",
"rs_queue":"Queue RS",
"min":"min.",
"forced_start":"forced start",
"call_rs":"call a RS",
"rs":"RS",
"you_subscribed_to_rs":"You're already subscribed to RS",
"to_add_to_queue_post":"to add you to the queue, post",
"you_subscribed_to_rs_ping":"you've subscribed to RS ping",
"you_not_subscribed_to_rs_ping":"you are not subscribed to RS ping",
"you_unsubscribed_from_rs_ping":"You've unsubscribed from RS ping",
"you_joined_queue":"You've joined the queue",
"another_one_needed_to_complete_queue":"another one is needed to complete the queue",
"queue_completed":"queue completed",
"go":"Go!",
"rs_queue_closed":"RS queue closed",
"you_out_of_queue":"You're out of queue",
"left_queue":"left the queue",
"was_deleted":"was deleted",
"info_set_emoji":"To set an emoji, enter the text \n'Emoji_{cell number 1-4}_{your emoji}'\for example, \n'Emoji 1 🚀'\n Your Slots:",
"your_emoji":"Your emoji:\n",
"for_event":"for the event",
"info_event_started":"The event has started. After each RS, one of the RS participants contributes the points into the database with the command 'K {RS number} {number of points scored}'.",
"event_mode_enabled":"Event mode is already enabled.",
"info_event_started":"Event started | The event can be stopped by the channel administrator.",
"event_stopped":"Event stopped.",
"info_event_not_active":"The event is not active. There is nothing to stop",
"rs_data_entered":"RS data has already been entered",
"points_added_to_database":"Points have been added to the database",
"info_points_cannot_be_added":"Points cannot be added. You are not a member of RS number",
"event_not_started":"The event hasn't started.",
"event_game":"event game",
"contributed":"contributed",
"info_time_almost_up":" time is almost up...\nTo extend the waiting time by 30m, click '+'\nTo exit the queue, click '-'",
"info_cannot_click_plus":"too early to click the '+', you're in the RS queue",
"info_cannot_click_minus":"too early to click the '-', you're in the RS queue",
"you_will_still":"you will still",
"timer_updated":"timer updated",
"empty":" empty ",
"no_active_queues":"No active queues ",
"info_forced_start_available":"Forced start is available for queue members.",
"was_launched_incomplete":"was launched incomplete",
"info_max_queue_time":"the maximum time in the queue is limited to 180 minutes\n  your time",
"timer_updated":" timer updated +30m",
"scan_db":"Scanning the database",
"no_history":" History not found ",
"form_list":"Forming a list ",
"top_participants":"TOP Participants",
"event":"event",
"you_subscribed_to":"You are now subscribed to",
"you_already_subscribed_to":"You are already subscribed to",
"error_rights_assign":"error: Insufficient rights to assign a role ",
"error_rights_remove":"error: insufficient rights to remove a role ",
"you_not_subscribed_to_role":"You are not subscribed to the role",
"role_not_exist":"The role doesn't exist",
"you_unsubscribed":"You've unsubscribed from the role",
"wishing_to":"Wishing to",
"to_add_to_queue":"to add to the queue",
"to_exit_the_queue":"to exit the queue",
"data_updated":"Data updated",
"info_help_text":"Get in queue: [4-11]+ or\n  [4-11]+[specify timeout in minutes]\n(rs level)+(waiting time)\n  9+ stand in queue for short circuit 9level.\n  9 + 60 get on short circuit 9level, waiting time is not more than 60 minutes.\nLeave queue: [4-11] -\n  9- exit the queue RS 9level.\nView list of active queues: q[4-11]\n  q9 output queue for your rs\n Get Role Rs: + [5-11]\n  +9 get the role of RS lvl 9.\n  -9 remove the role\n For Dark Red Stars\nTo start the queue\n  9*\nTo get a role\n  +d9",
"information":"Information",
"info_bot_delete_msg":"WARNING\n BOT DELETES USER MESSAGES\n AFTER 3 MINUTES.",
"info_activation_not_required":"Bot already works on your channel\nre-activation is not required.\nwrite Help",
"tranks_for_activation":"Thanks for the activation.",
"channel_not_connected":"your channel is not connected to the bot ",
"you_disabled_bot_functions":"you've disabled the bot's functions",
"DRS":"DRS",
"queue_drs":"Queue DRS",
"drs_queue_closed":" DRS queue closed",
"language_switched_to":"You switched the bot to English",
"select_module_level":"Module selected: %s, level: %d",
"delete_module_level":"Removed module: %s, level: %d",
"install_weapon":"Weapon installed: %s"}}`
}
