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

	"github.com/sapcc/concourse-swift-resource/pkg/versions"
)

func Check(ctx context.Context, request CheckRequest) ([]Version, error) {
	rsc := request.Resource
	regex, err := versions.Regexp(rsc.Regex)
	if err != nil {
		return nil, fmt.Errorf("invalid regular expression: %w", err)
	}

	client := NewClient(ctx, rsc)
	if cacheToken {
		defer CacheClientToken(client)
	}
	names, err := client.ObjectNamesAll(ctx, rsc.Container, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to enumerate objects: %w", err)
	}
	extractions, err := versions.Extract(names, regex)
	if err != nil {
		return nil, err
	}
	response := []Version{}
	if len(extractions) > 0 {
		if request.Version.Path == "" {
			response = append(response, Version{Path: extractions[len(extractions)-1].Path})
		} else {
			lastVersion, ok := versions.Parse(request.Version.Path, regex)
			if !ok {
				return nil, fmt.Errorf("invalid version. Can't parse %s", request.Version.Path)
			}
			for _, extraction := range extractions {
				if extraction.Version.GreaterThan(lastVersion.Version) {
					response = append(response, Version{Path: extraction.Path})
				}
			}
		}
	}
	return response, nil
}
