package utils

import (
	"fmt"
	"github.com/mentalisit/logger"
	"runtime"
	"time"
)

var log *logger.Logger

func PrintGoroutine(l *logger.Logger) {
	if log == nil && l != nil {
		log = l
	}
	goroutine := runtime.NumGoroutine()
	tm := time.Now()
	mdate := (tm.Format("2006-01-02"))
	mtime := (tm.Format("15:04"))
	text := fmt.Sprintf(" %s %s Горутин  %d\n", mdate, mtime, goroutine)
	if log != nil {
		if goroutine > 120 {
			log.Info(text)
		} else if goroutine > 50 && goroutine%10 == 0 {
			log.Info(text)
		}
	}

	fmt.Println(text)
}
func PrintGoroutinesStack() {
	buf := make([]byte, 1<<16)
	stacklen := runtime.Stack(buf, true)
	fmt.Printf("=== Goroutine Stack ===\n%s\n", buf[:stacklen])
}
