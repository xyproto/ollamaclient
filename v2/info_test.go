package ollamaclient

import "testing"

func TestGetShowInfo(t *testing.T) {
	oc := New("llama3.2")
	oc.Verbose = true

	err := oc.PullIfNeeded(true)
	if err != nil {
		t.Fatalf("Failed to pull model: %v", err)
	}

	if found, err := oc.Has("llama3.2"); err != nil || !found {
		t.Error("Expected to have 'llama3.2' model downloaded, but it's not present")
	}

	res, err := oc.GetShowInfo()
	if err != nil {
		t.Fatalf("Failed to get show info: %v", err)
	}

	if res.ModelInfo.LlamaContextLength != 131072 {
		t.Errorf("Expected context length to be 131072, but got %d", res.ModelInfo.LlamaContextLength)
	}
}
