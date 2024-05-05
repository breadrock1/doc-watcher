package logoper

type OcrJobErrorType int

const (
	Processing OcrJobErrorType = iota
	FailedResponse
)

type OcrJobError struct {
	Type    OcrJobErrorType
	Message string
}

type OcrJob struct {
	JobId string `json:"job_id"`
}
