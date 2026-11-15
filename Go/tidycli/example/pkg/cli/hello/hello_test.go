// © Roscoe Skeens <web@rskeens.com>
// SPDX-License-Identifier: MIT

package hello

import (
	"context"
	"testing"
)

func TestDo(t *testing.T) {
	t.Parallel()

	if err := Do(context.Background(), []string{"world"}); err != nil {
		t.Errorf("do error %v", err)
	}
}
