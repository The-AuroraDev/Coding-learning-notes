// © Roscoe Skeens <web@rskeens.com>
// SPDX-License-Identifier: MIT

// Package help centralises example apps's per command usage messages.
//
// Stored as exported errs lets each command package import only what it needs.
// Handed to [tidycli.Cmd.Help] without restating usage string at call site.
package help

import "errors"

// https://developers.google.com/style/code-syntax
var (
	// ErrHello means hello help.
	ErrHello = errors.New("hello-world hello NAME")

	// ErrGreet means greet help.
	ErrGreet = errors.New("hello-world greet { hi | bye } NAME")
)
