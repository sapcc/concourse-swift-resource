package resource

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ncw/swift"

	"github.com/sapcc/concourse-swift-resource/pkg/versions"
)

func In(request InRequest, destinationDir string) (*InResponse, error) {
	rsc := request.Resource
	regex, err := versions.Regexp(rsc.Regex)
	if err != nil {
		return nil, fmt.Errorf("error parsing regular expression: %s", err)
	}

	filename := request.Version.Path

	ver, ok := versions.Parse(filename, regex)
	if !ok {
		return nil, fmt.Errorf("can't extract version from %#v", filename)
	}

	if err = os.MkdirAll(destinationDir, 0755); err != nil {
		return nil, fmt.Errorf("can't create destination directory; %s", err)
	}

	client := NewClient(rsc)

	file, err := os.OpenFile(filepath.Join(destinationDir, filename), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("can't open file %s: %s", filepath.Join(destinationDir, filename), err)
	}
	headers, err := client.ObjectGet(rsc.Container, filename, file, true, swift.Headers{})
	file.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch object: %s", err)
	}

	if err = os.WriteFile(filepath.Join(destinationDir, "version"), []byte(ver.VersionNumber), 0644); err != nil {
		return nil, fmt.Errorf("failed to write version file: %s", err)
	}

	if err = os.WriteFile(filepath.Join(destinationDir, "filename"), []byte(filename), 0644); err != nil {
		return nil, fmt.Errorf("failed to write version file: %s", err)
	}

	response := InResponse{
		Version: Version{Path: filename},
		Metadata: []Metadata{
			{
				Name:  "Version",
				Value: ver.VersionNumber,
			},
			{
				Name:  "Size",
				Value: headers["Content-Length"],
			},
			{
				Name:  "Last Modified",
				Value: headers["Last-Modified"],
			},
		},
	}

	return &response, nil
}
