package models

type ConfigWebhook struct {
	Name         string
	WebhookGame  string
	ScoreChannel string
}

type Battles struct {
	EventId  int
	CorpName string
	Name     string
	Level    int
	Points   int
}
