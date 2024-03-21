package ocr

type OcrMode int

const (
	ReadRawFile OcrMode = iota
	LocalTesseract
	RemoteTesseract
	RemoteDedoc
)

type Options struct {
	Mode    OcrMode
	Address string
}

func GetModeFromString(mode string) OcrMode {
	switch mode {
	case "read-raw-file":
		return ReadRawFile
	case "local-tesseract":
		return LocalTesseract
	case "remote-tesseract":
		return RemoteTesseract
	case "dedoc":
		return RemoteDedoc
	default:
		return ReadRawFile
	}
}
