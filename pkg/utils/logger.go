package utils

import (
	"log"
	"sync"
	"time"
)

var logEntries []string
var logMutex sync.Mutex

func LogMessage(message string) {
	logMutex.Lock()
	defer logMutex.Unlock()

	entry := time.Now().Format("2006-01-02 15:04:05") + " - " + message
	logEntries = append(logEntries, entry)

	log.Println(entry)
}
