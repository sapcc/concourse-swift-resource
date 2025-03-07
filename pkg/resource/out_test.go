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

func TestOut(t *testing.T) {
	ctx := t.Context()
	cacheToken = false
	testVersion := "test_1.2.4"
	testServer, source, client, err := testServer(ctx, []testObject{})
	if err != nil {
		t.Fatal("Failed to setup swift mock ", err)
	}
	defer testServer.Close()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, testVersion), []byte("foo"), 0644); err != nil {
		t.Fatal("Failed to write test file ", err)
	}

	response, err := Out(
		ctx,
		OutRequest{
			Resource: source,
			Params: OutParams{
				From: "test_.*",
			},
		},
		dir,
	)

	if err != nil {
		t.Fatal("check failed: ", err)
	}

	expected := &OutResponse{
		Version: Version{Path: testVersion},
		Metadata: []Metadata{
			{Name: "Version", Value: "1.2.4"},
			{Name: "Size", Value: "3"},
		},
	}
	if !reflect.DeepEqual(expected, response) {
		t.Fatalf("Expected %v, got %v", expected, response)
	}

	content, err := client.ObjectGetString(ctx, testContainer, testVersion)
	if err != nil {
		t.Fatalf("Error fetching object %s from container %s", testVersion, err)
	}
	if content != "foo" {
		t.Fatalf("Expected object %s to contain %#v, got %#v", testVersion, "foo", content)
	}
}
