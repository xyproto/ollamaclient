# ollamaclient

[![Go Reference](https://pkg.go.dev/badge/github.com/xyproto/ollamaclient.svg)](https://pkg.go.dev/github.com/xyproto/ollamaclient/v2)

A Go package for using Ollama and large language models (LLMs).

### Example use

```go
package main

import (
    "fmt"

    "github.com/xyproto/ollamaclient/v2"
    "github.com/xyproto/usermodel"
)

func main() {
    oc := ollamaclient.New(usermodel.GetTextGenerationModel())
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

Example output (with verbosity set to `true`):

```
Sending request to http://localhost:11434/api/tags
Sending request to http://localhost:11434/api/generate: {"model":"gemma2:2b","prompt":"Write a haiku about the color of cows.","options":{"seed":256,"tempera}

Brown hides, gentle eyes,
Mooing low in grassy fields,
Milk flows, life's sweet hue.
```

Make sure to install and run Ollama first, or set `OLLAMA_HOST` to a valid host.

### Using images in the prompt, with the `llava` model

A simple way to describe images:

```go
package main

import (
    "fmt"
    "log"

    "github.com/xyproto/ollamaclient/v2"
    "github.com/xyproto/usermodel"
)

func main() {
    model := usermodel.GetVisionModel()
    oc := ollamaclient.New(model)
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
```

See `v2/cmd/describeimage` for an example that uses a custom prompt.

### Embeddings

* The `.Embeddings` method can be used to pass in a prompt and return a `[]float64`.

### Environment variables

These environment variables are supported:

* `OLLAMA_HOST` (`http://localhost:11434` by default)
* `OLLAMA_MODEL` (uses the model defined by [`llm-manager`](https://github.com/xyproto/llm-manager) by default)
* `OLLAMA_VERBOSE` (`false` by default)

### The `summarize` utility

Getting started:

1. Install `ollama` and start it as a service.
2. Install the `summarize` utility: `go install github.com/xyproto/ollamaclient/cmd/summarize@latest`
3. Summarize a README.md file and a source code file: `summarize README.md ollamaclient.go`. This will also download the model if it's the first run.
4. Write a poem about one or more files: `summarize --prompt "Write a poem about the following files:" README.md`

Usage:

```bash
./summarize [flags] <filename1> [<filename2> ...]
```

Flags:

- `-m`, `--model`: Specify an Ollama model.
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

`go test` depends on a local Ollama server being up and running, and will attempt to download and use various models.

### General info

* Version: 2.7.0
* License: Apache 2
