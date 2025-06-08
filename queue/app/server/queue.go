package server

import (
	"fmt"
	"queue/models"
	"queue/rsq"
	"sync"
)

func (s *Server) Queue(level string) map[string][]models.QueueStruct {
	m := make(map[string][]models.QueueStruct)
	if level != "" {
		m = s.QueueLevel(level)
	} else {
		var wg sync.WaitGroup

		wg.Add(3)
		fmt.Printf("Status: ")

		go func() {
			sborkzActive := s.kzbot.SelectSborkzActive()
			if len(sborkzActive) > 0 {
				m["RsBot"] = ConvertingSborkzToQueueStruct(sborkzActive)
			}
			fmt.Printf(" KzBot ")
			wg.Done()
		}()

		go func() {
			rsbotQueueAll := s.queue.GetQueueAll()
			if len(rsbotQueueAll) > 0 {
				m["RsSoyzBot"] = ConvertingTumchaToQueueStruct(rsbotQueueAll)
			}
			fmt.Printf(" rsSoyzBot ")
			wg.Done()
		}()

		go func() {
			dataCaprican := rsq.GetDataCaprican()
			if len(dataCaprican) > 0 {
				m["RSQ"] = dataCaprican
			}
			fmt.Printf(" RSQ ")
			wg.Done()
		}()

		wg.Wait()
		fmt.Printf(" ...complite all\n")
	}

	return m
}

func (s *Server) QueueLevel(level string) map[string][]models.QueueStruct {
	m := make(map[string][]models.QueueStruct)
	var wg sync.WaitGroup

	wg.Add(3)
	go func() {
		//My rsBot
		sborkzActive := s.kzbot.SelectSborkzActiveLevel(level)
		if len(sborkzActive) > 0 {
			m["RsBot"] = ConvertingSborkzToQueueStruct(sborkzActive)
		}
		wg.Done()
	}()

	go func() {
		//Tumcha Rs bot
		queueLevel := s.queue.GetQueueLevel(level)
		if len(queueLevel) > 0 {
			m["RsSoyzBot"] = ConvertingTumchaToQueueStruct(queueLevel)
		}
		wg.Done()
	}()

	go func() {
		dataCaprican := rsq.GetDataCaprican()
		if len(dataCaprican) > 0 {
			ll := func(level string) []models.QueueStruct {
				var qq []models.QueueStruct
				for _, queueStruct := range dataCaprican {
					if queueStruct.Level == level {
						qq = append(qq, queueStruct)
					}
				}
				return qq
			}
			result := ll(level)
			if len(result) > 0 {
				m["RSQ"] = result
			}
		}
		wg.Done()
	}()

	wg.Wait()
	return m
}
