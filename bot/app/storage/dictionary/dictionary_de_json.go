package dictionary

func getDictionaryDeJson() []byte {
	return []byte(`{"de":{
"you_in_queue":"You're already in the queue",
"temp_queue_started":"%s got in queue %s",
"rs_queue":"Queue RR",
"min":"min.",
"forced_start":"forced start",
"call_rs":"call a RR", 
"rs":"RR",
"you_subscribed_to_rs":"You're already subscribed to RR",
"to_add_to_queue_post":"to add you to the queue, post",
"you_subscribed_to_rs_ping":"you've subscribed to RR ping",
"you_not_subscribed_to_rs_ping":"you are not subscribed to RR ping",
"you_unsubscribed_from_rs_ping":"You've unsubscribed from RR ping",
"you_joined_queue":"You've joined the queue",
"another_one_needed_to_complete_queue":"Ein weiterer Spieler wird benötigt um den Gruppe zu vervollständigen.",
"queue_completed":"queue completed",
"go":"Start!",
"rs_queue_closed":"RS queue closed",
"you_out_of_queue":"You're out of queue",
"left_queue":"left the queue",
"was_deleted":"was deleted",
"info_set_emoji":"To set an emoji, enter the text \n'Emoji_{cell number 1-4}_{your emoji}'\nfor example, \n'Emoji 1 🚀'\n Your Slots:",
"your_emoji":"Your emoji:\n",
"for_event":"for the event",
"info_event_started":"The event has started. After each RR, one of the RR participants contributes the points into the database with the command 'K {RR number} {number of points scored}'.",
"event_mode_enabled":"Event mode is already enabled.",
"info_event_starting":"Event started | The event can be stopped by the channel administrator.",
"event_stopped": "Das Ereignis ist beendet.",
"info_event_not_active":"The event is not active. There is nothing to stop",
"rs_data_entered":"RR data has already been entered",
"points_added_to_database":"Points have been added to the database",
"info_points_cannot_be_added":"Points cannot be added. You are not a member of RR number",
"event_not_started": "Der Ereignismodus ist noch nicht aktiviert.",
"event_game":"event game",
"contributed":"contributed",
"info_time_almost_up":" time is almost up...\nTo extend the waiting time by 30m, click '+'\nTo exit the queue, click '-'",
"info_cannot_click_plus":"too early to click the '+', you're in the RR queue",
"info_cannot_click_minus":"too early to click the '-', you're in the RR queue",
"you_will_still":"you will still",
"timer_updated":"timer updated",
"empty":" empty ",
"no_active_queues":"No active queues ",
"info_forced_start_available":"Forced start is available for queue members.",
"was_launched_incomplete":"was launched incomplete", 
"info_max_queue_time":"the maximum time in the queue is limited to 180 minutes\n  your time",  
"scan_db":"Scanning the database",
"no_history":" History not found ",
"form_list":"Forming a list ",
"top_participants":"TOP Participants",
"event":"event",
"you_subscribed_to":"You are now subscribed to",
"you_already_subscribed_to":"You are already subscribed to", 
"error_rights_assign": "Fehler: Fehlende Berechtigung um die Rolle zuzuweisen:",
"error_rights_remove": "Fehler: Fehlende Berechtigung um die Rolle zu entfernen:",
"you_not_subscribed_to_role":"You are not subscribed to the role", 
"role_not_exist":"The role doesn't exist",
"you_unsubscribed":"You've unsubscribed from the role",
"wishing_to":"Wishing to",
"to_add_to_queue":"to add to the queue",
"to_exit_the_queue":"to exit the queue",
"data_updated":"Daten wurden aktualisiert.",
"info_help_text":"Get in queue: [4-11]+ or\n  [4-11]+[specify timeout in minutes]\n(rs level)+(waiting time)\n  9+ stand in queue for short circuit 9level.\n  9 + 60 get on short circuit 9level, waiting time is not more than 60 minutes.\nLeave queue: [4-11] -\n  9- exit the queue RS 9level.\nView list of active queues: q[4-11]\n  q9 output queue for your rs\n Get Role Rs: + [5-11]\n  +9 get the role of RS lvl 9.\n  -9 remove the role\n For Dark Red Stars\nTo start the queue\n  9*\nTo get a role\n  +d9",
"information":"Information",	
"info_bot_delete_msg":"WARNING\n BOT DELETES USER MESSAGES\n AFTER 3 MINUTES.",
"info_activation_not_required":"Bot already works on your channel\nre-activation is not required.\nwrite Help",
"tranks_for_activation":"Thanks for the activation.",
"channel_not_connected": "Dein Kanal ist nicht mit dem Bot verbunden",
"you_disabled_bot_functions":"you've disabled the bot's functions",
"drs":"DRR",
"queue_drs":"Queue DRR",
"drs_queue_closed":" DRR queue closed",
"language_switched_to":"You switched the bot to English",
"select_module_level":"Module selected: %s, level: %d",
"delete_module_level": "Modul %s, Level %d entfernt.",
"install_weapon":"Weapon installed: %s",
"temp1_queue":"Queue RS%s (%d)\n1️⃣ %s - %smin. (%d)\n\n%s++ - forced start"}}`)
}