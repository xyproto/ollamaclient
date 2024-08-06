package ollamaclient

import (
	"encoding/json"
	"testing"
)

func TestTools(t *testing.T) {
	oc := New("llama3.1")
	oc.Verbose = true

	err := oc.PullIfNeeded(true)
	if err != nil {
		t.Fatalf("Failed to pull model: %v", err)
	}

	if found, err := oc.Has("llama3.1"); err != nil || !found {
		t.Error("Expected to have 'llama3.1' model downloaded, but it's not present")
	}

	oc.SetSystemPrompt("You are a helpful assistant.")
	oc.SetRandom()

	var toolGetCurrentWeather Tool
	json.Unmarshal(json.RawMessage(`{
		"type": "function",
		"function": {
		  "name": "get_current_weather",
		  "description": "Get the current weather for a location",
		  "parameters": {
			"type": "object",
			"properties": {
			  "location": {
				"type": "string",
				"description": "The location to get the weather for, e.g. San Francisco, CA"
			  },
			  "format": {
				"type": "string",
				"description": "The format to return the weather in, e.g. 'celsius' or 'fahrenheit'",
				"enum": ["celsius", "fahrenheit"]
			  }
			},
			"required": ["location", "format"]
		  }
		}
	  }`), &toolGetCurrentWeather)

	oc.SetTool(toolGetCurrentWeather)

	const prompt = "What is the weather in Toronto?"

	generatedOutput := oc.MustGetChatResponse(prompt)
	if generatedOutput.Error != "" {
		t.Error(generatedOutput.Error)
	}
	if len(generatedOutput.ToolCalls) != 1 {
		t.Errorf("Expected 1 tool call, got %d", len(generatedOutput.ToolCalls))
	}
	if generatedOutput.ToolCalls[0].Function.Name != "get_current_weather" {
		t.Errorf("Expected tool call 'get_current_weather', got '%s'", generatedOutput.ToolCalls[0].Function.Name)
	}
}
