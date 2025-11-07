// SPDX-FileCopyrightText: 2015-2023 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

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
