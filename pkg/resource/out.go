package resource

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/ncw/swift"

	"github.com/sapcc/concourse-swift-resource/pkg/versions"
)

func Out(request OutRequest, sourceDir string) (*OutResponse, error) {
	if request.Params.From == "" {
		return nil, fmt.Errorf("Required parameter 'from' missing")
	}

	from, err := regexp.Compile(request.Params.From)
	if err != nil {
		return nil, fmt.Errorf("Invalid regex in from: %s", err)
	}

	rsc := request.Resource

	regex, err := versions.Regexp(rsc.Regex)
	if err != nil {
		return nil, fmt.Errorf("Error parsing regex parameter: %s", err)
	}

	//if the from param contains a literal prefix containing slashes
	//we move the search base to the deepest sub directory
	prefix, _ := from.LiteralPrefix()
	dir := regexp.MustCompile("^.*/").FindString(prefix)
	searchBase := filepath.Join(sourceDir, dir)

	fileSource := ""
	filepath.Walk(searchBase, func(path string, info os.FileInfo, err error) error {
		if from.MatchString(path) {
			fileSource = path
			return errors.New("Found")
		}
		return nil
	})

	if fileSource == "" {
		return nil, fmt.Errorf("No file found matching %s", request.Params.From)
	}

	filename := path.Base(fileSource)
	version, ok := versions.Parse(filename, regex)
	if !ok {
		return nil, fmt.Errorf("Can't parse version from %s", filename)
	}

	client := NewClient(rsc)

	file, err := os.Open(fileSource)
	if err != nil {
		return nil, fmt.Errorf("Can't open source file %s: %s", fileSource, err)
	}
	defer file.Close()
	if _, _, err := client.Container(rsc.Container); err != nil {
		if err := client.ContainerCreate(rsc.Container, nil); err != nil {
			return nil, fmt.Errorf("Couldn't create Container %s: %s", rsc.Container, err)
		}
	}
	if _, err := client.ObjectPut(rsc.Container, filename, file, true, "", "", swift.Headers{}); err != nil {
		return nil, fmt.Errorf("Failed to upload to swift: %s", err)
	}
	fi, _ := file.Stat()

	response := OutResponse{
		Version: Version{Path: filename},
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
	return &response, nil

}
