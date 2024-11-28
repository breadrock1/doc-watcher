package httpserv

func createStatusResponse(status int, msg string) *ResponseForm {
	return &ResponseForm{Status: status, Message: msg}
}

// ResponseForm example
type ResponseForm struct {
	Status  int    `json:"status" example:"200"`
	Message string `json:"message" example:"Done"`
}

// BadRequestForm example
type BadRequestForm struct {
	Status  int    `json:"status" example:"400"`
	Message string `json:"message" example:"Bad Request message"`
}

// ServerErrorForm example
type ServerErrorForm struct {
	Status  int    `json:"status" example:"503"`
	Message string `json:"message" example:"Server Error message"`
}

// AttachDirectoryForm example
type AttachDirectoryForm struct {
	BucketName string `json:"bucket_name" example:"test-folder"`
}

// FetchDocumentsList example
type FetchDocumentsList struct {
	FileNames []string `json:"file_names" example:"test-file.docx"`
}
