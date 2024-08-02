package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xyproto/ollamaclient/v2"
)

const (
	model  = "tinydolphin"
	prompt = "Write a silly saying, quote or joke like it could have been the output of the fortune command on Linux."
)

func main() {
	oc := ollamaclient.New(model)

	err := oc.PullIfNeeded(true)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to pull model: %v\n", err)
		os.Exit(1)
	}

	if found, err := oc.Has(model); err != nil || !found {
		fmt.Fprintf(os.Stderr, "Expected to have '%s' model downloaded, but it's not present\n", model)
		os.Exit(1)
	}

	oc.SetRandom()

	generatedOutput := oc.MustOutput(prompt)
	if generatedOutput == "" {
		log.Println("Could not generate output.")
	}

	fmt.Println(ollamaclient.Massage(generatedOutput))
}
