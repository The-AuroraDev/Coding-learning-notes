// © Roscoe Skeens <web@rskeens.com>
// SPDX-License-Identifier: MIT

// Command hello-world demonstrates a minimal tidycli app.
// One package per command.
//
// go run ./example/cmd/hello-world hello world
// go run ./example/cmd/hello-world greet hi alice
// go run ./example/cmd/hello-world                 # prints usage
package main

import (
	"context"
	"flag"
	"log"

	"github.com/rskeens/tidycli"
)

func main() {
	flag.Parse()

	if err := tidycli.Run(context.Background(), cmds); err != nil {
		log.Fatal(err)
	}
}
