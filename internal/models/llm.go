package models

import "fmt"

type SummaryRequest struct {
	NPredict         int             `json:"n_predict"`
	Temperature      float32         `json:"temperature"`
	Stop             []string        `json:"stop"`
	RepeatLastN      int             `json:"repeat_last_n"`
	RepeatPenalty    float32         `json:"repeat_penalty"`
	PenalizeNL       bool            `json:"penalize_nl"`
	TopK             int             `json:"top_k"`
	TopP             float32         `json:"top_p"`
	MinP             float32         `json:"min_p"`
	TFSz             int             `json:"tfs_z"`
	TypicalP         int             `json:"typical_p"`
	PresencePenalty  int             `json:"presence_penalty"`
	FrequencyPenalty int             `json:"frequency_penalty"`
	Mirostat         int             `json:"mirostat"`
	MirostatTAU      int             `json:"mirostat_tau"`
	MirostatETA      float32         `json:"mirostat_eta"`
	Grammar          string          `json:"grammar"`
	NProbs           int             `json:"n_probs"`
	MinKeep          int             `json:"min_keep"`
	RespFormat       *ResponseFormat `json:"response_format"`
	CachePROMPT      bool            `json:"cache_prompt"`
	APIKey           string          `json:"api_key"`
	SlotID           int             `json:"slot_id"`
	PROMPT           string          `json:"prompt"`
}

type ResponseFormat struct {
	Type string `json:"type"`
}

func NewLLM(content string) *SummaryRequest {
	insert := "```{\"summary\": \"summary of the content\", \"thematic\": \"determined class of document content\"}```"

	promtp := fmt.Sprintf(`
		You will be provided with the contents of a file along with its metadata. 
		Provide a summary of the contents. The purpose of the summary is to organize files based on their content. 
		To this end provide a concise but informative summary. Make the summary as specific to the file as possible. 
		And try determinate a document thematic to classify it by following:
			- military
 			- scientific
 			- other

		Write your response a JSON object with the following schema without any metadata or text:
		
		%s

		User: %s
		Llama:
	`, insert, content)

	return &SummaryRequest{
		NPredict:         400,
		Temperature:      0.1,
		Stop:             []string{"</s>", "Llama:", "User:"},
		RepeatLastN:      256,
		RepeatPenalty:    1.18,
		PenalizeNL:       false,
		TopK:             40,
		TopP:             0.95,
		MinP:             0.05,
		TFSz:             1,
		TypicalP:         1,
		PresencePenalty:  0,
		FrequencyPenalty: 0,
		Mirostat:         0,
		MirostatTAU:      5,
		MirostatETA:      0.1,
		Grammar:          "",
		NProbs:           0,
		MinKeep:          0,
		CachePROMPT:      false,
		APIKey:           "",
		SlotID:           -1,
		PROMPT:           promtp,
		RespFormat:       &ResponseFormat{Type: "json_object"},
	}
}
