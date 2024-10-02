package main

import (
	"fmt"
	"log"

	"github.com/xyproto/ollamaclient/v2"
	"github.com/xyproto/usermodel"
)

func main() {
	oc := ollamaclient.New(usermodel.GetVisionModel())
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
