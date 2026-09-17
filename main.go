// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Wouter van der Toorren

package main

import (
	"fmt"
	"os"

	"github.com/torreirow/soltty/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
