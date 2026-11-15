// © Roscoe Skeens <web@rskeens.com>
// SPDX-License-Identifier: MIT

package hello

import "testing"

func TestCmd(t *testing.T) {
	t.Parallel()

	if Cmd() == nil {
		t.Error("expected non nil cmd")
	}
}
