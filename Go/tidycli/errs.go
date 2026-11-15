// © Roscoe Skeens <web@rskeens.com>
// SPDX-License-Identifier: MIT

package tidycli

import "errors"

var (
	// ErrArgNone means no specified args.
	ErrArgNone = errors.New("tidycli: no args specified")
	
	// ErrArgInvalid means arg did not match command. Returned for
	// unknown top level arg and any unknown token while descending
	// sub command [Sub]. Wraps with offending token.
	ErrArgInvalid = errors.New("tidycli: invalid arg specified")

	// ErrFlagParse means [flag.FlagSet.Parse] returned err while
	// parsing [Cmd.Flags] against trailing args. Err from flag pkg is
	// wrapped to distinguish parse failure from routing failure or [Fn]
	// failure with [errors.Is].
	ErrFlagParse = errors.New("tidycli: flag parse failed")

	// ErrFlagReset means [flag.Value] returned err from [flag.Value.Set] when
	// [Sub.RouteArgs] tried to restore it to default before parsing.
	// Wrapped with offending flag name and the underlying [flag.Value.Set] error
	// to distinguish reset failure from parse failure with [errors.Is].
	ErrFlagReset = errors.New("tidycli: flag reset failed")

	// ErrCmdNoFn means resolved leaf command [Cmd.Fn] is nil.
	// Wraps with [Cmd.Help] to surface usage.
	ErrCmdNoFn = errors.New("tidycli: command has no func")
	
	// ErrCmdNoHelp means routing wanted to return [Cmd.Help] but [Cmd.Help] was nil.
	ErrCmdNoHelp = errors.New("tidycli: command has no help")

	// ErrCmdNil means [Sub] held nil [*Cmd] for key routing / rendering.
	// [Sub] vals must not be nil.
	ErrCmdNil = errors.New("tidycli: nil command")
)
