package models

type ProcessingDocuments struct {
	Done         []string `json:"done"`
	Processing   []string `json:"processing"`
	Unrecognized []string `json:"unrecognized"`
}
