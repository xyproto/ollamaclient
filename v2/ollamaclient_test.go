package ollamaclient

import (
	"fmt"
	"testing"
)

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

	err := oc.PullIfNeeded(true)
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
