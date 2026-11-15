// © Roscoe Skeens <web@rskeens.com>
// SPDX-License-Identifier: MIT

package hello

import (
	"context"
	"fmt"
)

// Do prints greeting for given name.
// hello-world hello NAME
func Do(_ context.Context, args []string) error {
	_, err := fmt.Println("hello,", args[0])

	return err
}
