package storage

import "discord/models"

type Scoreboard interface {
	ScoreboardInsertParam(p models.ScoreboardParams)
	ScoreboardUpdateParam(p models.ScoreboardParams)
	ScoreboardReadWebhookChannel(webhookChannel string) *models.ScoreboardParams
	ScoreboardReadName(name string) *models.ScoreboardParams
	ReadEventScheduleAndMessage() (nextDateStart, nextDateStop, message string)
}
