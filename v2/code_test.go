package ollamaclient

import (
	"fmt"
	"testing"
)

const codeCompleteModel = "deepseek-coder-v2:latest"

func TestBetween(t *testing.T) {
	const (
		codeStart = "def compute_gcd(a, b):"
		codeEnd   = "    return result"
	)
	oc := New(codeCompleteModel)
	oc.Verbose = true
	if err := oc.PullIfNeeded(true); err != nil {
		t.Fatalf("Failed to pull model: %v", err)
	}
	response, err := oc.GetBetweenResponse(codeStart, codeEnd)
	if err != nil {
		t.Fatalf("Failed to get code completion: %v", err)
	}
	fmt.Printf("%s\n%s\n%s\n", codeStart, response.Response, codeEnd)
}

func TestCodeCompletion(t *testing.T) {
	const (
		codeStart = "def compute_gcd(a, b):"
		codeEnd   = "    return result"
		verbose   = true
	)
	oc := New(codeCompleteModel)
	oc.Verbose = true
	codeBetween, err := oc.Complete(codeStart, codeEnd)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s\n%s\n%s\n", codeStart, codeBetween, codeEnd)
}
