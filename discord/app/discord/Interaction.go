package DiscordClient

import (
	"discord/discord/helpers"
	"discord/models"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log/slog"
	"strconv"
	"time"
)

// slash command module respond
func (d *Discord) handleModuleCommand(i *discordgo.InteractionCreate, locale string) {
	module := i.ApplicationCommandData().Options[0].StringValue()
	level := i.ApplicationCommandData().Options[1].IntValue()

	response := fmt.Sprintf(getText(locale, "select_module_level"), module, level)
	if level == 0 {
		response = fmt.Sprintf(getText(locale, "delete_module_level"), module, level)
	}

	user := i.Interaction.Member.User
	use := existCompendiumData(user.Username, user.ID, i.Interaction.GuildID)
	if use {
		response = getText(locale, "USE_COMPENDIUM")
	}
	// Отправка ответа
	err := d.S.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
	if err != nil {
		d.log.ErrorErr(err)
		return
	}
	go func() {
		time.Sleep(20 * time.Second)
		err = d.S.InteractionResponseDelete(i.Interaction)
		if err != nil {
			return
		}
	}()

	if use {
		return
	}

	d.updateModuleOrWeapon(i.Interaction.Member.User.Username, module, strconv.FormatInt(level, 10))
}

// slash command weapon respond
func (d *Discord) handleWeaponCommand(i *discordgo.InteractionCreate, locale string) {
	weapon := i.ApplicationCommandData().Options[0].StringValue()

	response := fmt.Sprintf(getText(locale, "install_weapon"), weapon)

	// Отправка ответа
	err := d.S.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
	if err != nil {
		d.log.ErrorErr(err)
		return
	}
	go func() {
		time.Sleep(20 * time.Second)
		err = d.S.InteractionResponseDelete(i.Interaction)
		if err != nil {
			return
		}
	}()
	d.updateModuleOrWeapon(i.Interaction.Member.User.Username, weapon, "")
}
func (d *Discord) updateModuleOrWeapon(username, module, level string) {
	rse := "<:rse:1199068829511335946> " + level
	genesis := "<:genesis:1199068748280242237> " + level
	enrich := "<:enrich:1199068793633251338> " + level
	if level == "0" {
		rse, genesis, enrich = "", "", ""
	}

	barrage := "<:barrage:1199084425393225782>"
	laser := "<:laser:1199084197571207339>"
	chainray := "<:chainray:1199073579577376888>"
	battery := "<:batteryw:1199072534562345021>"
	massbattery := "<:massbattery:1199072493760151593>"
	dartlauncher := "<:dartlauncher:1199072434674991145>"
	rocketlauncher := "<:rocketlauncher:1199071677548605562>"
	t, err := d.storage.Emoji.EmojiModuleReadUsers(username, "ds")
	if err != nil {
		slog.Error(err.Error())
	}
	if len(t.Name) == 0 {
		d.storage.Emoji.EmInsertEmpty("ds", username)
	}
	switch module {
	case "RSE":
		d.storage.Emoji.ModuleUpdate(username, "ds", "1", rse)
	case "GENESIS":
		d.storage.Emoji.ModuleUpdate(username, "ds", "2", genesis)
	case "ENRICH":
		d.storage.Emoji.ModuleUpdate(username, "ds", "3", enrich)
	case "barrage":
		d.storage.Emoji.WeaponUpdate(username, "ds", barrage)
	case "laser":
		d.storage.Emoji.WeaponUpdate(username, "ds", laser)
	case "chainray":
		d.storage.Emoji.WeaponUpdate(username, "ds", chainray)
	case "battery":
		d.storage.Emoji.WeaponUpdate(username, "ds", battery)
	case "massbattery":
		d.storage.Emoji.WeaponUpdate(username, "ds", massbattery)
	case "dartlauncher":
		d.storage.Emoji.WeaponUpdate(username, "ds", dartlauncher)
	case "rocketlauncher":
		d.storage.Emoji.WeaponUpdate(username, "ds", rocketlauncher)
	case "Remove":
		d.storage.Emoji.WeaponUpdate(username, "ds", "")
	}
}
func (d *Discord) handleButtonPressed(i *discordgo.InteractionCreate) {
	ok, config := d.CheckChannelConfigDS(i.ChannelID)
	if ok {
		user := i.Interaction.User
		if i.Interaction.Member != nil && i.Interaction.Member.User != nil {
			user = i.Interaction.Member.User
		}

		in := models.InMessage{
			Mtext:       i.MessageComponentData().CustomID,
			Tip:         "ds",
			Username:    user.Username,
			UserId:      user.ID,
			NameNick:    "",
			NameMention: user.Mention(),
			Ds: struct {
				Mesid   string
				Guildid string
				Avatar  string
			}{
				Mesid:   i.Interaction.Message.ID,
				Guildid: i.Interaction.GuildID,
				Avatar:  user.AvatarURL("128"),
			},
			Config: config,
			Option: models.Option{Reaction: true},
		}
		if i.Interaction.Member != nil && i.Interaction.Member.Nick != "" {
			in.NameNick = i.Interaction.Member.Nick
		}

		d.api.SendRsBotAppRecover(in)
		err := d.S.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredMessageUpdate,
		})
		if err != nil {
			d.log.ErrorErr(err)
			return
		}
	}
	ds, bridgeConfig := d.BridgeCheckChannelConfigDS(i.ChannelID)
	if ds {
		user := i.Interaction.User
		if i.Interaction.Member != nil && i.Interaction.Member.User != nil {
			user = i.Interaction.Member.User
		}
		in := models.ToBridgeMessage{
			Text:          ".poll " + i.MessageComponentData().CustomID,
			Sender:        user.Username,
			Tip:           "ds",
			ChatId:        i.Interaction.ChannelID,
			MesId:         i.Interaction.Message.ID,
			GuildId:       i.Interaction.GuildID,
			TimestampUnix: i.Interaction.Message.Timestamp.Unix(),
			Config:        &bridgeConfig,
		}
		err := d.S.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredMessageUpdate,
		})
		if err != nil {
			d.log.ErrorErr(err)
			return
		}
		d.api.SendBridgeAppRecover(in)
	}
}
func existCompendiumData(name, userid, guildid string) bool {
	genesis1, enrich1, rsextender1 := helpers.GetTechDataUserId(userid, guildid)
	genesis2, enrich2, rsextender2 := helpers.Get2TechDataUserId(name, userid, guildid)

	genesis := max(genesis1, genesis2)
	enrich := max(enrich1, enrich2)
	rsextender := max(rsextender1, rsextender2)
	if genesis == 0 && enrich == 0 && rsextender == 0 {
		return false
	}
	return true
}
func getText(lang, key string) string {
	t := make(map[string]map[string]string)
	t["en"] = make(map[string]string)
	t["ru"] = make(map[string]string)
	t["ua"] = make(map[string]string)

	t["ru"]["install_weapon"] = "Установлено оружие: %s"
	t["ua"]["install_weapon"] = "Встановлено зброю: %s"
	t["en"]["install_weapon"] = "Weapon installed: %s"

	t["ru"]["USE_COMPENDIUM"] = "используйте компендиум для установки модулей"
	t["ua"]["USE_COMPENDIUM"] = "використовуйте компендіум для встановлення модулів"
	t["en"]["USE_COMPENDIUM"] = "use compendium to install modules"

	t["ru"]["delete_module_level"] = "Удален модуль: %s, уровень: %d"
	t["ua"]["delete_module_level"] = "Видалений модуль: %s, рівень: %d"
	t["en"]["delete_module_level"] = "Removed module: %s, level: %d"

	t["ru"]["select_module_level"] = "Выбран модуль: %s, уровень: %d"
	t["ua"]["select_module_level"] = "Вибрано модуль: %s, рівень: %d"
	t["en"]["select_module_level"] = "Module selected: %s, level: %d"

	text := t[lang][key]
	if text == "" {
		return key
	}
	return text
}
