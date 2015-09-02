package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"resource"
	"version"

	"github.com/ncw/swift"
)

type outRequest struct {
	Resource resource.Source `json:"source"`
	Params   outParams       `json:"params"`
}

type outParams struct {
	From string `json:"from"`
}

type outResponse struct {
	Version  string     `json:"source"`
	Metadata []Metadata `json:"metadata"`
}

type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func main() {
	if len(os.Args) < 2 {
		resource.Sayf("usage: %s <build directory>\n", os.Args[0])
		os.Exit(1)
	}

	buildDir := os.Args[1]

	var request outRequest

	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		resource.Fatal("reading request from stdin", err)
	}

	if request.Params.From == "" {
		resource.Sayf("Required parameter from missing")
		os.Exit(1)
	}

	from, err := regexp.Compile(request.Params.From)
	if err != nil {
		resource.Fatal("Invalid regex in from", err)
	}

	rsc := request.Resource

	regex, err := versions.Regexp(rsc.Regex)
	if err != nil {
		resource.Fatal("Error parsing regex parameter", err)
	}

	//if the from param contains a literal prefix containing slashes
	//we move the search base to the deepest sub directory
	prefix, _ := from.LiteralPrefix()
	dir := regexp.MustCompile("^.*/").FindString(prefix)
	searchBase := filepath.Join(buildDir, dir)

	fileSource := ""
	filepath.Walk(searchBase, func(path string, info os.FileInfo, err error) error {
		if from.MatchString(path) {
			fileSource = path
			return errors.New("Found")
		}
		return nil
	})

	if fileSource == "" {
		resource.Sayf("No file found matching %s", request.Params.From)
		os.Exit(1)
	}

	filename := path.Base(fileSource)
	version, ok := versions.Parse(filename, regex)
	if !ok {
		resource.Sayf("Can't parse version from %s", filename)
		os.Exit(1)
	}

	client := resource.NewClient(rsc)

	file, err := os.Open(fileSource)
	if err != nil {
		resource.Fatal("Can't open source file", err)
	}
	client.ObjectPut(rsc.Container, filename, file, true, "", "", swift.Headers{})
	fi, _ := file.Stat()
	file.Close()

	response := outResponse{
		Version: filename,
		Metadata: []Metadata{
			Metadata{
				Name:  "Version",
				Value: version.VersionNumber,
			},
			Metadata{
				Name:  "Size",
				Value: fmt.Sprintf("%d", fi.Size()),
			},
		},
	}

	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		resource.Fatal("writing response to stdout", err)
	}

}
