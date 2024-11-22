package ocr

import "time"

type Config struct {
	Address   string
	EnableSSL bool
	Timeout   time.Duration
}
