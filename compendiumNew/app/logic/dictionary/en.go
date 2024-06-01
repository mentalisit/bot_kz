package dictionary

func getDictionaryEnJson() []byte {
	return []byte(`{"en":{
"HELP_TEXT_DS": "Currently available commands are:\n '**%help**' to get current help \n'**%connect**' to connect the application\n '**%t i**' to get an image with your modules \n '**%t @name i**' to get an image with another player's modules\n'**%t name i**' to get an image with alt modules\n '**%alts add NameAlt**' to create alt for technologies\n'**%alts del NameAlt**' to delete alt\n",
"HELP_TEXT_TG": "Currently available commands are:\n '%help' to get current help \n'%connect' to connect the application\n '%t i' to get an image with your modules \n '%t @name i' to get an image with another player's modules\n'%t name i' to get an image with alt modules\n '%alts add NameAlt' to create alt for technologies\n'%alts del NameAlt' to delete alt\n'%role create RoleName' creating a role for telegrams\n'%role delete Rolename' deleting a role for telegrams\n'%role s RoleName' for subscribing to a role\n'%role u RoleName' to delete a role subscription\n",
"CODE_FOR_CONNECT":"Code for connecting the application to the server %s.",
"ERROR_SEND":"%s please send me the start command in private messages, as a bot I cannot be the first to send you a private message. And then repeat the command.",
"INSTRUCTIONS_SEND":"%s, Instructions have been sent to you via DM.",
"PLEASE_PASTE_CODE":"Please paste the code into the application \n %s \n or simply follow the link for automatic authorization \n %s",
"DATA_NOT_FOUND":"data not found",
"ALREADY_EXISTS":"already exists",
"ALTO_ADDED":"alto added %s",
"LIST_ALTS":"List of your alts %+v",
"ALTO_REMOVED":"alto removed %s",
"NO_ALTOS_FOUND":"no altos found",
"SCHEDULED_RETURNS":"%s, Scheduled WS Returns",
"NO_SHIP_ARE_SCHEDULED":"%s, No ships are scheduled to return.",
"WILL_BE_ABLE_TO_RETURN":"%s's %s will be able to return to the White Star in 15 minutes.",
"IS_NOW_ABLE_TO_RETURN":"%s's %s is now able to return to the White Star",
"IS_DUE_TO_RETURN":"%s %s's %s is due to return at %s (%s)",
"TIME_HAS_ALREADY_PASSED":"time has already passed",
"H_M_S":"%dh %dm %ds",
"CODE_OUTDATED":"the code is outdated",
"I_COULD_NOT_FIND_ANY":"%s, I could not find any timezones matching '%s'",
"TIMEZONA_SET":"%s,Timezona for %s set to %s",
"TIMEZONA_IS_CURRENTLY":"%s, Timezona for %s is currently set to '%s'",
"LOCAL_TIME_FOR_EVERYONE":"%s Local time for everyone:",
"UNLISTED_MEMBERS":"Unlisted members have no timezone setting. They can use the %tz set -5 command to set it."
}}`)
}
