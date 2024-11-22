package watcher

import "time"

type Config struct {
	Address            string
	Username           string
	Password           string
	EnableSSL          bool
	WatchedDirectories []string
	CacheExpire        time.Duration
	CacheCleanInterval time.Duration
}
