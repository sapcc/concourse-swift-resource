package main

import (
	"encoding/json"
	"os"

	"github.com/sapcc/concourse-swift-resource/pkg/resource"
)

func main() {
	if len(os.Args) < 2 {
		resource.Sayf("usage: %s <source directory>\n", os.Args[0])
		os.Exit(1)
	}
	buildDir := os.Args[1]

	var request resource.OutRequest
	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		resource.Fatal("reading request from stdin", err)
	}

	response, err := resource.Out(request, buildDir)
	if err != nil {
		resource.Fail(err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		resource.Fatal("writing response to stdout", err)
	}
}
