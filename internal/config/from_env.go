package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func LoadEnv(enableDotenv bool) (*Config, error) {
	if enableDotenv {
		_ = godotenv.Load()
	}

	logLevel := loadString("LOGGER_LEVEL")
	logFilePath := loadString("LOGGER_FILE_PATH")
	enableFileLog := loadBool("LOGGER_ENABLE_FILE_LOG")

	loggerConfig := LoggerConfig{
		Level:         logLevel,
		FilePath:      logFilePath,
		EnableFileLog: enableFileLog,
	}

	watcherAddress := loadString("DOC_NOTIFIER_SERVICE_ADDRESS")
	watchedDirectories := strings.Split(loadString("DOC_NOTIFIER_WATCHED_DIRS"), ",")

	watcherConfig := WatcherConfig{
		Address:            watcherAddress,
		WatchedDirectories: watchedDirectories,
	}

	ocrAddress := loadString("OCR_SERVICE_ADDRESS")
	ocrMode := loadString("OCR_SERVICE_MODE")
	ocrTimeout := loadNumber("OCR_SERVICE_TIMEOUT", 64)

	ocrConfig := OcrConfig{
		Address: ocrAddress,
		Mode:    ocrMode,
		Timeout: time.Duration(ocrTimeout),
	}

	searcherAddress := loadString("DOC_SEARCHER_SERVICE_ADDRESS")
	searcherTimeout := loadNumber("DOC_SEARCHER_SERVICE_TIMEOUT", 64)

	searcherConfig := SearcherConfig{
		Address: searcherAddress,
		Timeout: time.Duration(searcherTimeout),
	}

	tokenizerAddress := loadString("TOKENIZER_SERVICE_ADDRESS")
	tokenizerMode := loadString("TOKENIZER_SERVICE_MODE")
	returnChunksFlag := loadBool("TOKENIZER_RETURN_CHUNKS")
	chunkSize := loadNumber("TOKENIZER_CHUNK_SIZE", 64)
	chunkOverlap := loadNumber("TOKENIZER_CHUNK_OVERLAP", 64)
	tokenizerTimeout := loadNumber("TOKENIZER_TIMEOUT_SECONDS", 64)
	chunkBySelfFlag := loadBool("TOKENIZER_CHUNK_BY_SELF")

	tokenizerConfig := TokenizerConfig{
		Address:      tokenizerAddress,
		Mode:         tokenizerMode,
		ChunkSize:    chunkSize,
		ChunkOverlap: chunkOverlap,
		ReturnChunks: returnChunksFlag,
		ChunkBySelf:  chunkBySelfFlag,
		Timeout:      time.Duration(tokenizerTimeout),
	}

	return &Config{
		Logger:    loggerConfig,
		Watcher:   watcherConfig,
		Ocr:       ocrConfig,
		Searcher:  searcherConfig,
		Tokenizer: tokenizerConfig,
	}, nil
}

func loadString(envName string) string {
	value, exists := os.LookupEnv(envName)
	if !exists {
		msg := fmt.Sprintf("faile to extract %s env var: %s", envName, value)
		log.Println(msg)
		return ""
	}
	return value
}

func loadNumber(envName string, bitSize int) int {
	value, exists := os.LookupEnv(envName)
	if !exists {
		msg := fmt.Sprintf("faile to extract %s env var: %s", envName, value)
		log.Println(msg)
		return 0
	}

	number, err := strconv.ParseInt(value, 10, bitSize)
	if err != nil {
		msg := fmt.Sprintf("faile to convert %s env var: %s", envName, value)
		log.Println(msg)
		return 0
	}

	return int(number)
}

func loadBool(envName string) bool {
	value, exists := os.LookupEnv(envName)
	if !exists {
		msg := fmt.Sprintf("faile to extract %s env var: %s", envName, value)
		log.Println(msg)
		return false
	}

	boolean, err := strconv.ParseBool(value)
	if err != nil {
		msg := fmt.Sprintf("faile to convert %s env var: %s", envName, value)
		log.Println(msg)
		return false
	}

	return boolean
}
