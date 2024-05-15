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
"NO_ALTOS_FOUND":"no altos found"
}}`)
}
