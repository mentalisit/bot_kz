package server

import (
	"bridge/models"
	"strconv"
)

func (b *Bridge) RemoveMessage() {
	if len(b.messages) > 0 {
		var mem []models.BridgeTempMemory
		for _, memory := range b.messages {
			if b.ifMessageIdDs(&memory, b.in.MesId) {
				for _, s := range memory.MessageDs {
					go b.discord.DeleteMessageDs(s.ChatId, s.MessageId)
				}
				for _, s := range memory.MessageTg {
					mid, err := strconv.Atoi(s.MessageId)
					if err != nil {
						return
					}
					go b.telegram.DeleteMessage(s.ChatId, strconv.Itoa(mid))
				}
			} else {
				mem = append(mem, memory)
			}
		}
		b.messages = mem
	}
}
func (b *Bridge) ifMessageIdDs(memory *models.BridgeTempMemory, MesId string) bool {
	for _, s := range memory.MessageDs {
		if s.MessageId == MesId {
			return true
		}
	}
	return false
}

func (b *Bridge) ifMessageIdTg(memory *models.BridgeTempMemory, MesId string) bool {
	for _, s := range memory.MessageTg {
		if s.MessageId == MesId {
			return true
		}
	}
	return false
}
