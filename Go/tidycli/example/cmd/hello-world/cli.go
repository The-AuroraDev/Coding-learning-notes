// © Roscoe Skeens <web@rskeens.com>
// SPDX-License-Identifier: MIT

package main

import (
	"github.com/rskeens/tidycli"
	"github.com/rskeens/tidycli/example/pkg/cli/greet"
	"github.com/rskeens/tidycli/example/pkg/cli/hello"
)

// cmds stores arg:cmd layout.
var cmds = tidycli.Sub{
	"hello": hello.Cmd(),
	"greet": greet.Cmd(),
}
