// Copyright 2024 SAP SE
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
