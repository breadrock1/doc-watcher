package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func FromFile(filePath string) (*Config, error) {
	config := &Config{}

	viperInstance := viper.New()
	viperInstance.AutomaticEnv()
	viperInstance.SetConfigFile(filePath)

	viperInstance.SetDefault("logger.Level", "INFO")
	viperInstance.SetDefault("logger.FilePath", "./logs/app.log")
	viperInstance.SetDefault("logger.EnableFileLog", false)

	viperInstance.SetDefault("watcher.Address", "http://0.0.0.0:2893")
	viperInstance.SetDefault("watcher.WatchedDirectories", []string{"./indexer"})

	viperInstance.SetDefault("ocr.Address", "http://localhost:8004")
	viperInstance.SetDefault("ocr.Mode", "assistant")

	viperInstance.SetDefault("searcher.Address", "http://localhost:2892")

	viperInstance.SetDefault("tokenizer.Address", "http://localhost:8001")
	viperInstance.SetDefault("tokenizer.Mode", "none")
	viperInstance.SetDefault("tokenizer.ChunkSize", 800)
	viperInstance.SetDefault("tokenizer.ChunkOverlap", 100)
	viperInstance.SetDefault("tokenizer.ReturnChunks", false)
	viperInstance.SetDefault("tokenizer.ChunkBySelf", false)
	viperInstance.SetDefault("tokenizer.Timeout", 300)

	viperInstance.SetDefault("office.Address", "http://localhost:8087")

	viperInstance.SetDefault("minio.Address", "0.0.0.0:2894")
	viperInstance.SetDefault("minio.MinioEndpoint", "localhost:9000")
	viperInstance.SetDefault("minio.BucketName", "indexer")
	viperInstance.SetDefault("minio.MinioRootUser", "<access-id>")
	viperInstance.SetDefault("minio.MinioRootPassword", "<secret-key>")
	viperInstance.SetDefault("minio.MinioUseSSL", false)

	viperInstance.SetDefault("samba.Address", "0.0.0.0:445")
	viperInstance.SetDefault("samba.Username", "admin")
	viperInstance.SetDefault("samba.Password", "admin")

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
