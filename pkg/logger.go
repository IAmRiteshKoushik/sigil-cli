package pkg

import (
	"log"
	"os"
	"time"
)

func SetupLogger() {
	logFile, err := os.OpenFile("sigil.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}

	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("[%s] Application started", time.Now().Format("2006-01-02 15:04:05"))
}
