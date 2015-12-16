package versions

import (
	"errors"
	"regexp"
	"sort"

	version "github.com/hashicorp/go-version"
)

type Extractions []Extraction

func (e Extractions) Len() int {
	return len(e)
}

func (e Extractions) Less(i int, j int) bool {
	return e[i].Version.LessThan(e[j].Version)
}

func (e Extractions) Swap(i int, j int) {
	e[i], e[j] = e[j], e[i]
}

type Extraction struct {
	// path to an object in bucket
	Path string

	// parsed semantic version
	Version *version.Version

	// the raw version match
	VersionNumber string
}

func Parse(filename string, pattern *regexp.Regexp) (Extraction, bool) {
	matches := pattern.FindStringSubmatch(filename)
	if matches == nil || len(matches) < 2 {
		return Extraction{}, false
	}
	ver, err := version.NewVersion(matches[1])
	if err != nil {
		return Extraction{}, false
	}
	return Extraction{
		Path:          filename,
		VersionNumber: matches[1],
		Version:       ver,
	}, true
}

func Extract(filenames []string, pattern *regexp.Regexp) (Extractions, error) {
	extractions := make(Extractions, 0, len(filenames))

	for _, filename := range filenames {
		if extraction, ok := Parse(filename, pattern); ok {
			extractions = append(extractions, extraction)
		}
	}

	sort.Sort(extractions)
	return extractions, nil
}

func Regexp(pattern string) (*regexp.Regexp, error) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	if regex.NumSubexp() != 1 {
		return nil, errors.New("Regular expression needs to have exactly one subexpression")
	}

	return regex, nil

}
