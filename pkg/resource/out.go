package resource

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/ncw/swift/v2"

	"github.com/sapcc/concourse-swift-resource/pkg/versions"
)

func Out(ctx context.Context, request OutRequest, sourceDir string) (*OutResponse, error) {
	fileSource, err := prepareFileSource(request, sourceDir)
	if err != nil {
		return nil, err
	}
	filename := path.Base(fileSource)

	version, err := parseVersion(request, filename)
	if err != nil {
		return nil, fmt.Errorf("parsing version failed: %w", err)
	}

	client := NewClient(ctx, request.Resource)

	file, err := os.Open(fileSource)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	container := request.Resource.Container
	if _, _, err := client.Container(ctx, container); err != nil {
		if err := client.ContainerCreate(ctx, container, nil); err != nil {
			return nil, fmt.Errorf("couldn't create Container %s: %w", container, err)
		}
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	headers := swift.Headers{}

	deleteAfter := request.Params.DeleteAfter
	shouldDeleteAfter := deleteAfter != 0
	if shouldDeleteAfter {
		headers["X-Delete-After"] = strconv.FormatInt(deleteAfter, 10)
	}

	bytes := stat.Size()
	if request.Params.SegmentSize == 0 {
		request.Params.SegmentSize = 1073741824
	}

	if bytes > request.Params.SegmentSize {
		if err := uploadLargeObject(ctx, request, client, file, filename, headers); err != nil {
			return nil, fmt.Errorf("failed to upload Large Object to swift: %w", err)
		}
	} else {
		if _, err := client.ObjectPut(ctx, container, filename, file, true, "", "", headers); err != nil {
			return nil, fmt.Errorf("failed to upload to swift: %w", err)
		}
	}

	response := OutResponse{
		Version: Version{Path: filename},
		Metadata: []Metadata{
			{
				Name:  "Version",
				Value: version.VersionNumber,
			},
			{
				Name:  "Size",
				Value: strconv.FormatInt(stat.Size(), 10),
			},
		},
	}

	if shouldDeleteAfter {
		response.Metadata = append(response.Metadata, Metadata{
			Name:  "DeleteAfter",
			Value: strconv.FormatInt(deleteAfter, 10),
		})
	}

	return &response, nil
}

func prepareFileSource(request OutRequest, sourceDir string) (string, error) {
	if request.Params.From == "" {
		return "", fmt.Errorf("required parameter 'from' missing")
	}

	from, err := regexp.Compile(request.Params.From)
	if err != nil {
		return "", fmt.Errorf("invalid regex in from: %w", err)
	}

	//if the from param contains a literal prefix containing slashes
	//we move the search base to the deepest sub directory
	prefix, _ := from.LiteralPrefix()
	dir := regexp.MustCompile("^.*/").FindString(prefix)
	searchBase := filepath.Join(sourceDir, dir)

	fileSource := ""
	err = filepath.Walk(searchBase, func(path string, _ os.FileInfo, _ error) error {
		if from.MatchString(path) {
			fileSource = path
			return nil
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	if fileSource == "" {
		return "", fmt.Errorf("no file found matching %s", request.Params.From)
	}

	return fileSource, nil
}

func parseVersion(request OutRequest, filename string) (versions.Extraction, error) {
	regex, err := versions.Regexp(request.Resource.Regex)
	if err != nil {
		return versions.Extraction{}, fmt.Errorf("error parsing regex parameter: %w", err)
	}

	version, ok := versions.Parse(filename, regex)
	if !ok {
		return versions.Extraction{}, fmt.Errorf("can't parse version from %s", filename)
	}

	return version, nil
}

func uploadLargeObject(ctx context.Context, request OutRequest, client *swift.Connection, file *os.File, filename string, headers swift.Headers) error {
	rsc := request.Resource

	if request.Params.SegmentContainer == "" {
		request.Params.SegmentContainer = rsc.Container + "_segments"
	}
	if _, _, err := client.Container(ctx, request.Params.SegmentContainer); err != nil {
		if err := client.ContainerCreate(ctx, request.Params.SegmentContainer, nil); err != nil {
			return fmt.Errorf("couldn't create Container %s: %w", request.Params.SegmentContainer, err)
		}
	}
	fileHeader := make([]byte, 512)
	if _, err := file.Read(fileHeader); err != nil {
		return fmt.Errorf("couldn't read header information: %w", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("couldn't reset file pointer: %w", err)
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

	out, err := client.StaticLargeObjectCreateFile(ctx, &opts)
	if err != nil {
		return fmt.Errorf("failed to create Static large Object: %w", err)
	}
	_, err = io.Copy(out, file)
	if err != nil {
		return fmt.Errorf("error writing Large Object : %w", err)
	}

	err = out.Close()
	if err != nil {
		return fmt.Errorf("error closing Large Object : %w", err)
	}

	return nil
}
