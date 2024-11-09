package main

import (
	"fmt"

	"github.com/xyproto/ollamaclient/v2"
	"github.com/xyproto/usermodel"
)

func main() {
	oc := ollamaclient.New(usermodel.GetChatModel())
	if err := oc.PullIfNeeded(true); err != nil {
		fmt.Println("Error:", err)
		return
	}
	prompt := "Write a haiku about the color of cows."
	output, err := oc.GetOutput(prompt)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(output)
}
