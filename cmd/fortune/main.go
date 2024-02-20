package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xyproto/ollamaclient"
)

const (
	model  = "tinydolphin"
	prompt = "Write a silly saying, quote or joke like it could have been the output of the fortune command on Linux."
)

func main() {
	oc := ollamaclient.NewWithModel(model)

	err := oc.PullIfNeeded(true)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to pull model: %v", err)
		os.Exit(1)
	}

	if !oc.Has(model) {
		fmt.Fprintf(os.Stderr, "Expected to have '%s' model downloaded, but it's not present\n", model)
		os.Exit(1)
	}

	oc.SetRandomOutput()

	generatedOutput := oc.MustOutput(prompt)
	if generatedOutput == "" {
		log.Println("fail")
	}

	fmt.Println(ollamaclient.Massage(generatedOutput))
}
