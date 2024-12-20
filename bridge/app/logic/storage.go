package logic

import (
	"bridge/models"
	"fmt"
)

func (b *Bridge) LoadConfig() {
	var i = 0
	var bridge string
	bc := b.storage.DBReadBridgeConfig()
	for _, configBridge := range bc {
		b.configs[configBridge.NameRelay] = configBridge
		i++
		bridge = bridge + fmt.Sprintf("%s, ", configBridge.NameRelay)
	}
	fmt.Printf("Загружено конфиг мостов %d : %s\n", i, bridge)
}
func (b *Bridge) CacheNameBridge(nameRelay string) (bool, models.BridgeConfig) {
	if len(b.configs) != 0 {
		for _, config := range b.configs {
			if config.NameRelay == nameRelay {
				return true, config
			}
		}
	}
	return false, models.BridgeConfig{}
}
func (b *Bridge) AddNewBridgeConfig() {
	b.configs[b.in.Config.NameRelay] = *b.in.Config
	b.storage.InsertBridgeChat(*b.in.Config)
}
func (b *Bridge) AddBridgeConfig() {
	b.configs[b.in.Config.NameRelay] = *b.in.Config
	b.storage.UpdateBridgeChat(*b.in.Config)
}

func (b *Bridge) CacheCheckChannelConfigDS(chatIdDs string) (bool, models.BridgeConfig) {
	for _, config := range b.configs {
		for _, ds := range config.ChannelDs {
			if ds.ChannelId == chatIdDs {
				return true, config
			}
		}
	}
	return false, models.BridgeConfig{}
}
func (b *Bridge) CacheCheckChannelConfigTg(chatIdTg string) (bool, models.BridgeConfig) {
	for _, config := range b.configs {
		for _, tg := range config.ChannelTg {
			if tg.ChannelId == chatIdTg {
				return true, config
			}
		}
	}
	return false, models.BridgeConfig{}
}
