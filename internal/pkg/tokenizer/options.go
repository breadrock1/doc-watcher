package tokenizer

import "time"

type Mode int

const (
	Assistant Mode = iota
	LangChain
	None
)

type Options struct {
	Mode         Mode
	Address      string
	Timeout      time.Duration
	ChunkedFlag  bool
	ChunkSize    int
	ChunkOverlap int
}

func GetModeFromString(mode string) Mode {
	switch mode {
	case "assistant":
		return Assistant
	case "langchain":
		return LangChain
	case "none":
		return None
	default:
		return None
	}
}
