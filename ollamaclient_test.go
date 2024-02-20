package ollamaclient

import (
	"testing"
)

func TestPullTinyLlamaIntegration(t *testing.T) {
	oc := NewWithModel("tinyllama")

	_, err := oc.Pull(true)
	if err != nil {
		t.Fatalf("Failed to pull model: %v", err)
	}

	if !oc.Has("tinyllama") {
		t.Error("Expected to have 'tinyllama' model downloaded, but it's not present")
	} else {
		t.Log("Pull operation completed successfully, 'tinyllama' model is present")
	}
}
