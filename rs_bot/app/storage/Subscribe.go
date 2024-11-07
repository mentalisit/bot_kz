package storage

type Subscribe interface {
	SubscribePing(nameMention, lvlkz string, tipPing int, TgChannel string) string
	CheckSubscribe(name, lvlkz string, TgChannel string, tipPing int) int
	Subscribe(name, nameMention, lvlkz string, tipPing int, TgChannel string)
	Unsubscribe(name, lvlkz string, TgChannel string, tipPing int)
}
