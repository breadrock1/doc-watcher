package config

import "time"

type Config struct {
	Logger    LoggerConfig
	Watcher   WatcherConfig
	Ocr       OcrConfig
	Searcher  SearcherConfig
	Tokenizer TokenizerConfig
	Storage   StorageConfig
	Office    OfficeConfig
	Minio     MinioConfig
	Samba     SambaConfig
}

type LoggerConfig struct {
	Level         string
	FilePath      string
	EnableFileLog bool
}

type WatcherConfig struct {
	Address            string
	WatchedDirectories []string
}

type OcrConfig struct {
	Address string
	Mode    string
	Timeout time.Duration
}

type SearcherConfig struct {
	Address string
	Timeout time.Duration
}

type TokenizerConfig struct {
	Address      string
	Mode         string
	ChunkSize    int
	ChunkOverlap int
	ReturnChunks bool
	ChunkBySelf  bool
	Timeout      time.Duration
}

type StorageConfig struct {
	DriverName string
	User       string
	Password   string
	Address    string
	Port       int
	DbName     string
	EnableSSL  string
	AddressLLM string
}

type OfficeConfig struct {
	Address    string
	WatcherDir string
}

type MinioConfig struct {
	Address           string
	MinioEndpoint     string
	BucketName        string
	MinioRootUser     string
	MinioRootPassword string
	MinioUseSSL       bool
}

type SambaConfig struct {
	Address  string
	Username string
	Password string
}
