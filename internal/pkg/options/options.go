package options

import (
	"doc-notifier/internal/pkg/server"
	"errors"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"strings"
)

type Options struct {
	ServerAddress     string
	DocSearchAddress  string
	OcrServiceAddress string
	LlmServiceAddress string
	WatchDirectories  []string
	StoreChunksFlag   bool
	ReadRawFileFlag   bool
}

func LoadFromEnv() (*Options, error) {
	_ = godotenv.Load()

	var parseErr error
	var envExists bool
	var watcherPath []string
	var storeChunksFlag, readRawFileFlag bool
	var serverAddr, llmAddr, ocrAddr, docSearchAddr string

	storeChunksTmp, envExists := os.LookupEnv("IS_LOAD_CHUNKS")
	if !envExists {
		return nil, errors.New("")
	}
	if storeChunksFlag, parseErr = strconv.ParseBool(storeChunksTmp); parseErr != nil {
		return nil, errors.New("")
	}

	readRawFileTmp, envExists := os.LookupEnv("READ_RAW_FILE")
	if !envExists {
		return nil, errors.New("")
	}
	readRawFileFlag, parseErr = strconv.ParseBool(readRawFileTmp)

	if serverAddr, envExists = os.LookupEnv("host-address"); !envExists {
		return nil, errors.New("")
	}
	if llmAddr, envExists = os.LookupEnv("LLM_ADDRESS"); !envExists {
		return nil, errors.New("")
	}
	if ocrAddr, envExists = os.LookupEnv("OCR_ADDRESS"); !envExists {
		return nil, errors.New("")
	}
	if docSearchAddr, envExists = os.LookupEnv("DOC_ADDRESS"); !envExists {
		return nil, errors.New("")
	}

	watcherPathTmp, envExists := os.LookupEnv("WATCHER_DIR_PATHS")
	if !envExists {
		return nil, errors.New("")
	}
	watcherPath = strings.Split(watcherPathTmp, ",")

	return &Options{
		ServerAddress:     serverAddr,
		LlmServiceAddress: llmAddr,
		DocSearchAddress:  docSearchAddr,
		OcrServiceAddress: ocrAddr,
		WatchDirectories:  watcherPath,
		StoreChunksFlag:   storeChunksFlag,
		ReadRawFileFlag:   readRawFileFlag,
	}, nil
}

func ParseServerAddress(serverAddr string) *server.ServerOptions {
	tmp := strings.Split(serverAddr, ":")
	servicePort, _ := strconv.Atoi(tmp[1])
	return server.BuildOptions(tmp[0], servicePort)
}
