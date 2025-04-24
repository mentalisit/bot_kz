package storage

import "rs/models"

type Scoreboard interface {
	ScoreboardInsertParam(p models.ScoreboardParams)
	ScoreboardUpdateParam(p models.ScoreboardParams)
	ScoreboardUpdateParamLastMessageId(p models.ScoreboardParams)
	ScoreboardUpdateParamScoreChannels(p models.ScoreboardParams)
	ScoreboardReadWebhookChannel(webhookChannel string) *models.ScoreboardParams
	ScoreboardReadName(name string) *models.ScoreboardParams
	ScoreboardReadAll() []models.ScoreboardParams
	ReadEventScheduleAndMessage() (nextDateStart, nextDateStop, message string)
}
