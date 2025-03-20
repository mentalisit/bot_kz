package storage

import "rs/models"

type ConfigWebhook interface {
	ConfigWebhookInsert(u models.ConfigWebhook) error
	ConfigWebhookGetAll() ([]models.ConfigWebhook, error)
}
