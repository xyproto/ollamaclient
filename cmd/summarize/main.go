package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/xyproto/ollamaclient"
	"github.com/xyproto/wordwrap"
	"golang.org/x/term"
)

const defaultModel = "nous-hermes:latest"

var verbose bool

// Only print the provided data when in verbose mode
func logVerbose(format string, a ...interface{}) {
	if verbose {
		fmt.Printf(format, a...)
	}
}

func main() {
	// Flags
	var promptHeader, outputFile, model string
	var wrapWidth int

	// Get the current terminal width as default
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		// Default to 79 if unable to get terminal size
		width = 79
	}

	pflag.BoolVarP(&verbose, "verbose", "V", false, "verbose output")
	pflag.StringVarP(&promptHeader, "prompt", "p", "Write a short summary of what a project that contains the following files is:", "Provide a custom prompt header")
	pflag.StringVarP(&outputFile, "output", "o", "", "Specify an output file")
	pflag.StringVarP(&model, "model", "m", defaultModel, "Specify the Ollama model to use")
	pflag.IntVarP(&wrapWidth, "wrap", "w", width, "Word wrap at specified width")
	pflag.Parse()

	// Retrieve non-flag arguments
	filenames := pflag.Args()
	if len(filenames) < 1 {
		fmt.Println("Usage: summarize [--prompt <customPrompt>] [--output <outputFile>] [--wrap <width>] [--model <ollamaModel>] <filename1> [<filename2> ...]")
		fmt.Println("Error: Please provide at least one filename.")
		os.Exit(1)
	}

	// Build a prompt by reading in all given filenames
	var sb strings.Builder
	readCount := 0
	for _, filename := range filenames {
		logVerbose("[%s] Reading... ", filename)
		data, err := os.ReadFile(filename)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			continue
		}
		readCount++
		logVerbose("OK\n")
		sb.WriteString(filename + ":\n")
		sb.Write(data)
	}

	if readCount == 0 {
		fmt.Println("Error: No files could be read.")
		os.Exit(1)
	}

	prompt := promptHeader + "\n\n" + sb.String()

	// Generate text with Ollama
	oc := ollamaclient.New()
	if model != defaultModel {
		oc.Model = model
	}
	logVerbose("[%s] Generating... ", oc.Model)
	output, err := oc.GetOutput(prompt)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

	// Output the result
	logVerbose("OK\n")
	trimmedOutput := strings.TrimSpace(output)

	if outputFile != "" {
		err := os.WriteFile(outputFile, []byte(trimmedOutput), 0o644)
		if err != nil {
			fmt.Printf("error writing to file: %s\n", err)
			os.Exit(1)
		}
		return
	}

	lines, err := wordwrap.WordWrap(trimmedOutput, wrapWidth)
	if err != nil {
		fmt.Println(trimmedOutput)
	}
	for _, line := range lines {
		fmt.Println(line)
	}
}
