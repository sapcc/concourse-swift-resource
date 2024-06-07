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
	"context"
	"reflect"
	"testing"

	"github.com/ncw/swift/v2"
	"github.com/ncw/swift/v2/swifttest"
)

func TestCheckResource(t *testing.T) {
	ctx := context.TODO()
	cacheToken = false
	testServer, source, _, err := testServer(ctx, []testObject{
		{Path: "test_0.3.0"},
		{Path: "test_1.2.3"},
		{Path: "test_1.0.0"},
	})
	if err != nil {
		t.Fatal("Failed to setup swift mock ", err)
	}
	defer testServer.Close()

	// No version given
	versions, err := Check(ctx, CheckRequest{Resource: source})
	if err != nil {
		t.Fatal("check failed: ", err)
	}

	expected := []Version{{Path: "test_1.2.3"}}
	if !reflect.DeepEqual(versions, expected) {
		t.Fatalf("Expected %v, got %v", expected, versions)
	}

	// version given
	versions, err = Check(ctx, CheckRequest{
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

	// No new version
	versions, err = Check(ctx, CheckRequest{
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

func testServer(ctx context.Context, objects []testObject) (server *swifttest.SwiftServer, source Source, client *swift.Connection, err error) {
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

	if err = client.Authenticate(ctx); err != nil {
		return
	}

	for _, o := range objects {
		container := testContainer
		if o.Container != "" {
			container = o.Container
		}
		if err = client.ContainerCreate(ctx, container, nil); err != nil {
			return
		}
		if err = client.ObjectPutString(ctx, container, o.Path, o.Content, ""); err != nil {
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
