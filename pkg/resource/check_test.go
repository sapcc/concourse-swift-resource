package resource

import (
	"reflect"
	"testing"

	"github.com/ncw/swift"
	"github.com/ncw/swift/swifttest"
)

func TestCheckResource(t *testing.T) {
	cacheToken = false
	testServer, source, _, err := testServer([]testObject{
		{Path: "test_0.3.0"},
		{Path: "test_1.2.3"},
		{Path: "test_1.0.0"},
	})
	if err != nil {
		t.Fatal("Failed to setup swift mock ", err)
	}
	defer testServer.Close()

	//No version given
	versions, err := Check(CheckRequest{Resource: source})
	if err != nil {
		t.Fatal("check failed: ", err)
	}

	expected := []Version{{Path: "test_1.2.3"}}
	if !reflect.DeepEqual(versions, expected) {
		t.Fatalf("Expected %v, got %v", expected, versions)
	}

	//version given
	versions, err = Check(CheckRequest{
		Resource: source,
		Version:  Version{Path: "test_0.3.0"},
	})
	if err != nil {
		t.Fatal("check failed: ", err)
	}

	expected = []Version{
		{Path: "test_1.0.0"},
		{Path: "test_1.2.3"},
	}
	if !reflect.DeepEqual(versions, expected) {
		t.Fatalf("Expected %v, got %v", expected, versions)
	}

	//No new version
	versions, err = Check(CheckRequest{
		Resource: source,
		Version:  Version{Path: "test_1.2.3"},
	})
	if err != nil {
		t.Fatal("check failed: ", err)
	}

	expected = []Version{}
	if !reflect.DeepEqual(versions, expected) {
		t.Fatalf("Expected %v, got %v", expected, versions)
	}
}

type testObject struct {
	Container string
	Path      string
	Content   string
}

var testContainer = "container"

func testServer(objects []testObject) (server *swifttest.SwiftServer, source Source, client *swift.Connection, err error) {
	if server, err = swifttest.NewSwiftServer("localhost"); err != nil {
		return
	}
	defer func() {
		if err != nil {
			server.Close()
		}
	}()

	client = &swift.Connection{
		UserName: swifttest.TEST_ACCOUNT,
		ApiKey:   swifttest.TEST_ACCOUNT,
		AuthUrl:  server.AuthURL,
	}

	if err = client.Authenticate(); err != nil {
		return
	}

	for _, o := range objects {
		container := testContainer
		if o.Container != "" {
			container = o.Container
		}
		if err = client.ContainerCreate(container, nil); err != nil {
			return
		}
		if err = client.ObjectPutString(container, o.Path, o.Content, ""); err != nil {
			return
		}
	}

	source.Username = swifttest.TEST_ACCOUNT
	source.AuthURL = server.AuthURL
	source.APIKey = swifttest.TEST_ACCOUNT
	source.Container = testContainer
	source.Regex = "^test_(.*)$"

	return
}
