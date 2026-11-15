// © Roscoe Skeens <web@rskeens.com>
// SPDX-License-Identifier: MIT

package greet

import (
	"context"
	"fmt"
)

// Hi prints "hi" greeting for given name.
// hello-world greet hi NAME
func Hi(_ context.Context, args []string) error {
	_, err := fmt.Println("hi,", args[0])

	return err
}

// Bye prints "bye" farewell for given name.
// hello-world greet bye NAME
func Bye(_ context.Context, args []string) error {
	_, err := fmt.Println("bye,", args[0])

	return err
}
