package ollamaclient

import (
	"fmt"
	"testing"
)

func TestPullTinyLlamaIntegration(t *testing.T) {
	oc := NewWithModel("tinyllama")

	err := oc.PullIfNeeded(true)
	if err != nil {
		t.Fatalf("Failed to pull model: %v", err)
	}

	if !oc.Has("tinyllama") {
		t.Error("Expected to have 'tinyllama' model downloaded, but it's not present")
	} else {
		t.Log("Pull operation completed successfully, 'tinyllama' model is present")
	}

	//oc.SetReproducibleOutput()
	oc.SetRandomOutput()

	generatedOutput := oc.MustOutput("Generate an imperative sentence. Keep it brief. Only output the sentence itself. Skip explanations, introductions or preamble.")
	if generatedOutput == "" {
		t.Fail()
	}

	fmt.Println(Massage(generatedOutput))
}
