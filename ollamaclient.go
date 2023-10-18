// Package ollamaclient can be used for communicating with the Ollama service
package ollamaclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/xyproto/env/v2"
)

const defaultModel = "nous-hermes:7b-llama2-q2_K"

// Config represents configuration details for communicating with the Ollama API
type Config struct {
	API     string
	Model   string
	Verbose bool
}

// GenerateRequest represents the request payload for generating output
type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// GenerateResponse represents the response data from the generate API call
type GenerateResponse struct {
	Model              string `json:"model"`
	CreatedAt          string `json:"created_at"`
	Response           string `json:"response"`
	Done               bool   `json:"done"`
	Context            []int  `json:"context,omitempty"`
	TotalDuration      int64  `json:"total_duration,omitempty"`
	LoadDuration       int64  `json:"load_duration,omitempty"`
	SampleCount        int    `json:"sample_count,omitempty"`
	SampleDuration     int64  `json:"sample_duration,omitempty"`
	PromptEvalCount    int    `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration int64  `json:"prompt_eval_duration,omitempty"`
	EvalCount          int    `json:"eval_count,omitempty"`
	EvalDuration       int64  `json:"eval_duration,omitempty"`
}

// EmbeddingsRequest represents the request payload for getting embeddings
type EmbeddingsRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// EmbeddingsResponse represents the response data containing embeddings
type EmbeddingsResponse struct {
	Embeddings []float64 `json:"embedding"`
}

// PullRequest represents the request payload for pulling a model
type PullRequest struct {
	Name     string `json:"name"`
	Insecure bool   `json:"insecure,omitempty"`
	Stream   bool   `json:"stream,omitempty"`
}

// PullResponse represents the response data from the pull API call
type PullResponse struct {
	Status string `json:"status"`
	Digest string `json:"digest"`
	Total  int64  `json:"total"`
}

// New initializes a new Config using environment variables
func New() *Config {
	return &Config{
		env.Str("OLLAMA_HOST", "http://localhost:11434"),
		env.Str("OLLAMA_MODEL", defaultModel),
		env.Bool("OLLAMA_VERBOSE"),
	}
}

// NewWithModel initializes a new Config using a specified model and environment variables
func NewWithModel(model string) *Config {
	return &Config{
		env.Str("OLLAMA_HOST", "http://localhost:11434"),
		model,
		env.Bool("OLLAMA_VERBOSE"),
	}
}

// GetOutput sends a request to the Ollama API and returns the generated output
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
		fmt.Printf("Sending request to /api/generate: %s\n", string(reqBytes))
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

// AddEmbedding sends a request to get embeddings for a given prompt
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
		fmt.Printf("Sending request to /api/embeddings: %s\n", string(reqBytes))
	}

	resp, err := http.Post(c.API+"/api/embeddings", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return []float64{}, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var embResp EmbeddingsResponse
	if err := decoder.Decode(&embResp); err != nil {
		return []float64{}, err
	}
	return embResp.Embeddings, nil
}

// Pull sends a request to pull a specified model from the Ollama API
func (c *Config) Pull() (string, error) {
	reqBody := PullRequest{
		Name:   c.Model,
		Stream: false,
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	if c.Verbose {
		fmt.Printf("Sending request to /api/pull: %s\n", string(reqBytes))
	}
	resp, err := http.Post(c.API+"/api/pull", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var pullResp PullResponse
	if err := decoder.Decode(&pullResp); err != nil {
		return "", err
	}
	return pullResp.Status, nil
}
