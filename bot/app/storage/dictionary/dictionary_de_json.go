package dictionary

func getDictionaryDeJson() []byte {
	return []byte(`{"de":{
"you_in_queue":"Sie befinden sich bereits in der Warteschlange",
  "temp_queue_started":"%s ist in der Warteschlange %s angekommen",
  "rs_queue":"Warteschlange RR",
  "min":"Mindest.",
  "forced_start":"Zwangsstart",
  "call_rs":"Rufen Sie einen RR an",
  "rs":"RR",
  "you_subscribed_to_rs":"Sie haben RR bereits abonniert",
  "to_add_to_queue_post":"Um Sie zur Warteschlange hinzuzufügen, posten Sie",
  "you_subscribed_to_rs_ping":"Sie haben RR-Ping abonniert",
  "you_not_subscribed_to_rs_ping":"Sie haben RR-Ping nicht abonniert",
  "you_unsubscribed_from_rs_ping":"Sie haben sich vom RR-Ping abgemeldet",
  "you_joined_queue":"Sie haben sich der Warteschlange angeschlossen",
  "another_one_needed_to_complete_queue":"Ein weiterer Spieler wird benötigt um den Gruppe zu vervollständigen.",
  "queue_completed":"Warteschlange abgeschlossen",
  "go":"Start!",
  "rs_queue_closed":"RR-Warteschlange geschlossen",
  "you_out_of_queue":"Sie befinden sich außerhalb der Warteschlange",
  "left_queue":"verließ die Warteschlange",
  "was_deleted":"wurde gelöscht",
  "info_set_emoji":"Um ein Emoji festzulegen, geben Sie den Text ein\n„Emoji_{Zellennummer 1-4}_{Ihr Emoji}“ oder Beispiel,\n„Emoji 1 \uD83D\uDE80“\n Ihre Slots:",
  "your_emoji":"Dein Emoji:\n",
  "for_event":"Für das Event",
  "info_event_started":"The event has started. After each RR, one of the RR participants contributes the points into the database with the command 'K {RR number} {number of points scored}'.",
  "event_mode_enabled":"Der Ereignismodus ist bereits aktiviert.",
  "info_event_starting":"Veranstaltung gestartet | Das Ereignis kann vom Kanaladministrator gestoppt werden.",
  "event_stopped": "Das Ereignis ist beendet.",
  "info_event_not_active":"Das Ereignis ist nicht aktiv. Es gibt nichts zu stoppen",
  "rs_data_entered":"RR-Daten wurden bereits eingegeben",
  "points_added_to_database":"Der Datenbank wurden Punkte hinzugefügt",
  "info_points_cannot_be_added":"Es können keine Punkte hinzugefügt werden. Sie sind kein Mitglied der RR-Nummer",
  "event_not_started": "Der Ereignismodus ist noch nicht aktiviert.",
  "event_game":"Event-Spiel",
  "contributed":"beigetragen",
  "info_time_almost_up":" Die Zeit ist fast abgelaufen...\nUm die Wartezeit um 30 Minuten zu verlängern, klicken Sie auf „+“.\nUm die Warteschlange zu verlassen, klicken Sie auf „-“",
  "info_cannot_click_plus":"Wenn es zu früh ist, um auf das „+“ zu klicken, befinden Sie sich in der RR-Warteschlange",
  "info_cannot_click_minus":"Wenn es zu früh ist, auf das „-“ zu klicken, befinden Sie sich in der RR-Warteschlange",
  "you_will_still":"Du wirst es immer noch tun",
  "timer_updated":"Timer aktualisiert",
  "empty":" leer ",
  "no_active_queues":"Keine aktiven Warteschlangen ",
  "info_forced_start_available":"Für Warteschlangenmitglieder ist ein erzwungener Start verfügbar.",
  "was_launched_incomplete":"wurde unvollständig gestartet",
  "info_max_queue_time":"Die maximale Zeit in der Warteschlange ist auf 180 Minuten begrenzt\n  deine Zeit",
  "scan_db":"Scannen der Datenbank",
  "no_history":" Verlauf nicht gefunden ",
  "form_list":"Eine Liste erstellen ",
  "top_participants":"TOP-Teilnehmer",
  "event":"Ereignis",
  "you_subscribed_to":"Sie sind jetzt abonniert",
  "you_already_subscribed_to":"Sie sind bereits abonniert",
  "error_rights_assign": "Fehler: Fehlende Berechtigung um die Rolle zuzuweisen:",
  "error_rights_remove": "Fehler: Fehlende Berechtigung um die Rolle zu entfernen:",
  "you_not_subscribed_to_role":"Sie haben die Rolle nicht abonniert",
  "role_not_exist":"Die Rolle existiert nicht",
  "you_unsubscribed":"Sie haben sich von der Rolle abgemeldet",
  "wishing_to":"Ich wünsche es",
  "to_add_to_queue":"zur Warteschlange hinzuzufügen",
  "to_exit_the_queue":"um die Warteschlange zu verlassen",
  "data_updated":"Daten wurden aktualisiert.",
  "info_help_text":"In die Warteschlange kommen: [4-11]+ oder\n  [4-11]+[Timeout in Minuten angeben]\n(RS-Level)+(Wartezeit)\n  9+ stehen in der Warteschlange für Kurzschluss 9Level.\n  9 + 60 steigen Sie in die Kurzschlussstufe 9 ein, die Wartezeit beträgt nicht mehr als 60 Minuten.\nWarteschlange verlassen: [4-11] -\n  9- Verlassen Sie die Warteschlange RS 9Ebene.\nListe der aktiven Warteschlangen anzeigen: q[4-11]\n  q9-Ausgabewarteschlange für Ihre RS\n Holen Sie sich Rollen-Rs: + [5-11]\n  +9 erhält die Rolle des RS Level 9.\n  -9 Entfernen Sie die Rolle\n Für dunkelrote Sterne\nUm die Warteschlange zu starten\n  9*\nUm eine Rolle zu bekommen\n  +d9",
  "information":"Information",
  "info_bot_delete_msg":"WARNUNG\n BOT löscht Benutzernachrichten\n NACH 3 MINUTEN.",
  "info_activation_not_required":"Der Bot funktioniert bereits auf deinem Kanal\nEine erneute Aktivierung ist nicht erforderlich.\nHilfe schreiben(help)",
  "tranks_for_activation":"Danke für die Aktivierung.",
  "channel_not_connected": "Dein Kanal ist nicht mit dem Bot verbunden",
  "you_disabled_bot_functions":"Sie haben die Funktionen des Bots deaktiviert",
  "drs":"DRR",
  "queue_drs":"Warteschlange DRR",
  "drs_queue_closed":" DRR-Warteschlange geschlossen",
  "language_switched_to":"Sie haben den Bot auf Deutsch umgestellt",
  "select_module_level":"Module selected: %s, level: %d",
  "delete_module_level": "Modul %s, Level %d entfernt.",
  "install_weapon":"Installierte Waffe: %s",
  "temp1_queue":"Warteschlange RR%s (%d)\n1\uFE0F⃣ %s - %smin. (%D)\n\n%s++ – erzwungener Start"}}`)
}
