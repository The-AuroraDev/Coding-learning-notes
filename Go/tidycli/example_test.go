// © Roscoe Skeens <web@rskeens.com>
// SPDX-License-Identifier: MIT

package tidycli_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/rskeens/tidycli"
)

// Example demonstrates min cmd tree routed through [tidycli.RunArgs].
// "hello" cmd requires one positional arg and prints greeting.
func Example() {
	cmds := tidycli.Sub{
		"hello": {
			Help: errors.New("hello NAME"),
			Min:  1,

			Fn: func(_ context.Context, args []string) error {
				fmt.Println("hello,", args[0])

				return nil
			},
		},
	}

	if err := tidycli.RunArgs(context.Background(), cmds, []string{"hello", "world"}); err != nil {
		fmt.Println("error:", err)
	}

	// Output: hello, world
}

// ExampleRunArgs demonstrates [tidycli.RunArgs]. Returns [tidycli.ErrArgInvalid]
// for unknown cmd keyword without proc exit.
func ExampleRunArgs() {
	cmds := tidycli.Sub{
		"echo": {
			Help: errors.New("echo TEXT"),
			Min:  1,

			Fn: func(_ context.Context, args []string) error {
				fmt.Println(args[0])

				return nil
			},
		},
	}

	err := tidycli.RunArgs(context.Background(), cmds, []string{"no"})

	fmt.Println(errors.Is(err, tidycli.ErrArgInvalid))

	// Output: true
}
