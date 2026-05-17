package models

import "time"

type ChatChannel struct {
	ID          string    `json:"id"`
	GID         string    `json:"gid"`
	Name        string    `json:"name"`
	CreatorUUID string    `json:"creator_uuid"`
	CreatedAt   time.Time `json:"created_at"`
}
