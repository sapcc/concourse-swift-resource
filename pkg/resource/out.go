package resource

import (
	"errors"
	"fmt"
	"io"
	"net/http"
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

	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("Can't stats of file %s: %s", fileSource, err)
	}

	headers := swift.Headers{}

	expires := request.Params.Expires
	shouldExpire := expires != 0
	if shouldExpire {
		headers["X-Delete-After"] = fmt.Sprintf("%v", expires)
	}

	var bytes int64
	bytes = stat.Size()
	if request.Params.SegmentSize == 0 {
		request.Params.SegmentSize = 1073741824
	}

	if bytes > request.Params.SegmentSize {
		if err := uploadLargeObject(request, client, file, filename, headers); err != nil {
			return nil, fmt.Errorf("Failed to upload Large Object to swift: %s", err)
		}
	} else {
		if _, err := client.ObjectPut(rsc.Container, filename, file, true, "", "", headers); err != nil {
			return nil, fmt.Errorf("Failed to upload to swift: %s", err)
		}
	}

	response := OutResponse{
		Version: Version{Path: filename},
		Metadata: []Metadata{
			Metadata{
				Name:  "Version",
				Value: version.VersionNumber,
			},
			Metadata{
				Name:  "Size",
				Value: fmt.Sprintf("%d", stat.Size()),
			},
		},
	}
	return &response, nil
}

func uploadLargeObject(request OutRequest, client *swift.Connection, file *os.File, filename string, headers swift.Headers) error {
	rsc := request.Resource

	if request.Params.SegmentContainer == "" {
		request.Params.SegmentContainer = rsc.Container + "_segments"
	}
	if _, _, err := client.Container(request.Params.SegmentContainer); err != nil {
		if err := client.ContainerCreate(request.Params.SegmentContainer, nil); err != nil {
			return fmt.Errorf("Couldn't create Container %s: %s", request.Params.SegmentContainer, err)
		}
	}
	fileHeader := make([]byte, 512)
	if _, err := file.Read(fileHeader); err != nil {
		return fmt.Errorf("Couldn't read header information: %s", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("Couldn't reset file pointer: %s", err)
	}

	opts := swift.LargeObjectOpts{
		Container:        rsc.Container,
		ObjectName:       filename,
		ContentType:      http.DetectContentType(fileHeader),
		Headers:          headers,
		ChunkSize:        request.Params.SegmentSize,
		MinChunkSize:     request.Params.SegmentSize,
		SegmentContainer: request.Params.SegmentContainer,
	}

	out, err := client.StaticLargeObjectCreateFile(&opts)
	if err != nil {
		return fmt.Errorf("Failed to create Static large Object: %s", err)
	}
	_, err = io.Copy(out, file)
	if err != nil {
		return fmt.Errorf("Error writing Large Object : %s", err)
	}

	err = out.Close()
	if err != nil {
		return fmt.Errorf("Error closing Large Object : %s", err)
	}

	return nil
}
