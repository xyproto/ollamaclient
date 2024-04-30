package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/xyproto/ollamaclient/v2"
	"github.com/xyproto/wordwrap"
	"golang.org/x/term"
)

const (
	versionString    = "DescribeImage 1.0.0"
	defaultModel     = "llava"
	defaultTermWidth = 79
)

var verbose bool

// Only print the provided data when in verbose mode
func logVerbose(format string, a ...interface{}) {
	if verbose {
		fmt.Printf(format, a...)
	}
}

func getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return defaultTermWidth
	}
	return width
}

func main() {
	var (
		promptHeader, outputFile, model string
		wrapWidth                       int
		showVersion                     bool
	)

	pflag.BoolVarP(&verbose, "verbose", "V", false, "verbose output")
	pflag.StringVarP(&promptHeader, "prompt", "p", "Write a short summary of what a project that contains the following files is:", "Provide a custom prompt header")
	pflag.StringVarP(&outputFile, "output", "o", "", "Specify an output file")
	pflag.StringVarP(&model, "model", "m", defaultModel, "Specify the Ollama model to use")
	pflag.IntVarP(&wrapWidth, "wrap", "w", 0, "Word wrap at specified width. Use '-1' for terminal width")
	pflag.BoolVarP(&showVersion, "version", "v", false, "display version")

	pflag.Parse()

	if showVersion {
		fmt.Println(versionString)
		return
	}

	if wrapWidth == -1 {
		wrapWidth = getTerminalWidth()
	}

	filenames := pflag.Args()
	if len(filenames) < 1 {
		fmt.Println("Usage: describeimage [--prompt <customPrompt>] [--output <outputFile>] [--wrap <width>|-1] [--model <ollamaModel>] <image_filename1> [<image_filename2> ...]")
		os.Exit(1)
	}

	var images []string
	for _, filename := range filenames {
		logVerbose("[%s] Reading... ", filename)
		base64image, err := ollamaclient.Base64EncodeFile(filename)
		if err == nil { // success
			images = append(images, base64image)
			logVerbose("OK\n")
		} else {
			logVerbose("FAILED: " + err.Error() + "\n")
		}
	}

	var prompt string
	switch len(images) {
	case 0:
		fmt.Println("Error: no images to describe")
		os.Exit(1)
	case 1:
		prompt = "Describe this image:"
	default:
		prompt = "Describe these images:"
	}
	if promptHeader != "" {
		prompt = promptHeader
	}

	oc := ollamaclient.New()
	oc.ModelName = model

	if err := oc.PullIfNeeded(verbose); err != nil {
		fmt.Println("Error:", err)
		return
	}
	oc.SetReproducible()

	promptAndImages := append([]string{prompt}, images...)

	logVerbose("[%s] Generating... ", oc.ModelName)
	output, err := oc.GetOutput(promptAndImages...)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
	logVerbose("OK\n")

	if output == "" {
		fmt.Printf("Generated output for the prompt %s is empty.\n", prompt)
		os.Exit(1)
	}

	if wrapWidth > 0 {
		lines, err := wordwrap.WordWrap(output, wrapWidth)
		if err == nil { // success
			output = strings.Join(lines, "\n")
		}
	}

	if outputFile != "" {
		err := os.WriteFile(outputFile, []byte(output), 0o644)
		if err != nil {
			fmt.Printf("error writing to file: %s\n", err)
			os.Exit(1)
		}
		return
	}

	fmt.Println(output)
}
