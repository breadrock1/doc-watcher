package artifacts

type OcrArtifacts struct {
	AllDocTypes DocTypes `json:"doc_types"`
}

type DocTypes struct {
	TN         OcrDocType `json:"tn"`
	Smgs       OcrDocType `json:"smgs"`
	Conosament OcrDocType `json:"bill_of_landing"`
}

type OcrDocType struct {
	Name           string `json:"name"`
	JsonName       string `json:"json_name"`
	SampleFileName string `json:"sample_file_name"`
	Artifacts      []Arts `json:"artifacts"`
}

type Arts struct {
	GroupName     string `json:"name"`
	GroupJsonName string `json:"json_name"`
	Type          string `json:"type"`
}
