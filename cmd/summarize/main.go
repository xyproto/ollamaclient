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

const (
	versionString    = "Summarize 1.4.0"
	defaultModel     = "nous-hermes:latest"
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
		noHeader, showVersion           bool
	)

	pflag.BoolVarP(&verbose, "verbose", "V", false, "verbose output")
	pflag.StringVarP(&promptHeader, "prompt", "p", "Write a short summary of what a project that contains the following files is:", "Provide a custom prompt header")
	pflag.StringVarP(&outputFile, "output", "o", "", "Specify an output file")
	pflag.StringVarP(&model, "model", "m", defaultModel, "Specify the Ollama model to use")
	pflag.IntVarP(&wrapWidth, "wrap", "w", 0, "Word wrap at specified width. Use '-1' for terminal width")
	pflag.BoolVarP(&noHeader, "no-header", "n", false, "Do not include filenames in the prompt")
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
		fmt.Println("Usage: summarize [--prompt <customPrompt>] [--output <outputFile>] [--wrap <width>|-1] [--model <ollamaModel>] [--no-header] <filename1> [<filename2> ...]")
		os.Exit(1)
	}

	var sb strings.Builder
	for _, filename := range filenames {
		logVerbose("[%s] Reading... ", filename)
		data, err := os.ReadFile(filename)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		logVerbose("OK\n")

		if !noHeader {
			sb.WriteString(filename + ":\n")
		}
		sb.Write(data)
	}

	prompt := promptHeader + "\n\n" + sb.String()

	oc := ollamaclient.New()
	if model != defaultModel {
		oc.Model = model
	}

	logVerbose("[%s] Downloading model, if needed... ", oc.Model)
	if _, err := oc.Pull(); err != nil {
		fmt.Println("Error:", err)
		return
	}
	logVerbose("OK\n")

	logVerbose("[%s] Generating... ", oc.Model)
	output, err := oc.GetOutput(prompt)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

	logVerbose("OK\n")
	trimmedOutput := strings.TrimSpace(output)

	if wrapWidth > 0 {
		lines, err := wordwrap.WordWrap(trimmedOutput, wrapWidth)
		if err == nil { // success
			trimmedOutput = strings.Join(lines, "\n")
		}
	}

	if outputFile != "" {
		err := os.WriteFile(outputFile, []byte(trimmedOutput), 0o644)
		if err != nil {
			fmt.Printf("error writing to file: %s\n", err)
			os.Exit(1)
		}
		return
	}

	fmt.Println(trimmedOutput)
}
