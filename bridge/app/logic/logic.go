package logic

import (
	"bridge/models"
	"fmt"
	"strings"
	"time"
)

func (b *Bridge) Logic(m models.ToBridgeMessage) {
	b.in = m

	if b.in.Tip == "matrix" {
		if ok, conf := b.CacheCheckChannelConfigMatrix(b.in.ChatId); ok {
			b.in.Config = &conf
		} else {
			roomName := b.matrix.GetRoomName(b.in.ChatId)
			if roomName != "" {
				if ok, conf := b.CacheNameBridge(roomName); ok {
					b.in.Config = &conf
					//fmt.Printf("[Matrix] Resolved room %s to relay %s by name\n", b.in.ChatId, roomName)
				}
			}
			if b.in.Config == nil {
				fmt.Printf("[Matrix] Config not found for room %s (name: %s)\n", b.in.ChatId, roomName)
			}
		}
	}

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
