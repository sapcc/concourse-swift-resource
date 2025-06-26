// SPDX-FileCopyrightText: 2015-2023 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package resource

import (
	"fmt"
	"os"

	"github.com/mitchellh/colorstring"
)

func Fail(err error) {
	fmt.Fprintf(os.Stderr, colorstring.Color("[red]%s\n"), err)
	os.Exit(1)
}

func Fatal(doing string, err error) {
	Sayf(colorstring.Color("[red]%s: %s\n"), doing, err)
	os.Exit(1)
}

func Sayf(message string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, message, args...)
}
