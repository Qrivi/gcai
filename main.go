// Copyright 2024 Kristof Dewilde <code@qrivi.dev>. All rights reserved.
// This is the first time I write in Go so you probably don't wanna steal this code anyway.
package main

import (
	"fmt"
	"github.com/qrivi/gcai/cmd"
	"os"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
