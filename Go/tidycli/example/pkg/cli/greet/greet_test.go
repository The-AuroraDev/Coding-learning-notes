// © Roscoe Skeens <web@rskeens.com>
// SPDX-License-Identifier: MIT

package greet

import (
	"context"
	"testing"
)

func TestHi(t *testing.T) {
	t.Parallel()

	if err := Hi(context.Background(), []string{"alice"}); err != nil {
		t.Errorf("hi error %v", err)
	}
}

func TestBye(t *testing.T) {
	t.Parallel()

	if err := Bye(context.Background(), []string{"alice"}); err != nil {
		t.Errorf("bye error %v", err)
	}
}
