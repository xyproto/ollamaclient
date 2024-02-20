package ollamaclient

import (
	"fmt"
	"strings"
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

	generatedOutput := oc.MustOutput("Generate an imperative sentence. Keep it brief. Only output the sentence itself. Skip explanations, introductions or preamble.")
	if generatedOutput == "" {
		t.Fail()
	}
	s := generatedOutput

	// Keep the part after ":", if applicable
	if strings.Contains(s, ":") {
		parts := strings.SplitN(s, ":", 2)
		s = strings.TrimSpace(parts[1])
	}
	// Keep the part within double quotes, if applicable
	if strings.Count(s, "\"") == 2 {
		parts := strings.SplitN(s, "\"", 3)
		s = strings.TrimSpace(parts[1])
	}
	// Keep the part after ":", if applicable
	if strings.Contains(s, ":") {
		parts := strings.SplitN(s, ":", 2)
		s = strings.TrimSpace(parts[1])
	}
	// Remove stray quotes
	s = strings.TrimPrefix(s, "\"")
	s = strings.TrimPrefix(s, "'")
	s = strings.TrimSuffix(s, "\"")
	s = strings.TrimSuffix(s, "'")

	// Keep the last line
	if strings.Count(s, "\n") > 1 {
		lines := strings.Split(s, "\n")
		s = lines[len(lines)-1]
	}

	// Keep the part before the exclamation mark, if applicable
	if strings.Contains(s, "!") {
		parts := strings.SplitN(s, "!", 2)
		s = strings.TrimSpace(parts[0]) + "!"
	}

	if s == "" {
		s = generatedOutput
	}

	fmt.Println(s)
}
