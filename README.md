# ollamaclient

A Go package for using Ollama and large language models (LLMs).

### Example use

```go
package main

import (
    "fmt"

    "github.com/xyproto/ollamaclient/v2"
)

func main() {
    oc := ollamaclient.New()

    // Select a model. The default model is "tinyllama".
    // See: https://ollama.com/library/
    //oc.ModelName = "gemma:latest"
    //oc.ModelName = "mixtral:instruct"
    //oc.ModelName = "nous-hermes:7b-llama2-q2_K"

    oc.Verbose = true

    if err := oc.PullIfNeeded(); err != nil {
        fmt.Println("Error:", err)
        return
    }

    prompt := "Write a haiku about the color of cows."
    output, err := oc.GetOutput(prompt)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Printf("\n%s\n", output)
}
```

Example output:

```
Sending request to /api/tags
Sending request to /api/generate: {"model":"mistral:instruct","prompt":"Write a haiku about the color of cows."}

Majestic brown cows
Grazing in green fields so serene
Nature's masterpiece
```

Make sure to install and run Ollama first.

### Using images in the prompt, with the `llava` model

Example use, loading in images from a slice of filenames and building a prompt with base64 encoded images, where the output is a description of the images:

```go
imageFilenames := []string{"example1.png", "example2.png"}

var images []string
for _, imageFilename := range imageFilenames {
    if base64image, err := ollamaclient.Base64EncodeFile(imageFilename); err == nil { // success
        // append the base64 encoded image to the "images" string slice
        images = append(images, base64image)
    }
}

var prompt string
switch len(images) {
case 0:
    log.Fatalln("Error: no images to describe")
case 1:
    prompt = "Describe this image:"
default:
    prompt = "Describe these images:"
}

oc := ollamaclient.New()
oc.ModelName = "llava"

if err := oc.PullIfNeeded(verbose); err != nil {
    fmt.Println("Error:", err)
    return
}
oc.SetReproducible()

promptAndImages := append([]string{prompt}, images...)

output, err := oc.GetOutput(promptAndImages...)
```

See `v2/cmd/describeimage` for a runnable example.

### Environment variables

These environment variables are supported:

* `OLLAMA_HOST` (`http://localhost:11434` by default)
* `OLLAMA_MODEL` (`nous-hermes:7b-llama2-q2_K` by default)
* `OLLAMA_VERBOSE` (`false` by default)

### The `summarize` utility

Getting started:

1. Install `ollama` and start it as a service.
2. Run `ollama pull nous-hermes:7b-llama2-q2_K` to fetch the `nous-hermes:7b-llama2-q2_K` model.
3. Install the `summarize` utility: `go install github.com/xyproto/ollamaclient/cmd/summarize@latest`
4. Summarize a README.md file and a source code file: `summarize README.md ollamaclient.go`
5. Write a poem about one or more files: `summarize --prompt "Write a poem about the following files:" README.md`

Usage:

```bash
./summarize [flags] <filename1> [<filename2> ...]
```

Flags:

- `-m`, `--model`: Specify an Ollama model. The default is `nous-hermes:latest`.
- `-o`, `--output`: Define an output file to store the summary.
- `-p`, `--prompt`: Specify a custom prompt header for summary. The default is `Write a short summary of a project that contains the following files:`
- `-w`, `--wrap`: Set the word wrap width. Use -1 to detect the terminal width.
- `-v`, `--version`: Display the current version.
- `-V`, `--verbose`: Enable verbose logging.

Generate a summary with a custom prompt:

```bash
./summarize -w -1 -p "Summarize these files:" README.md CONFIG.md
```

Generate a summary, saving the output to a file:

```bash
./summarize -o output.txt README.md CONFIG.md
```

Generate a summary with custom word wrap width:

```bash
./summarize -w 100 README.md
```

### Testing

`go test` depends on a local Ollama server being up and running, and will attempt to download the `tinyllama` model (637 MiB).

### General info

* Version: 2.1.0
* License: Apache 2
* Author: Alexander F. Rødseth
