package DiscordClient

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"kz_bot/bot/helpers"
	"kz_bot/models"
	"strconv"
	"time"
)

// slash command module respond
func (d *Discord) handleModuleCommand(i *discordgo.InteractionCreate, locale string) {
	module := i.ApplicationCommandData().Options[0].StringValue()
	level := i.ApplicationCommandData().Options[1].IntValue()

	response := fmt.Sprintf(d.getLanguage(locale, "select_module_level"), module, level)
	if level == 0 {
		response = fmt.Sprintf(d.getLanguage(locale, "delete_module_level"), module, level)
	}

	user := i.Interaction.Member.User
	use := existCompendiumData(user.Username, user.ID, i.Interaction.GuildID)
	if use {
		response = d.getLanguage(locale, "USE_COMPENDIUM")
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

	response := fmt.Sprintf(d.getLanguage(locale, "install_weapon"), weapon)

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
	t := d.storage.Emoji.EmojiModuleReadUsers(username, "ds")
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

		d.ChanRsMessage <- in
		err := d.S.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredMessageUpdate,
		})
		if err != nil {
			d.log.ErrorErr(err)
			return
		}
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
