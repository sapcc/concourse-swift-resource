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
	"fmt"
	"os"
	"path/filepath"

	"github.com/ncw/swift/v2"

	"github.com/sapcc/concourse-swift-resource/pkg/versions"
)

func In(ctx context.Context, request InRequest, destinationDir string) (*InResponse, error) {
	rsc := request.Resource
	regex, err := versions.Regexp(rsc.Regex)
	if err != nil {
		return nil, fmt.Errorf("error parsing regular expression: %w", err)
	}

	filename := request.Version.Path

	ver, ok := versions.Parse(filename, regex)
	if !ok {
		return nil, fmt.Errorf("can't extract version from %#v", filename)
	}

	if err = os.MkdirAll(destinationDir, 0755); err != nil {
		return nil, err
	}

	client := NewClient(ctx, rsc)

	file, err := os.OpenFile(filepath.Join(destinationDir, filename), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}
	headers, err := client.ObjectGet(ctx, rsc.Container, filename, file, true, swift.Headers{})
	file.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch object: %w", err)
	}

	if err = os.WriteFile(filepath.Join(destinationDir, "version"), []byte(ver.VersionNumber), 0644); err != nil {
		return nil, err
	}

	if err = os.WriteFile(filepath.Join(destinationDir, "filename"), []byte(filename), 0644); err != nil {
		return nil, err
	}

	response := InResponse{
		Version: Version{Path: filename},
		Metadata: []Metadata{
			{
				Name:  "Version",
				Value: ver.VersionNumber,
			},
			{
				Name:  "Size",
				Value: headers["Content-Length"],
			},
			{
				Name:  "Last Modified",
				Value: headers["Last-Modified"],
			},
		},
	}

	return &response, nil
}
