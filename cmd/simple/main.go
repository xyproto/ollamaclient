package main

import (
	"fmt"

	"github.com/xyproto/ollamaclient"
)

func main() {
	oc := ollamaclient.NewWithModel("tinyllama")

	if err := oc.PullIfNeeded(true); err != nil {
		fmt.Println("Error:", err)
		return
	}

	prompt := "Write a haiku about the color of cows."
	output, err := oc.GetOutput(prompt, true)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(output)
}
