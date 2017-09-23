package versions

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderedExtraction(t *testing.T) {

	filenames := []string{
		"file-2.0",
		"file-1.0",
		"file-äüö", // ignores non-matching
		"file-3",
	}
	expected := []string{
		"file-1.0",
		"file-2.0",
		"file-3",
	}

	extractions, err := Extract(filenames, regexp.MustCompile("file-([.0-9]+)"))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(expected), len(extractions))
	for i, path := range expected {
		assert.Equal(t, path, extractions[i].Path)
	}
}

func TestUnorderedExtraction(t *testing.T) {
	filenames := []string{
		"file-1.0", // 1.0 is sorted last, the 0.0.0+ numbers come first
		"file-b123",
		"file-a123",
		"file-34asd",
		"file-x",
	}
	expected := []string{
		"b123",
		"a123",
		"34asd",
		"x",
		"1.0",
	}

	extractions, err := Extract(filenames, regexp.MustCompile("file-([a-z0-9.]+)"))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(expected), len(extractions))

	for i, vnumber := range expected {
		extracted := extractions[i]

		assert.Equal(t, "file-"+vnumber, extracted.Path)
		assert.Equal(t, vnumber, extracted.VersionNumber)
		if vnumber != "1.0" {
			assert.Equal(t, "0.0.0+"+vnumber, extracted.Version.String())
		}
	}
}

func TestExtractionEdgeConditions(t *testing.T) {
	filenames := []string{
		"file-1.0", // valid
		"file-äüö", // ignores zero-legnth version group
		"äüö",      // ignores non-matching full string
	}
	expected := []string{
		"1.0",
	}

	// note: the regex allows a zero-length version group here
	extractions, err := Extract(filenames, regexp.MustCompile("file-([0-9.]*)"))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(expected), len(extractions))
	for i, vnumber := range expected {
		extracted := extractions[i]

		assert.Equal(t, "file-"+vnumber, extracted.Path)
		assert.Equal(t, vnumber, extracted.VersionNumber)
		if vnumber != "1.0" {
			assert.Equal(t, "0.0.0+"+vnumber, extracted.Version.String())
		}
	}
}
