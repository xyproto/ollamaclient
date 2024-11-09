package ollamaclient

import (
	"fmt"
	"strings"
	"testing"
)

func TestVisionImage(t *testing.T) {

	oc := New("llama3.2-vision")
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

	if found, err := oc.Has("llama3.2-vision"); err != nil || !found {
		t.Error("Expected to have 'llama3.2-vision' model downloaded, but it's not present")
	}

	oc.SetReproducible()

	imageFilename := "img/puppy.png"

	base64image, err := Base64EncodeFile(imageFilename)
	if err != nil {
		t.Fatalf("%s is missing or empty\n", imageFilename)
	}

	prompt := "How many puppy are there? Only numbers."
	generatedOutput, err := oc.GetOutputChatVision(prompt, base64image)
	if err != nil {
		t.Fatalf("Failed to get output: %v", err)
	}
	if len(generatedOutput.Content) == 0 {
		t.Fatalf("Generated output for the prompt %s is empty.\n", prompt)
	}
	if !strings.Contains(generatedOutput.Content, "1") {
		t.Fatalf("Generated output for the prompt %s does not contain '1'. Output: %s", prompt, generatedOutput.Content)
	}
}
