package logger

import (
	"log"
	"os"
)

func EnableFileLogTranslating() {
	log.SetFlags(log.Ldate | log.Ltime)
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	log.SetOutput(file)
}
