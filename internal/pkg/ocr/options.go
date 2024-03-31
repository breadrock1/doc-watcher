package ocr

import "time"

type Mode int

const (
	ReadRawFile Mode = iota
	AssistantMode
	DedocWrapper
)

type Options struct {
	Mode    Mode
	Address string
	Timeout time.Duration
}

func GetModeFromString(mode string) Mode {
	switch mode {
	case "read-raw-file":
		return ReadRawFile
	case "assistant":
		return AssistantMode
	case "dedoc-wrapper":
		return DedocWrapper
	default:
		return ReadRawFile
	}
}
