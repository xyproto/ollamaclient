package ollamaclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/xyproto/env/v2"
)

type Config struct {
	API     string
	Model   string
	Verbose bool
}

func New() *Config {
	return &Config{
		env.Str("OLLAMA_HOST", "http://localhost:11434"),
		env.Str("OLLAMA_MODEL", "nous-hermes:latest"),
		env.Bool("OLLAMA_VERBOSE"),
	}
}

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

func (c *Config) GetOutput(prompt string) (string, error) {
	reqBody := GenerateRequest{
		Model:  c.Model,
		Prompt: prompt,
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	if c.Verbose {
		fmt.Printf("Sending request: %s\n", string(reqBytes))
	}
	resp, err := http.Post(c.API+"/api/generate", "application/json", bytes.NewBuffer(reqBytes))
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

func (c *Config) AddEmbedding(prompt string) ([]float64, error) {
	reqBody := EmbeddingsRequest{
		Model:  c.Model,
		Prompt: prompt,
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return []float64{}, err
	}

	if c.Verbose {
		fmt.Printf("Sending request: %s\n", string(reqBytes))
	}

	resp, err := http.Post(c.API+"/api/embeddings", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return []float64{}, err
	}
	defer resp.Body.Close()
	//var sb strings.Builder
	decoder := json.NewDecoder(resp.Body)
	var embResp EmbeddingsResponse
	if err := decoder.Decode(&embResp); err != nil {
		return []float64{}, err
	}
	//sb.WriteString(fmt.Sprintf("%v", embResp.Embeddings))
	return embResp.Embeddings, nil
}
