package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/mentalisit/restapi/models"
)

func (d *Db) loadConfig() {
	d.reloadRSBotConfig()
	d.reloadBridgeConfig()
	d.reloadKZBotConfig()
}

func (d *Db) StartConfigWatcher(ctx context.Context) {
	for {
		listener := pq.NewListener(
			d.dns,
			10*time.Second,
			time.Minute,
			func(event pq.ListenerEventType, err error) {
				if err != nil {
					fmt.Println("Listener error:", err)
				}
			},
		)

		err := listener.Listen("config_updates")
		if err != nil {
			fmt.Println("Ошибка LISTEN:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		fmt.Println("Успешно подписались на уведомления БД (канал: config_updates)")

	loop:
		for {
			select {
			case <-ctx.Done():
				_ = listener.Close()
				return

			case notification := <-listener.Notify:
				if notification == nil {
					continue
				}

				table := notification.Extra

				fmt.Printf("Получено уведомление об изменении в: %s\n", table)

				switch table {

				case "rs_bot.config_rs":
					d.reloadRSBotConfig()

				case "rs_bot.bridge_config":
					d.reloadBridgeConfig()

				case "kzbot.config":
					d.reloadKZBotConfig()

				default:
					d.log.Info("Неизвестная таблица в уведомлении: " + table)
				}
			}

			// Проверка соединения
			err = listener.Ping()
			if err != nil {
				fmt.Println("Потеряно соединение с PostgreSQL:", err)
				break loop
			}
		}

		fmt.Println("Переподключение через 5 секунд...")

		select {
		case <-ctx.Done():
			return

		case <-time.After(5 * time.Second):
		}
	}
}

func (d *Db) reloadKZBotConfig() {
	configRs, _ := d.ReadConfigRs()
	m := make(map[string]models.CorporationConfig)
	for _, r := range configRs {
		if r.DsChannel != "" {
			m[r.DsChannel] = r
		}
		if r.TgChannel != "" {
			m[r.TgChannel] = r
		}
		if r.WaChannel != "" {
			m[r.WaChannel] = r
		}
	}
	d.Lock()
	d.KzBotConfig = m
	d.Unlock()
	fmt.Println("Кэш KZBotConfig обновлен")
}

func (d *Db) reloadBridgeConfig() {
	bridgeConfig := d.DBReadBridgeConfig()
	m := make(map[string]models.Bridge2Config)
	for _, config := range bridgeConfig {
		for _, configs := range config.Channel {
			for _, ch := range configs {
				m[ch.ChannelId] = config
			}
		}
	}
	d.Lock()
	d.BridgeConfig = m
	d.Unlock()
	fmt.Println("Кэш BridgeConfig обновлен")
}

func (d *Db) reloadRSBotConfig() {
	configV2 := d.ReadConfigV2()
	m := make(map[string]models.CorporationConfigV2)
	for _, v2 := range configV2 {
		for ch, _ := range v2.Channels {
			m[ch] = v2
		}
	}
	d.Lock()
	d.RsBotConfig = m
	d.Unlock()
	fmt.Println("Кэш RSBotConfig обновлен")
}

func (d *Db) CheckBridgeChannel(channelId string) (bool, models.Bridge2Config) {
	d.RLock()
	defer d.RUnlock()
	config, exist := d.BridgeConfig[channelId]
	if exist {
		return true, config
	}
	return false, models.Bridge2Config{}
}

func (d *Db) CheckKzChannel(channelId string) (bool, models.CorporationConfig) {
	d.RLock()
	defer d.RUnlock()
	config, exist := d.KzBotConfig[channelId]
	if exist {
		return true, config
	}
	return false, models.CorporationConfig{}
}

func (d *Db) CheckRsChannel(channelId string) (bool, models.CorporationConfigV2) {
	d.RLock()
	defer d.RUnlock()
	config, exist := d.RsBotConfig[channelId]
	if exist {
		return true, config
	}
	return false, models.CorporationConfigV2{}
}
