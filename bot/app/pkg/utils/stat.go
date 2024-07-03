package utils

import (
	"fmt"
	"github.com/mentalisit/logger"
	"runtime"
)

var log *logger.Logger

func PrintGoroutine(l *logger.Logger) {
	if log == nil && l != nil {
		log = l
	}
	goroutine := runtime.NumGoroutine()
	text := fmt.Sprintf("Горутин  %d\n", goroutine)
	if log != nil {
		if goroutine > 120 {
			log.Info(text)
		} else if goroutine > 50 && goroutine%10 == 0 {
			log.Info(text)
		}
	}

	fmt.Println(text)
}
