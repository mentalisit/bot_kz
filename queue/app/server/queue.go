package server

type QueueStruct struct {
	CorpName string
	Level    string
	Count    int
}

func (s *Server) Queue(level string) map[string][]QueueStruct {
	m := make(map[string][]QueueStruct)
	if level != "" {
		m = s.QueueLevel(level)
	} else {
		//
		rsbotQueueAll := s.queue.GetQueueAll()
		if len(rsbotQueueAll) > 0 {
			m["RsSoyzBot"] = ConvertingTumchaToQueueStruct(rsbotQueueAll)
		}

		//
		sborkzActive := s.kzbot.SelectSborkzActive()
		if len(sborkzActive) > 0 {
			m["RsBot"] = ConvertingSborkzToQueueStruct(sborkzActive)
		}

		if len(merged1) > 0 {
			m["Hades' Star RS Q"] = Merging(merged1)
		}
	}

	return m
}

func (s *Server) QueueLevel(level string) map[string][]QueueStruct {
	m := make(map[string][]QueueStruct)
	//Tumcha Rs bot
	queueLevel := s.queue.GetQueueLevel(level)
	if len(queueLevel) > 0 {
		m["RsSoyzBot"] = ConvertingTumchaToQueueStruct(queueLevel)
	}

	//My rsBot
	sborkzActive := s.kzbot.SelectSborkzActiveLevel(level)
	if len(sborkzActive) > 0 {
		m["RsBot"] = ConvertingSborkzToQueueStruct(sborkzActive)
	}

	if len(merged1) > 0 {
		ll := func(level string) []QueueStruct {
			var qq []QueueStruct
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
