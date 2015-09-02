package main

import (
	"encoding/json"
	"fmt"
	"os"
	"version"

	"resource"
)

type checkRequest struct {
	Resource resource.Source `json:"source"`
	Version  string          `json:"version"`
}

func main() {
	var request checkRequest

	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		resource.Fatal("reading request from stdin", err)
	}
	rsc := request.Resource
	regex, err := versions.Regexp(rsc.Regex)
	if err != nil {
		resource.Fatal("Error parsing regular expression", err)
	}

	client := resource.NewClient(rsc)
	names, err := client.ObjectNamesAll(rsc.Container, nil)
	//names := []string{"file-3.0", "file-2.1", "file-2.3"}
	if err != nil {
		resource.Fatal("Failed to enumerate objects", err)
	}
	extractions, err := versions.Extract(names, regex)
	if err != nil {
		resource.Fatal("Error", err)
	}
	response := []string{}
	if len(extractions) > 0 {
		if request.Version == "" {
			response = append(response, extractions[len(extractions)-1].Path)
		} else {

			lastVersion, ok := versions.Parse(request.Version, regex)
			if !ok {
				resource.Fatal("Invalid version", fmt.Errorf("Can't parse %s", request.Version))
			}
			for _, extraction := range extractions {
				if extraction.Version.GreaterThan(lastVersion.Version) {
					response = append(response, extraction.Path)
				}
			}
		}
	}

	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		resource.Fatal("writing response to stdout", err)
	}
}
