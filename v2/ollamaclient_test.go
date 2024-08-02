package ollamaclient

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestTools(t *testing.T) {
	oc := New()
	oc.ModelName = "llama3.1"
	oc.Verbose = true

	err := oc.PullIfNeeded(true)
	if err != nil {
		t.Fatalf("Failed to pull model: %v", err)
	}

	if found, err := oc.Has("llama3.1"); err != nil || !found {
		t.Error("Expected to have 'llama3.1' model downloaded, but it's not present")
	}

	oc.SetRandom()
	oc.SetTool(json.RawMessage(`{
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
	  }`))

	prompt := "What is the weather in Toronto?"
	generatedOutput := oc.MustOutputChat(prompt)
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

func TestPullTinyLlamaIntegration(t *testing.T) {
	oc := New()
	oc.ModelName = "tinyllama"
	oc.Verbose = true

	err := oc.PullIfNeeded(true)
	if err != nil {
		t.Fatalf("Failed to pull model: %v", err)
	}

	if found, err := oc.Has("tinyllama"); err != nil || !found {
		t.Error("Expected to have 'tinyllama' model downloaded, but it's not present")
	}

	oc.SetRandom()

	prompt := "Generate an imperative sentence. Keep it brief. Only output the sentence itself. Skip explanations, introductions or preamble."
	generatedOutput := oc.MustOutput(prompt)
	if generatedOutput == "" {
		t.Fatalf("Generated output for the prompt %s is empty.\n", prompt)
	}
	fmt.Println(Massage(generatedOutput))
}

func TestDescribeImage(t *testing.T) {

	oc := New()
	oc.ModelName = "llava"
	oc.Verbose = true

	ver, err := oc.Version()
	if err != nil {
		t.Fatalf("Failed to fetch the current Ollama version: %v", err)
	}
	fmt.Printf("Ollama version %s\n", ver)

	err = oc.PullIfNeeded(true)
	if err != nil {
		t.Fatalf("Failed to pull model: %v", err)
	}

	if found, err := oc.Has("llava"); err != nil || !found {
		t.Error("Expected to have 'llava' model downloaded, but it's not present")
	}

	oc.SetReproducible()

	imageFilename := "img/puppy.png"

	base64image, err := Base64EncodeFile(imageFilename)
	if err != nil {
		t.Fatalf("%s is missing or empty\n", imageFilename)
	}

	prompt := "Describe this image:"
	generatedOutput := oc.MustOutput(prompt, base64image)
	if generatedOutput == "" {
		t.Fatalf("Generated output for the prompt %s is empty.\n", prompt)
	}
	fmt.Println(Massage(generatedOutput))
}
