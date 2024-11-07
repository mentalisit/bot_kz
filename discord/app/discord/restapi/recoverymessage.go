package restapi

import (
	"discord/discord/restapi/bridge"
	"discord/discord/restapi/rs_bot"
	"discord/models"
	"fmt"
	"github.com/mentalisit/logger"
	"time"
)

type Recover struct {
	log               *logger.Logger
	bridgeMessage     []models.ToBridgeMessage
	compendiumMessage []models.IncomingMessage
	rsBotMessage      []models.InMessage
	bridge            *bridge.Client
	rs                *rs_bot.Client
}

func NewRecover(log *logger.Logger) *Recover {
	r := &Recover{
		log:    log,
		bridge: bridge.NewClient(log),
		rs:     rs_bot.NewClient(log),
	}
	go r.trySend()
	return r
}

func (r *Recover) SendBridgeAppRecover(m models.ToBridgeMessage) {
	fmt.Printf("%s SendBridgeApp :%+v\n", time.Now().Format(time.DateTime), m)
	err := r.bridge.SendToBridge(m)
	if err != nil {
		r.log.InfoStruct("SendBridgeApp err "+err.Error(), m)
		r.bridgeMessage = append(r.bridgeMessage, m)
	}
}

func (r *Recover) SendCompendiumAppRecover(m models.IncomingMessage) {
	fmt.Printf("%s SendCompendiumApp :%+v\n", time.Now().Format(time.DateTime), m)
	err := SendCompendiumApp(m)
	if err != nil {
		r.log.InfoStruct("SendCompendiumApp err "+err.Error(), m)
		r.compendiumMessage = append(r.compendiumMessage, m)
	}
}
func (r *Recover) SendRsBotAppRecover(m models.InMessage) {
	fmt.Printf("%s SendRsBotApp :%+v\n", time.Now().Format(time.DateTime), m)
	err := r.rs.SendToRs(m)
	if err != nil {
		r.log.InfoStruct("SendRsBotApp err "+err.Error(), m)
		r.rsBotMessage = append(r.rsBotMessage, m)
	}
}
func (r *Recover) trySend() {
	for {
		// Проверка и отправка сообщений в rsBot
		if len(r.rsBotMessage) > 0 {
			for i := 0; i < len(r.rsBotMessage); i++ {
				message := r.rsBotMessage[i]
				err := r.rs.SendToRs(message)
				if err == nil {
					// Если отправка успешна, удаляем сообщение из слайса
					r.rsBotMessage = append(r.rsBotMessage[:i], r.rsBotMessage[i+1:]...)
					i-- // Сдвигаем индекс назад, чтобы корректно обработать оставшиеся элементы
				}
				time.Sleep(1 * time.Second)
			}
		}

		// Проверка и отправка сообщений в compendium
		if len(r.compendiumMessage) > 0 {
			for i := 0; i < len(r.compendiumMessage); i++ {
				message := r.compendiumMessage[i]
				err := SendCompendiumApp(message)
				if err == nil {
					// Если отправка успешна, удаляем сообщение из слайса
					r.compendiumMessage = append(r.compendiumMessage[:i], r.compendiumMessage[i+1:]...)
					i-- // Сдвигаем индекс назад
				}
				time.Sleep(1 * time.Second)
			}
		}

		// Проверка и отправка сообщений в bridge
		if len(r.bridgeMessage) > 0 {
			for i := 0; i < len(r.bridgeMessage); i++ {
				message := r.bridgeMessage[i]
				err := r.bridge.SendToBridge(message)
				if err == nil {
					// Если отправка успешна, удаляем сообщение из слайса
					r.bridgeMessage = append(r.bridgeMessage[:i], r.bridgeMessage[i+1:]...)
					i-- // Сдвигаем индекс назад
				}
				time.Sleep(1 * time.Second)
			}
		}

		time.Sleep(10 * time.Second)
	}
}
func (r *Recover) Close() {
	err := r.bridge.Close()
	if err != nil {
		r.log.ErrorErr(err)
		return
	}
	err = r.rs.Close()
	if err != nil {
		r.log.ErrorErr(err)
		return
	}

}
