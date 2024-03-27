package ocr

type OcrMode int

const (
	ReadRawFile OcrMode = iota
	AssistantMode
	DedocWrapper
)

type Options struct {
	Mode    OcrMode
	Address string
}

func GetModeFromString(mode string) OcrMode {
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
