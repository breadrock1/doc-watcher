package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"doc-watcher/internal/embeddings"
	"doc-watcher/internal/ocr"
	"doc-watcher/internal/searcher"
	"doc-watcher/internal/server"
	"doc-watcher/internal/watcher"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Ocr        ocr.Config
	Searcher   searcher.Config
	Server     server.Config
	Embeddings embeddings.Config
	Watcher    watcher.Config
}

func FromFile(filePath string) (*Config, error) {
	config := &Config{}

	viperInstance := viper.New()
	viperInstance.AutomaticEnv()
	viperInstance.SetConfigFile(filePath)

	viperInstance.SetDefault("server.Address", "0.0.0.0:2893")
	viperInstance.SetDefault("server.LoggerLevel", "INFO")

	viperInstance.SetDefault("ocr.Address", "ocr:8004")
	viperInstance.SetDefault("ocr.EnableSSL", false)
	viperInstance.SetDefault("ocr.Timeout", 300)

	viperInstance.SetDefault("searcher.Address", "doc-searcher:2892")
	viperInstance.SetDefault("searcher.EnableSSL", false)

	viperInstance.SetDefault("embeddings.Address", "embeddings:8001")
	viperInstance.SetDefault("embeddings.EnableSSL", false)
	viperInstance.SetDefault("embeddings.ChunkSize", 800)
	viperInstance.SetDefault("embeddings.ChunkOverlap", 100)
	viperInstance.SetDefault("embeddings.ReturnChunks", false)
	viperInstance.SetDefault("embeddings.ChunkBySelf", false)

	viperInstance.SetDefault("watcher.Address", "cloud-storage:2894")
	viperInstance.SetDefault("watcher.Username", "minio-root")
	viperInstance.SetDefault("watcher.Password", "minio-root")
	viperInstance.SetDefault("watcher.EnableSSL", false)
	viperInstance.SetDefault("watcher.WatchedDirectories", []string{"common-folder"})
	viperInstance.SetDefault("watcher.CacheExpire", 10)
	viperInstance.SetDefault("watcher.CacheCleanInterval", 30)

	if err := viperInstance.ReadInConfig(); err != nil {
		confErr := fmt.Errorf("failed while reading config file %s: %w", filePath, err)
		return config, confErr
	}

	if err := viperInstance.Unmarshal(config); err != nil {
		confErr := fmt.Errorf("failed while unmarshaling config file %s: %w", filePath, err)
		return config, confErr
	}

	return config, nil
}

func LoadEnv(enableDotenv bool) (*Config, error) {
	if enableDotenv {
		_ = godotenv.Load()
	}

	serverAddr := loadString("DOC_WATCHER_SERVER_ADDRESS")
	loggerLevel := loadString("DOC_WATCHER_SERVER_LOGGER_LEVEL")
	serverConfig := server.Config{
		Address:     serverAddr,
		LoggerLevel: loggerLevel,
	}

	ocrAddr := loadString("DOC_WATCHER_OCR_ADDRESS")
	enableSSL := loadBool("DOC_WATCHER_OCR_ENABLE_SSL")
	ocrTimeout := loadNumber("DOC_WATCHER_OCR_TIMEOUT")
	ocrConfig := ocr.Config{
		Address:   ocrAddr,
		EnableSSL: enableSSL,
		Timeout:   time.Duration(ocrTimeout),
	}

	searchAddr := loadString("DOC_WATCHER_SEARCHER_ADDRESS")
	searchEnableSSL := loadBool("DOC_WATCHER_SEARCHER_ENABLE_SSL")
	searchConfig := searcher.Config{
		Address:   searchAddr,
		EnableSSL: searchEnableSSL,
	}

	embAddress := loadString("DOC_WATCHER_EMBEDDINGS_ADDRESS")
	embEnableSSL := loadBool("DOC_WATCHER_EMBEDDINGS_ENABLE_SSL")
	embChunkSize := loadNumber("DOC_WATCHER_EMBEDDINGS_CHUNK_SIZE")
	embChunkOverlap := loadNumber("DOC_WATCHER_EMBEDDINGS_CHUNK_OVERLAP")
	embRetChunks := loadBool("DOC_WATCHER_EMBEDDINGS_RETURN_CHUNKS")
	embChunkBySelf := loadBool("DOC_WATCHER_EMBEDDINGS_SELF_CHUNK")
	embConfig := embeddings.Config{
		Address:      embAddress,
		EnableSSL:    embEnableSSL,
		ChunkSize:    embChunkSize,
		ChunkOverlap: embChunkOverlap,
		ReturnChunks: embRetChunks,
		ChunkBySelf:  embChunkBySelf,
	}

	watchAddress := loadString("DOC_WATCHER_ADDRESS")
	watchEnableSSL := loadBool("DOC_WATCHER_ENABLE_SSL")
	watchUsername := loadString("DOC_WATCHER_USERNAME")
	watchPassword := loadString("DOC_WATCHER_PASSWORD")
	watchCacheExpire := loadNumber("DOC_WATCHER_CACHE_EXPIRE")
	watchCacheCleanInterval := loadNumber("DOC_WATCHER_CACHE_CLEAN_INTERVAL")

	watchDirectories := strings.Split(loadString("DOC_WATCHER_WATCHED_DIRS"), ",")
	watchConfig := watcher.Config{
		Address:            watchAddress,
		EnableSSL:          watchEnableSSL,
		Username:           watchUsername,
		Password:           watchPassword,
		CacheExpire:        time.Duration(watchCacheExpire),
		CacheCleanInterval: time.Duration(watchCacheCleanInterval),
		WatchedDirectories: watchDirectories,
	}

	return &Config{
		Ocr:        ocrConfig,
		Searcher:   searchConfig,
		Server:     serverConfig,
		Embeddings: embConfig,
		Watcher:    watchConfig,
	}, nil
}

func loadString(envName string) string {
	value, exists := os.LookupEnv(envName)
	if !exists {
		msg := fmt.Sprintf("failed to extract %s env var: %s", envName, value)
		log.Println(msg)
		return ""
	}
	return value
}

func loadNumber(envName string) int {
	value, exists := os.LookupEnv(envName)
	if !exists {
		msg := fmt.Sprintf("failed to extract %s env var: %s", envName, value)
		log.Println(msg)
		return 0
	}

	number, err := strconv.Atoi(value)
	if err != nil {
		msg := fmt.Sprintf("failed to convert %s env var: %s", envName, value)
		log.Println(msg)
		return 0
	}

	return number
}

func loadBool(envName string) bool {
	value, exists := os.LookupEnv(envName)
	if !exists {
		msg := fmt.Sprintf("failed to extract %s env var: %s", envName, value)
		log.Println(msg)
		return false
	}

	boolean, err := strconv.ParseBool(value)
	if err != nil {
		msg := fmt.Sprintf("failed to convert %s env var: %s", envName, value)
		log.Println(msg)
		return false
	}

	return boolean
}
