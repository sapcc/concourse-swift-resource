package resource

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestOut(t *testing.T) {
	cacheToken = false
	testVersion := "test_1.2.4"
	testServer, source, client, err := testServer([]testObject{})
	if err != nil {
		t.Fatal("Failed to setup swift mock ", err)
	}
	defer testServer.Close()

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("Failed to create source directory %s: %s", dir, err)
	}
	defer os.RemoveAll(dir)
	if err := ioutil.WriteFile(filepath.Join(dir, testVersion), []byte("foo"), 0644); err != nil {
		t.Fatal("Failed to write test file ", err)
	}

	response, err := Out(
		OutRequest{
			Resource: source,
			Params: OutParams{
				From: "test_.*",
			},
		},
		dir)

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

	content, err := client.ObjectGetString(testContainer, testVersion)
	if err != nil {
		t.Fatalf("Error fetching object %s from container %s", testVersion, err)
	}
	if content != "foo" {
		t.Fatalf("Expected object %s to contain %#v, got %#v", testVersion, "foo", content)
	}

}
