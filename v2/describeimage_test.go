package ollamaclient

import (
	"fmt"
	"testing"
)

func TestDescribeImage(t *testing.T) {

	oc := New("llava")
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
	if generatedOutput.Response == "" {
		t.Fatalf("Generated output for the prompt %s is empty.\n", prompt)
	}
	fmt.Println(Massage(generatedOutput.Response))
}
