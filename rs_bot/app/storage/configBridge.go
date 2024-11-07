package storage

import "rs/models"

type ConfigBridge interface {
	DBReadBridgeConfig() []models.BridgeConfig
	FindBridgeConfigByChannelId(channelId string) (*models.BridgeConfig, error)
}
