package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const defaultAPI = "http://localhost:11434"
//const defaultModel = "codeup:latest"
const defaultModel = "nous-hermes:latest"

type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type GenerateResponse struct {
	Model              string `json:"model"`
	CreatedAt          string `json:"created_at"`
	Response           string `json:"response"`
	Done               bool   `json:"done"`
	Context            []int  `json:"context,omitempty"`              // optional field
	TotalDuration      int64  `json:"total_duration,omitempty"`       // optional field
	LoadDuration       int64  `json:"load_duration,omitempty"`        // optional field
	SampleCount        int    `json:"sample_count,omitempty"`         // optional field
	SampleDuration     int64  `json:"sample_duration,omitempty"`      // optional field
	PromptEvalCount    int    `json:"prompt_eval_count,omitempty"`    // optional field
	PromptEvalDuration int64  `json:"prompt_eval_duration,omitempty"` // optional field
	EvalCount          int    `json:"eval_count,omitempty"`           // optional field
	EvalDuration       int64  `json:"eval_duration,omitempty"`        // optional field
}

type EmbeddingsRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type EmbeddingsResponse struct {
	Embeddings []float64 `json:"embedding"`
}

func getOutput(api, prompt string) (string, error) {
	reqBody := GenerateRequest{
		Model:  defaultModel,
		Prompt: prompt,
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	resp, err := http.Post(api+"/api/generate", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var sb strings.Builder
	decoder := json.NewDecoder(resp.Body)
	for {
		var genResp GenerateResponse
		if err := decoder.Decode(&genResp); err != nil {
			break
		}
		sb.WriteString(genResp.Response)
		if genResp.Done {
			break
		}
	}
	return sb.String(), nil
}

func addEmbedding(api, prompt string) (string, error) {
	reqBody := EmbeddingsRequest{
		Model:  defaultModel,
		Prompt: prompt,
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// Log the request for debugging
	fmt.Printf("Sending request: %s\n", string(reqBytes))

	resp, err := http.Post(api+"/api/embeddings", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var sb strings.Builder
	decoder := json.NewDecoder(resp.Body)
	var embResp EmbeddingsResponse
	if err := decoder.Decode(&embResp); err != nil {
		return "", err
	}
	sb.WriteString(fmt.Sprintf("%v", embResp.Embeddings))
	return sb.String(), nil
}

func main() {
	api := defaultAPI
	if host := os.Getenv("OLLAMA_HOST"); host != "" {
		api = host
	}

	//fmt.Print("Add an embedding by prompt: ")
	//var input1 string
	//fmt.Scanln(&input1)
	input1 := "All cangaroos are red."
	output, err := addEmbedding(api, input1)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Got:", output)

	//fmt.Print("Make a request by prompt: ")
	//var input2 string
	//fmt.Scanln(&input2)
	input2 := "Which color does the cangaroos have?"
	output, err = getOutput(api, input2)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Got:", output)
}
