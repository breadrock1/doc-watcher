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
		return nil, errors.New("failed while parsing IS_LOAD_CHUNKS env variable")
	}
	if storeChunksFlag, parseErr = strconv.ParseBool(storeChunksTmp); parseErr != nil {
		return nil, errors.New("failed while parsing IS_LOAD_CHUNKS env variable")
	}

	readRawFileTmp, envExists := os.LookupEnv("READ_RAW_FILE")
	if !envExists {
		return nil, errors.New("failed while parsing READ_RAW_FILE env variable")
	}
	if readRawFileFlag, parseErr = strconv.ParseBool(readRawFileTmp); parseErr != nil {
		return nil, errors.New("failed while parsing READ_RAW_FILE env variable")
	}

	if serverAddr, envExists = os.LookupEnv("HOST_ADDRESS"); !envExists {
		return nil, errors.New("failed while parsing HOST_ADDRESS env variable")
	}
	if llmAddr, envExists = os.LookupEnv("LLM_ADDRESS"); !envExists {
		return nil, errors.New("failed while parsing LLM_ADDRESS env variable")
	}
	if ocrAddr, envExists = os.LookupEnv("OCR_ADDRESS"); !envExists {
		return nil, errors.New("failed while parsing OCR_ADDRESS env variable")
	}
	if docSearchAddr, envExists = os.LookupEnv("DOC_ADDRESS"); !envExists {
		return nil, errors.New("failed while parsing DOC_ADDRESS env variable")
	}

	watcherPathTmp, envExists := os.LookupEnv("WATCHER_DIR_PATHS")
	if !envExists {
		return nil, errors.New("failed while parsing WATCHER_DIR_PATHS env variable")
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
