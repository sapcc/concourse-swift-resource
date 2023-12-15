package resource

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestIn(t *testing.T) {
	ctx := context.TODO()
	cacheToken = false
	testServer, source, _, err := testServer(ctx, []testObject{
		{Path: "test_1.2.3", Content: "foo"},
	})
	if err != nil {
		t.Fatal("Failed to setup swift mock ", err)
	}
	defer testServer.Close()
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatalf("Failed to create test directory %s: %s", dir, err)
	}
	defer os.RemoveAll(dir)

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
	//Clear out the Last modified metadata from response
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
