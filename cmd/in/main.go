// SPDX-FileCopyrightText: 2015-2023 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/sapcc/concourse-swift-resource/pkg/resource"
)

func main() {
	if len(os.Args) < 2 {
		resource.Sayf("usage: %s <dest directory>\n", os.Args[0])
		os.Exit(1)
	}
	destinationDir := os.Args[1]

	var request resource.InRequest
	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		resource.Fatal("reading request from stdin", err)
	}

	response, err := resource.In(context.TODO(), request, destinationDir)
	if err != nil {
		resource.Fail(err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		resource.Fatal("writing response to stdout", err)
	}
}
