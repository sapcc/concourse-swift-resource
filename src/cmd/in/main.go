package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"resource"
	"version"

	"github.com/ncw/swift"
)

type inRequest struct {
	Resource resource.Source `json:"source"`
	Version  string          `json:"version"`
}

type inResponse struct {
	Version  string     `json:"source"`
	Metadata []Metadata `json:"metadata"`
}

type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func main() {
	if len(os.Args) < 2 {
		resource.Sayf("usage: %s <dest directory>\n", os.Args[0])
		os.Exit(1)
	}

	destinationDir := os.Args[1]

	var request inRequest

	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		resource.Fatal("reading request from stdin", err)
	}
	rsc := request.Resource
	regex, err := versions.Regexp(rsc.Regex)
	if err != nil {
		resource.Fatal("Error parsing regular expression", err)
	}

	filename := request.Version

	ver, ok := versions.Parse(request.Version, regex)
	if !ok {
		resource.Sayf("Can't extract version from %s", filename)
		os.Exit(1)
	}

	if err = os.MkdirAll(destinationDir, 0755); err != nil {
		resource.Fatal("Can't create destination directory", err)
	}

	client := resource.NewClient(rsc)

	file, err := os.OpenFile(filepath.Join(destinationDir, filename), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		resource.Fatal("Can't open file", err)
	}
	headers, err := client.ObjectGet(rsc.Container, request.Version, file, true, swift.Headers{})
	file.Close()
	if err != nil {
		resource.Fatal("Failed to fetch object", err)
	}

	if err = ioutil.WriteFile(filepath.Join(destinationDir, "version"), []byte(ver.VersionNumber), 0644); err != nil {
		resource.Fatal("Failed to write version file", err)
	}

	if err = ioutil.WriteFile(filepath.Join(destinationDir, "filename"), []byte(filename), 0644); err != nil {
		resource.Fatal("Failed to write version file", err)
	}

	response := inResponse{
		Version: filename,
		Metadata: []Metadata{
			Metadata{
				Name:  "Version",
				Value: ver.VersionNumber,
			},
			Metadata{
				Name:  "Size",
				Value: headers["Content-Length"],
			},
			Metadata{
				Name:  "Last Modified",
				Value: headers["Last-Modified"],
			},
		},
	}

	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		resource.Fatal("writing response to stdout", err)
	}
}
