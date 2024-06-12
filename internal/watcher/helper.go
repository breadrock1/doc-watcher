package watcher

import (
	"bytes"
	"context"
	"doc-notifier/internal/reader"
	"doc-notifier/internal/sender"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

func (nw *NotifyWatcher) recognizeTriggeredDoc(documents []*reader.Document) {
	wg := sync.WaitGroup{}
	for _, document := range documents {
		wg.Add(1)
		document := document
		go func() {
			defer wg.Done()
			nw.recognizeCallback(document)
			<-time.After(2 * time.Second)
		}()
	}

	wg.Wait()
}

func (nw *NotifyWatcher) recognizeCallback(document *reader.Document) {
	document.SetQuality(0)
	if err := nw.Ocr.Ocr.RecognizeFile(document); err != nil {
		log.Println(err)
		return
	}

	document.ComputeMd5Hash()
	document.ComputeSsdeepHash()
	document.SetEmbeddings([]*reader.Embeddings{})

	log.Println("Computing tokens for extracted text: ", document.DocumentName)
	tokenVectors, _ := nw.Tokenizer.Tokenizer.TokenizeTextData(document.Content)
	for chunkID, chunkData := range tokenVectors.Vectors {
		text := tokenVectors.ChunkedText[chunkID]
		document.AppendContentVector(text, chunkData)
	}

	log.Println("Storing document to searcher: ", document.DocumentName)
	if err := nw.Searcher.StoreDocument(document); err != nil {
		log.Println("Failed while storing document: ", err)
	}

	ctx := context.Background()
	nw.loadSummary(document)
	if _, err := nw.Storage.Create(ctx, document); err != nil {
		log.Println("Failed while storing metadata to psql: ", err)
	}
}

type SummaryWrapper struct {
	Content string `json:"content"`
}

type SummaryResponse struct {
	FilePath string `json:"file_path"`
	Summary  string `json:"summary"`
	Class    string `json:"thematic"`
}

type ResponseFormat struct {
	Type string `json:"type"`
}

type SummaryRequest struct {
	//Stream           bool     `json:"stream"`
	NPredict         int      `json:"n_predict"`
	Temperature      float32  `json:"temperature"`
	Stop             []string `json:"stop"`
	RepeatLastN      int      `json:"repeat_last_n"`
	RepeatPenalty    float32  `json:"repeat_penalty"`
	PenalizeNL       bool     `json:"penalize_nl"`
	TopK             int      `json:"top_k"`
	TopP             float32  `json:"top_p"`
	MinP             float32  `json:"min_p"`
	TFSz             int      `json:"tfs_z"`
	TypicalP         int      `json:"typical_p"`
	PresencePenalty  int      `json:"presence_penalty"`
	FrequencyPenalty int      `json:"frequency_penalty"`
	Mirostat         int      `json:"mirostat"`
	MirostatTAU      int      `json:"mirostat_tau"`
	MirostatETA      float32  `json:"mirostat_eta"`
	Grammar          string   `json:"grammar"`
	NProbs           int      `json:"n_probs"`
	MinKeep          int      `json:"min_keep"`
	//ImageData        []byte   `json:"image_data"`
	RespFormat  *ResponseFormat `json:"response_format"`
	CachePROMPT bool            `json:"cache_prompt"`
	APIKey      string          `json:"api_key"`
	SlotID      int             `json:"slot_id"`
	PROMPT      string          `json:"prompt"`
}

func (nw *NotifyWatcher) loadSummary(document *reader.Document) {
	insert := "```{\"summary\": \"summary of the content\", \"thematic\": \"determined class of document content\"}```"

	text := fmt.Sprintf(`
		You will be provided with the contents of a file along with its metadata. 
		Provide a summary of the contents. The purpose of the summary is to organize files based on their content. 
		To this end provide a concise but informative summary. Make the summary as specific to the file as possible. 
		And try determinate a document thematic to classify it by following:
			- military
 			- scientific
 			- other

		Write your response a JSON object with the following schema without any metadata or text:
		
		%s

		%s
	`, insert, document.Content)

	summaryRequest := &SummaryRequest{
		//Stream:           false,
		NPredict:         400,
		Temperature:      0.7,
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
		//ImageData:        make([]byte, 0),
		CachePROMPT: false,
		APIKey:      "",
		SlotID:      -1,
		PROMPT:      text,
		RespFormat:  &ResponseFormat{Type: "json_object"},
	}

	jsonData, err := json.Marshal(summaryRequest)
	if err != nil {
		log.Println("Failed while marshaling doc: ", err)
		return
	}

	reqBody := bytes.NewBuffer(jsonData)

	method := "POST"
	targetURL := "http://192.168.0.59:8081/completion"
	mimeType := "application/json"
	respData, recErr := sender.SendRequest(reqBody, &targetURL, &method, &mimeType, 300*time.Second)
	if recErr != nil {
		log.Println("failed send request: ", recErr)
		return
	}

	var summaryResponse *SummaryWrapper
	if err := json.Unmarshal(respData, &summaryResponse); err != nil {
		log.Println("Failed while reading response reqBody: ", err)
		return
	}

	//document.Content = summaryResponse.Content.Summary
	//document.DocumentClass = summaryResponse.Content.Class
}
