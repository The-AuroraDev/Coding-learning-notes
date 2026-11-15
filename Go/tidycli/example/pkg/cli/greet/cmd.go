// © Roscoe Skeens <web@rskeens.com>
// SPDX-License-Identifier: MIT

// Package greet implements "greet" parent command and "hi" / "bye" leaves.
package greet

import (
	"github.com/rskeens/tidycli"
	"github.com/rskeens/tidycli/example/pkg/cli/help"
)

// Cmd returns cmd.
func Cmd() *tidycli.Cmd {
	return &tidycli.Cmd{
		Help: help.ErrGreet,
		Min:  1,

		Sub: tidycli.Sub{
			"hi": {
				Help: help.ErrGreet,
				Min:  1,
				Fn:   Hi,
			},

			"bye": {
				Help: help.ErrGreet,
				Min:  1,
				Fn:   Bye,
			},
		},
	}
}
