package ollamaclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type VisionRequest struct {
	Prompt string   `json:"prompt"`
	Images []string `json:"images"`
}

func (oc *Config) GetOutputChatVision(promptAndOptionalImages ...string) (OutputChat, error) {
	var (
		temperature float64
		seed        = oc.SeedOrNegative
	)
	if len(promptAndOptionalImages) == 0 {
		return OutputChat{}, errors.New("at least one prompt must be given (and then optionally, base64 encoded JPG or PNG image strings)")
	}
	prompt := promptAndOptionalImages[0]
	var images []string
	if len(promptAndOptionalImages) > 1 {
		images = promptAndOptionalImages[1:]
	}
	if seed < 0 {
		temperature = oc.TemperatureIfNegativeSeed
	}
	messages := []Message{}
	if oc.SystemPrompt != "" {
		messages = append(messages, Message{
			Role:    "system",
			Content: oc.SystemPrompt,
		})
	}
	messages = append(messages, Message{
		Role:    "user",
		Content: prompt,
		Images:  images,
	})

	reqBody := GenerateChatRequest{
		Model:    oc.ModelName,
		Messages: messages,
		Tools:    oc.Tools,
		Options: RequestOptions{
			Seed:        seed,        // set to -1 to make it random
			Temperature: temperature, // set to 0 together with a specific seed to make output reproducible
		},
	}

	if oc.ContextLength != 0 {
		reqBody.Options.ContextLength = oc.ContextLength
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return OutputChat{}, err
	}
	if oc.Verbose {
		fmt.Printf("Sending request to %s/api/chat: %s\n", oc.ServerAddr, string(reqBytes))
	}
	HTTPClient := &http.Client{
		Timeout: oc.HTTPTimeout,
	}
	resp, err := HTTPClient.Post(oc.ServerAddr+"/api/chat", mimeJSON, bytes.NewBuffer(reqBytes))
	if err != nil {
		return OutputChat{}, err
	}
	defer resp.Body.Close()
	var res = OutputChat{}
	var sb strings.Builder
	decoder := json.NewDecoder(resp.Body)
	for {
		var genResp GenerateChatResponse
		if err := decoder.Decode(&genResp); err != nil {
			break
		}
		sb.WriteString(genResp.Message.Content)
		if genResp.Done {
			res.Role = genResp.Message.Role
			res.ToolCalls = genResp.Message.ToolCalls
			break
		}
	}
	res.Content = strings.TrimPrefix(sb.String(), "\n")
	if oc.TrimSpace {
		res.Content = strings.TrimSpace(res.Content)
	}
	return res, nil
}
