package versions

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderedExtraction(t *testing.T) {

	filenames := []string{
		"file-2.0",
		"file-1.0",
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

	for i, path := range expected {
		assert.Equal(t, path, extractions[i].Path)
	}
}

func TestUnorderedExtraction(t *testing.T) {
	filenames := []string{
		"file-b123",
		"file-a123",
		"file-34asd",
		"file-x",
	}

	extractions, err := Extract(filenames, regexp.MustCompile("file-([a-z0-9]+)"))
	if err != nil {
		t.Fatal(err)
	}

	for i, path := range filenames {
		vnumber := strings.Replace(path, "file-", "", 1)

		assert.Equal(t, path, extractions[i].Path)
		assert.Equal(t, "0.0.0+"+vnumber, extractions[i].Version.String())
		assert.Equal(t, vnumber, extractions[i].VersionNumber)
	}
}
