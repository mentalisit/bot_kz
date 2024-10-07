package restapi

import (
	"discord/models"
	"fmt"
	"github.com/mentalisit/logger"
	"time"
)

type Recover struct {
	log        *logger.Logger
	bridge     []models.ToBridgeMessage
	compendium []models.IncomingMessage
	rsBot      []models.InMessage
}

func NewRecover(log *logger.Logger) *Recover {
	r := &Recover{
		log: log,
	}
	go r.trySend()
	return r
}

func (r *Recover) SendBridgeAppRecover(m models.ToBridgeMessage) {
	fmt.Printf("%s SendBridgeApp :%+v\n", time.Now().Format(time.DateTime), m)
	err := SendBridgeApp(m)
	if err != nil {
		r.log.InfoStruct("SendBridgeApp", m)
		r.log.ErrorErr(err)
		r.bridge = append(r.bridge, m)
	}
}

func (r *Recover) SendCompendiumAppRecover(m models.IncomingMessage) {
	fmt.Printf("%s SendCompendiumApp :%+v\n", time.Now().Format(time.DateTime), m)
	err := SendCompendiumApp(m)
	if err != nil {
		r.log.InfoStruct("SendCompendiumApp", m)
		r.log.ErrorErr(err)
		r.compendium = append(r.compendium, m)
	}
}
func (r *Recover) SendRsBotAppRecover(m models.InMessage) {
	fmt.Printf("%s SendRsBotApp :%+v\n", time.Now().Format(time.DateTime), m)
	err := SendRsBotApp(m)
	if err != nil {
		r.log.InfoStruct("SendRsBotApp", m)
		r.log.ErrorErr(err)
		r.rsBot = append(r.rsBot, m)
	}
}
func (r *Recover) trySend() {
	for {
		// Проверка и отправка сообщений в rsBot
		if len(r.rsBot) > 0 {
			for i := 0; i < len(r.rsBot); i++ {
				message := r.rsBot[i]
				err := SendRsBotApp(message)
				if err == nil {
					// Если отправка успешна, удаляем сообщение из слайса
					r.rsBot = append(r.rsBot[:i], r.rsBot[i+1:]...)
					i-- // Сдвигаем индекс назад, чтобы корректно обработать оставшиеся элементы
				}
			}
		}

		// Проверка и отправка сообщений в compendium
		if len(r.compendium) > 0 {
			for i := 0; i < len(r.compendium); i++ {
				message := r.compendium[i]
				err := SendCompendiumApp(message)
				if err == nil {
					// Если отправка успешна, удаляем сообщение из слайса
					r.compendium = append(r.compendium[:i], r.compendium[i+1:]...)
					i-- // Сдвигаем индекс назад
				}
			}
		}

		// Проверка и отправка сообщений в bridge
		if len(r.bridge) > 0 {
			for i := 0; i < len(r.bridge); i++ {
				message := r.bridge[i]
				err := SendBridgeApp(message)
				if err == nil {
					// Если отправка успешна, удаляем сообщение из слайса
					r.bridge = append(r.bridge[:i], r.bridge[i+1:]...)
					i-- // Сдвигаем индекс назад
				}
			}
		}

		time.Sleep(10 * time.Second)
	}
}
