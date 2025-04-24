package server

import (
	"queue/models"
)

func (s *Server) Queue(level string) map[string][]models.QueueStruct {
	m := make(map[string][]models.QueueStruct)
	if level != "" {
		m = s.QueueLevel(level)
	} else {

		sborkzActive := s.kzbot.SelectSborkzActive()
		if len(sborkzActive) > 0 {
			m["RsBot"] = ConvertingSborkzToQueueStruct(sborkzActive)
		}

		rsbotQueueAll := s.queue.GetQueueAll()
		if len(rsbotQueueAll) > 0 {
			m["RsSoyzBot"] = ConvertingTumchaToQueueStruct(rsbotQueueAll)
		}

		if len(merged1) > 0 {
			m["Hades' Star RS Q"] = Merging(merged1)
		}
	}

	return m
}

func (s *Server) QueueLevel(level string) map[string][]models.QueueStruct {
	m := make(map[string][]models.QueueStruct)

	//My rsBot
	sborkzActive := s.kzbot.SelectSborkzActiveLevel(level)
	if len(sborkzActive) > 0 {
		m["RsBot"] = ConvertingSborkzToQueueStruct(sborkzActive)
	}

	//Tumcha Rs bot
	queueLevel := s.queue.GetQueueLevel(level)
	if len(queueLevel) > 0 {
		m["RsSoyzBot"] = ConvertingTumchaToQueueStruct(queueLevel)
	}

	if len(merged1) > 0 {
		ll := func(level string) []models.QueueStruct {
			var qq []models.QueueStruct
			for _, queueStruct := range merged1 {
				if queueStruct.Level == level {
					qq = append(qq, queueStruct)
				}
			}
			return Merging(qq)
		}
		result := ll(level)
		if len(result) > 0 {
			m["Hades' Star RS Q"] = result
		}
	}

	return m
}
