package versions

import (
	"regexp"
	"testing"
)

func TestExtraction(t *testing.T) {
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
		if extractions[i].Path != path {
			t.Errorf("Expected %#v got %#v", path, extractions[i].Path)
		}
	}
}
