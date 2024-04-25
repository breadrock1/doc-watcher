package processing

import "doc-notifier/internal/pkg/reader"

type ProcessJob struct {
	JobId    string           `json:"job_id"`
	Status   bool             `json:"status"`
	Document *reader.Document `json:"document"`
}
