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
	WatcherServiceAddress string
	WatchedDirectories    []string

	OcrServiceAddress string
	OcrServiceMode    string

	DocSearchAddress string

	TokenizerServiceAddress string
	TokenizerServiceMode    string
	TokenizerChunkSize      int
	TokenizerChunkOverlap   int
	TokenizerReturnChunks   bool
	TokenizerChunkBySelf    bool
	TokenizerTimeout        uint
}

func LoadFromEnv(disabledDotenv bool) (*Options, error) {
	if !disabledDotenv {
		_ = godotenv.Load()
	}

	var envExists bool
	var tmpOptionVar string
	var parseOptionErr error

	var watchedDirectories []string
	var tokenizerTimeout uint64
	var chunkSize, chunkOverlap int64
	var returnChunksFlag, chunkBySelfFlag bool
	var tokenizerServiceAddr, tokenizerServiceMode string
	var notifierAddr, docSearchAddr, ocrServiceAddr, ocrServiceMode string

	if notifierAddr, envExists = os.LookupEnv("DOC_NOTIFIER_SERVICE_ADDRESS"); !envExists {
		return nil, errors.New("failed while parsing DOC_NOTIFIER_SERVICE_ADDRESS env variable")
	}
	tmpOptionVar, envExists = os.LookupEnv("DOC_NOTIFIER_WATCHED_DIRS")
	if !envExists {
		return nil, errors.New("failed while parsing DOC_NOTIFIER_WATCHED_DIRS env variable")
	}
	watchedDirectories = strings.Split(tmpOptionVar, ",")

	if ocrServiceAddr, envExists = os.LookupEnv("OCR_SERVICE_ADDRESS"); !envExists {
		return nil, errors.New("failed while parsing OCR_SERVICE_ADDRESS env variable")
	}
	if ocrServiceMode, envExists = os.LookupEnv("OCR_SERVICE_MODE"); !envExists {
		return nil, errors.New("failed while parsing OCR_SERVICE_MODE env variable")
	}

	if docSearchAddr, envExists = os.LookupEnv("DOC_SEARCHER_SERVICE_ADDRESS"); !envExists {
		return nil, errors.New("failed while parsing DOC_SEARCHER_SERVICE_ADDRESS env variable")
	}

	if tokenizerServiceAddr, envExists = os.LookupEnv("TOKENIZER_SERVICE_ADDRESS"); !envExists {
		return nil, errors.New("failed while parsing TOKENIZER_SERVICE_ADDRESS env variable")
	}
	if tokenizerServiceMode, envExists = os.LookupEnv("TOKENIZER_MODE"); !envExists {
		return nil, errors.New("failed while parsing TOKENIZER_MODE env variable")
	}

	tmpOptionVar, envExists = os.LookupEnv("TOKENIZER_RETURN_CHUNKS")
	if !envExists {
		return nil, errors.New("failed while parsing TOKENIZER_RETURN_CHUNKS env variable")
	}
	if returnChunksFlag, parseOptionErr = strconv.ParseBool(tmpOptionVar); parseOptionErr != nil {
		return nil, errors.New("failed while parsing TOKENIZER_RETURN_CHUNKS env variable")
	}

	tmpOptionVar, envExists = os.LookupEnv("TOKENIZER_CHUNK_SIZE")
	if !envExists {
		return nil, errors.New("failed while parsing TOKENIZER_CHUNK_SIZE env variable")
	}
	if chunkSize, parseOptionErr = strconv.ParseInt(tmpOptionVar, 10, 64); parseOptionErr != nil {
		return nil, errors.New("failed while parsing TOKENIZER_CHUNK_SIZE env variable")
	}

	tmpOptionVar, envExists = os.LookupEnv("TOKENIZER_CHUNK_OVERLAP")
	if !envExists {
		return nil, errors.New("failed while parsing TOKENIZER_CHUNK_OVERLAP env variable")
	}
	if chunkOverlap, parseOptionErr = strconv.ParseInt(tmpOptionVar, 10, 64); parseOptionErr != nil {
		return nil, errors.New("failed while parsing TOKENIZER_CHUNK_OVERLAP env variable")
	}

	tmpOptionVar, envExists = os.LookupEnv("TOKENIZER_TIMEOUT_SECONDS")
	if !envExists {
		return nil, errors.New("failed while parsing TOKENIZER_TIMEOUT_SECONDS env variable")
	}
	if tokenizerTimeout, parseOptionErr = strconv.ParseUint(tmpOptionVar, 10, 64); parseOptionErr != nil {
		return nil, errors.New("failed while parsing TOKENIZER_TIMEOUT_SECONDS env variable")
	}

	tmpOptionVar, envExists = os.LookupEnv("TOKENIZER_CHUNK_BY_SELF")
	if !envExists {
		return nil, errors.New("failed while parsing TOKENIZER_CHUNK_BY_SELF env variable")
	}
	if chunkBySelfFlag, parseOptionErr = strconv.ParseBool(tmpOptionVar); parseOptionErr != nil {
		return nil, errors.New("failed while parsing TOKENIZER_CHUNK_BY_SELF env variable")
	}

	return &Options{
		WatcherServiceAddress: notifierAddr,
		WatchedDirectories:    watchedDirectories,

		OcrServiceAddress: ocrServiceAddr,
		OcrServiceMode:    ocrServiceMode,

		DocSearchAddress: docSearchAddr,

		TokenizerServiceAddress: tokenizerServiceAddr,
		TokenizerServiceMode:    tokenizerServiceMode,
		TokenizerChunkSize:      int(chunkSize),
		TokenizerChunkOverlap:   int(chunkOverlap),
		TokenizerReturnChunks:   returnChunksFlag,
		TokenizerChunkBySelf:    chunkBySelfFlag,
		TokenizerTimeout:        uint(tokenizerTimeout),
	}, nil
}

func ParseServerAddress(serverAddr string) *server.Options {
	tmp := strings.Split(serverAddr, ":")
	servicePort, _ := strconv.Atoi(tmp[1])
	return server.BuildOptions(tmp[0], servicePort)
}
