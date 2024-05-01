package ollamaclient

import (
	"fmt"
	"testing"
)

const prompt = "Write a haiku about nine llamas."

func TestStream(t *testing.T) {
	oc := New()
	oc.ModelName = "tinyllama"
	oc.Verbose = true

	err := oc.PullIfNeeded(true)
	if err != nil {
		t.Fatalf("Failed to pull model: %v", err)
	}

	callbackFunction := func(partialResult string, streamingDone bool) {
		if !streamingDone {
			fmt.Printf("%s\n", partialResult)
		} else {
			fmt.Println("DONE!")
		}
	}

	if err := oc.StreamOutput(callbackFunction, prompt); err != nil {
		t.Fatalf("Failed to get streamed output: %v", err)
	}
}
