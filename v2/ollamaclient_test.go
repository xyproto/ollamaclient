package ollamaclient

import (
	"fmt"
	"testing"
)

func TestPullTinyLlamaIntegration(t *testing.T) {
	oc := New()
	oc.ModelName = "tinyllama"

	err := oc.PullIfNeeded(true)
	if err != nil {
		t.Fatalf("Failed to pull model: %v", err)
	}

	if found, err := oc.Has("tinyllama"); err != nil || !found {
		t.Error("Expected to have 'tinyllama' model downloaded, but it's not present")
	} else {
		t.Log("Pull operation completed successfully, 'tinyllama' model is present")
	}

	//oc.SetReproducible()
	oc.SetRandom()

	generatedOutput := oc.MustOutput("Generate an imperative sentence. Keep it brief. Only output the sentence itself. Skip explanations, introductions or preamble.")
	if generatedOutput == "" {
		t.Fail()
	}

	fmt.Println(Massage(generatedOutput))
}
