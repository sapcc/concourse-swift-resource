/*******************************************************************************
*
* Copyright 2015-2023 SAP SE
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You should have received a copy of the License along with this
* program. If not, you may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*
*******************************************************************************/

package resource

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestIn(t *testing.T) {
	ctx := t.Context()
	cacheToken = false
	testServer, source, _, err := testServer(ctx, []testObject{
		{Path: "test_1.2.3", Content: "foo"},
	})
	if err != nil {
		t.Fatal("Failed to setup swift mock ", err)
	}
	defer testServer.Close()
	dir := t.TempDir()

	response, err := In(ctx, InRequest{Resource: source, Version: Version{Path: "test_1.2.3"}}, dir)
	if err != nil {
		t.Fatal("check failed: ", err)
	}

	expected := &InResponse{
		Version: Version{Path: "test_1.2.3"},
		Metadata: []Metadata{
			{Name: "Version", Value: "1.2.3"},
			{Name: "Size", Value: "3"},
			{Name: "Last Modified", Value: ""},
		},
	}
	// Clear out the Last modified metadata from response
	for i, m := range response.Metadata {
		if m.Name == "Last Modified" {
			response.Metadata[i].Value = ""
		}
	}

	if !reflect.DeepEqual(expected, response) {
		t.Fatalf("Expected %v, got %v", expected, response)
	}

	if content, err := os.ReadFile(filepath.Join(dir, "version")); err != nil || string(content) != "1.2.3" {
		t.Fatalf("Expected to find file %s with content %s", filepath.Join(dir, "version"), content)
	}
	if content, err := os.ReadFile(filepath.Join(dir, "filename")); err != nil || string(content) != "test_1.2.3" {
		t.Fatalf("Expected to find file %s with content %s", filepath.Join(dir, "filename"), content)
	}
	if content, err := os.ReadFile(filepath.Join(dir, "test_1.2.3")); err != nil || string(content) != "foo" {
		t.Fatalf("Expected to find file %s with content %s", filepath.Join(dir, "test_1.2.3"), content)
	}
}
