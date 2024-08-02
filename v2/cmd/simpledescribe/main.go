package main

import (
	"fmt"
	"log"

	"github.com/xyproto/ollamaclient/v2"
)

func main() {
	oc := ollamaclient.New("llava")

	oc.SetReproducible()
	if err := oc.PullIfNeeded(true); err != nil {
		log.Fatalln(err)
	}
	imageFilenames := []string{"carrot1.png", "carrot2.png"}
	const desiredWordCount = 7
	description, err := oc.DescribeImages(imageFilenames, desiredWordCount)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(description)
}
