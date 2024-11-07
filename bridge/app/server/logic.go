package server

import (
	"bridge/models"
	"fmt"
	"strings"
	"time"
)

func (b *Bridge) Logic(m models.ToBridgeMessage) {
	b.in = m

	fmt.Printf("in bridge: %s %s relay %s channel %s lenFile:%d\n", m.Sender, m.Text, m.Config.HostRelay, m.ChatId, len(m.Extra))

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
