// © Roscoe Skeens <web@rskeens.com>
// SPDX-License-Identifier: MIT

// Package hello implements "hello" leaf command.
package hello

import (
	"github.com/rskeens/tidycli"
	"github.com/rskeens/tidycli/example/pkg/cli/help"
)

// Cmd returns cmd.
func Cmd() *tidycli.Cmd {
	return &tidycli.Cmd{
		Help: help.ErrHello,
		Min:  1,
		Fn:   Do,
	}
}
