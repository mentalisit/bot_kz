package restapi

import (
	"fmt"
	"time"
	"whatsapp/models"
	"whatsapp/whatsapp/restapi/bridge"
	"whatsapp/whatsapp/restapi/compendium"
	"whatsapp/whatsapp/restapi/rs_bot2"

	"github.com/mentalisit/logger"
)

type Recover struct {
	log               *logger.Logger
	bridgeMessage     []models.ToBridgeMessage
	compendiumMessage []models.IncomingMessage
	rsBotV2Message    []models.InMessageV2
	bridge            *bridge.Client
	rs                *rs_bot2.Client
	compendiumNew     *compendium.Client
}

func NewRecover(log *logger.Logger) *Recover {
	r := &Recover{
		log:           log,
		bridge:        bridge.NewClient(log),
		rs:            rs_bot2.NewClient(log),
		compendiumNew: compendium.NewClient(log),
	}
	go r.trySend()
	return r
}

func (r *Recover) SendBridgeAppRecover(m models.ToBridgeMessage) {
	fmt.Printf("%s SendBridgeApp :NameRelay %s Sender %s Text %s ExtraLen %d\n",
		time.Now().Format(time.DateTime), m.Config.NameRelay, m.Sender, m.Text, len(m.Extra))
	err := r.bridge.SendToBridge(m)
	if err != nil {
		r.log.InfoStruct("SendBridgeApp err "+err.Error(), m)
		r.bridgeMessage = append(r.bridgeMessage, m)
	}
}

func (r *Recover) SendCompendiumAppRecover(m models.IncomingMessage) {
	fmt.Printf("%s SendCompendiumApp :%+v\n", time.Now().Format(time.DateTime), m)
	err := r.compendiumNew.SendToCompendium(m)
	if err != nil {
		r.log.InfoStruct("SendCompendiumApp err "+err.Error(), m)
		r.compendiumMessage = append(r.compendiumMessage, m)
	}
}

func (r *Recover) SendRsBotV2AppRecover(m models.InMessageV2) {
	err := r.rs.SendToRs2(m)
	if err != nil {
		r.log.InfoStruct("SendRsBotV2App err "+err.Error(), m)
		r.rsBotV2Message = append(r.rsBotV2Message, m)
	}
}

func (r *Recover) trySend() {
	for {
		// Проверка и отправка сообщений в rsBot2
		if len(r.rsBotV2Message) > 0 {
			for i := 0; i < len(r.rsBotV2Message); i++ {
				message := r.rsBotV2Message[i]
				err := r.rs.SendToRs2(message)
				if err == nil {
					// Если отправка успешна, удаляем сообщение из слайса
					r.rsBotV2Message = append(r.rsBotV2Message[:i], r.rsBotV2Message[i+1:]...)
					i-- // Сдвигаем индекс назад, чтобы корректно обработать оставшиеся элементы
				}
				time.Sleep(1 * time.Second)
			}
		}

		// Проверка и отправка сообщений в compendium
		if len(r.compendiumMessage) > 0 {
			for i := 0; i < len(r.compendiumMessage); i++ {
				message := r.compendiumMessage[i]
				err := r.compendiumNew.SendToCompendium(message)
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
	err = r.compendiumNew.Close()
	if err != nil {
		r.log.ErrorErr(err)
		return
	}

}
