package logic

import (
	"bridge/models"
	"fmt"
)

func (b *Bridge) LoadConfig() {
	var i = 0
	var bridge string
	//bc := b.storage.DBReadBridgeConfig()
	//if len(bc) > 0 {
	//	for _, conf := range bc {
	//		newconfig := models.Bridge2Config{
	//			NameRelay:         conf.NameRelay,
	//			HostRelay:         conf.HostRelay,
	//			Role:              conf.Role,
	//			Channel:           make(map[string][]models.Bridge2Configs),
	//			ForbiddenPrefixes: conf.ForbiddenPrefixes,
	//		}
	//		if len(conf.ChannelDs) > 0 {
	//			for _, d := range conf.ChannelDs {
	//				if newconfig.Channel["ds"] == nil {
	//					newconfig.Channel["ds"] = []models.Bridge2Configs{}
	//				}
	//				newconfig.Channel["ds"] = append(newconfig.Channel["ds"], models.Bridge2Configs{
	//					ChannelId:       d.ChannelId,
	//					GuildId:         d.GuildId,
	//					CorpChannelName: d.CorpChannelName,
	//					AliasName:       d.AliasName,
	//					MappingRoles:    d.MappingRoles,
	//				})
	//			}
	//		}
	//		if len(conf.ChannelTg) > 0 {
	//			for _, d := range conf.ChannelTg {
	//				if newconfig.Channel["tg"] == nil {
	//					newconfig.Channel["tg"] = []models.Bridge2Configs{}
	//				}
	//				spl := strings.Split(d.ChannelId, "/")
	//				newconfig.Channel["tg"] = append(newconfig.Channel["tg"], models.Bridge2Configs{
	//					ChannelId:       d.ChannelId,
	//					GuildId:         spl[0],
	//					CorpChannelName: d.CorpChannelName,
	//					AliasName:       d.AliasName,
	//					MappingRoles:    d.MappingRoles,
	//				})
	//			}
	//		}
	//		b.storage.InsertBridge2Chat(newconfig)
	//		b.storage.DeleteBridgeChat(conf)
	//	}
	//}
	bcNew := b.storage.DBReadBridgeConfig2()
	for _, configBridge := range bcNew {
		b.configs[configBridge.NameRelay] = configBridge
		i++
		bridge = bridge + fmt.Sprintf("%s, ", configBridge.NameRelay)
	}
	fmt.Printf("Загружено конфиг мостов %d : %s\n", i, bridge)

}

func (b *Bridge) CacheNameBridge(nameRelay string) (bool, models.Bridge2Config) {
	if len(b.configs) != 0 {
		for _, config := range b.configs {
			if config.NameRelay == nameRelay {
				return true, config
			}
		}
	}
	return false, models.Bridge2Config{}
}
func (b *Bridge) AddNewBridgeConfig() {
	b.configs[b.in.Config.NameRelay] = *b.in.Config
	b.storage.InsertBridge2Chat(*b.in.Config)
}
func (b *Bridge) AddBridgeConfig() {
	b.configs[b.in.Config.NameRelay] = *b.in.Config
	b.storage.UpdateBridge2Chat(*b.in.Config)
}

func (b *Bridge) CacheCheckChannelConfigDS(chatIdDs string) (bool, models.Bridge2Config) {
	for _, config := range b.configs {
		for _, ds := range config.Channel["ds"] {
			if ds.ChannelId == chatIdDs {
				return true, config
			}
		}
	}
	return false, models.Bridge2Config{}
}
func (b *Bridge) CacheCheckChannelConfigTg(chatIdTg string) (bool, models.Bridge2Config) {
	for _, config := range b.configs {
		for _, tg := range config.Channel["tg"] {
			if tg.ChannelId == chatIdTg {
				return true, config
			}
		}
	}
	return false, models.Bridge2Config{}
}
func (b *Bridge) CacheCheckChannelConfigWA(chatIdTg string) (bool, models.Bridge2Config) {
	for _, config := range b.configs {
		for _, tg := range config.Channel["wa"] {
			if tg.ChannelId == chatIdTg {
				return true, config
			}
		}
	}
	return false, models.Bridge2Config{}
}
