package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/xyproto/ollamaclient"
)

func main() {
	oc := ollamaclient.New()

	//oc.Verbose = true

	filenames := []string{"../../README.md", "../../ollamaclient.go"}

	var sb strings.Builder
	for _, filename := range filenames {
		fmt.Printf("Reading %s...", filename)
		data, err := os.ReadFile(filename)
		if err != nil {
			fmt.Println("Skipping %s, got error: %s\n", filename, err)
			continue
		}
		fmt.Println("ok")
		sb.WriteString(filename + ":\n")
		sb.Write(data)
	}

	prompt := "Write a short summary of what a project that contains the following files is:\n\n" + sb.String()

	fmt.Printf("Sending request to Ollama, using the %s model...\n", oc.Model)

	output, err := oc.GetOutput(prompt)
	if err != nil {
		fmt.Println("Error:", err)
		return

	}
	fmt.Printf("\n%s\n", strings.TrimSpace(output))
}
