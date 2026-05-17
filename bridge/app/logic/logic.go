package logic

import (
	"bridge/models"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (b *Bridge) Logic(m models.ToBridgeMessage) {
	guildV2, err := b.db.GuildGetChannel(m.GuildId)
	if err == nil && guildV2 != nil && guildV2.GId != uuid.Nil {

	}

	b.in = m

	if b.in.Config == nil {
		fmt.Printf("Config is nil for Tip: %s, ChatId: %s\n", b.in.Tip, b.in.ChatId)
		return
	}

	fmt.Printf("in bridge: %s %s relay %s channel %s lenFile:%d\n", b.in.Sender, b.in.Text, b.in.Config.HostRelay, b.in.ChatId, len(b.in.Extra))

	if strings.HasPrefix(b.in.Text, ".poll") {
		b.ifPoll()
		return
	}

	if strings.HasPrefix(b.in.Text, ".") {
		go func() {
			b.Command()
			time.Sleep(5 * time.Second)
			b.LoadConfig()
		}()
		return
	} else {
		b.logicMessage()
	}

}

func (b *Bridge) logicMessage() {
	if b.checkingForIdenticalMessage() {
		return
	}
	if b.in.Tip == "delDs" || b.in.Tip == "delWa" {
		b.RemoveMessage()
		return
	}
	if b.in.Tip == "dse" || b.in.Tip == "tge" || b.in.Tip == "wae" {
		//todo need tested
		b.EditMessage()
	}

	b.logicSendMessage()
}
