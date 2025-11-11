package storage

import "rs/models"

type Subscribe interface {
	SubscribePing(s models.Subscribe) (subscribes []models.Subscribe)
	CheckSubscribe(s models.Subscribe) int
	Subscribe(s models.Subscribe)
	Unsubscribe(s models.Subscribe)
}
