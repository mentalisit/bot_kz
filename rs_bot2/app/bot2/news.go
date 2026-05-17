package bot2

import (
	"fmt"
	"rs/models"
	"time"
)

func (b *Bot) readNews() {
	if time.Now().Minute()%5 == 0 {
		en, ru, ua := b.client.Ds.ReadNews()
		if en != "" && ru != "" && ua != "" {

			getText := func(n models.News) string {
				if n.Language == "en" {
					return en
				} else if n.Language == "ua" {
					return ua
				}
				return ru
			}

			for _, n := range b.storage.ReadNewsAll() {
				text := fmt.Sprintf("%s \n%s", "Hades' Star Official", getText(n))
				if n.MessengerType == ds {
					b.client.Ds.Send(n.ChatId, getText(n))
				} else if n.MessengerType == tg {
					b.client.Tg.SendChannel(n.ChatId, text)
				} else if n.MessengerType == wa {
					b.client.Wa.SendChannel(n.ChatId, text)
				} else {
					b.log.Info("please implement news for " + n.MessengerType)
				}

			}
		}
	}
}
