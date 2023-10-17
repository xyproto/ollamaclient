# ollamaclient

This is a Go package for using Ollama.

The default model is `nous-hermes:latest`.

Example programs are in the `cmd` folder.

### Getting started

1. Install `ollama` and start it as a service.
2. Run `ollama pull nous-hermes` to fetch the `nous-hermes` model.
3. Install the summarizer utility: `go install github.com/xyproto/ollamaclient/cmd/summarizer@latest`
4. Summarize a README.md file and a source code file: `summarize README.md ollamaclient.go`
5. Write a poem about one or more files: `summarize --prompt "Write a poem about the following files:" README.md`

### General info

* Version: 1.0.0
* License: Apache2
* Author: Alexander F. Rødseth
