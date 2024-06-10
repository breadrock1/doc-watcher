package processing

import (
	"doc-notifier/internal/reader"
)

type ProcessJob struct {
	JobId    string           `json:"job_id"`
	Status   bool             `json:"status"`
	Document *reader.Document `json:"document"`
}

type Processor interface {
	GetProcessingJobs() map[string]*ProcessJob
	GetProcessingJob(jobId string) *ProcessJob
}
