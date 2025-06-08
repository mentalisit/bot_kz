package storage

import "discord/models"

type Scoreboard interface {
	ScoreboardUpdateParamLastMessageId(p models.ScoreboardParams)
	ScoreboardReadWebhookChannel(webhookChannel string) *models.ScoreboardParams
	ScoreboardReadAll() []models.ScoreboardParams
	ReadEventScheduleAndMessage() (nextDateStart, nextDateStop, message string)
	ScoreboardInsertParam(p models.ScoreboardParams)
}
