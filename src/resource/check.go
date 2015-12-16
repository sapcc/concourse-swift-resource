package resource

import (
	"fmt"
	"version"
)

func Check(request CheckRequest) ([]Version, error) {
	rsc := request.Resource
	regex, err := versions.Regexp(rsc.Regex)

	client := NewClient(rsc)
	if cacheToken {
		defer CacheClientToken(client)
	}
	names, err := client.ObjectNamesAll(rsc.Container, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to enumerate objects: %s", err)
	}
	extractions, err := versions.Extract(names, regex)
	if err != nil {
		return nil, fmt.Errorf("Error: %s", err)
	}
	response := []Version{}
	if len(extractions) > 0 {
		if request.Version.Path == "" {
			response = append(response, Version{Path: extractions[len(extractions)-1].Path})
		} else {

			lastVersion, ok := versions.Parse(request.Version.Path, regex)
			if !ok {
				return nil, fmt.Errorf("Invalid version. Can't parse %s", request.Version.Path)
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
