// © Roscoe Skeens <web@rskeens.com>
// SPDX-License-Identifier: MIT

package main

import (
	"flag"
	"os"
	"os/exec"
	"testing"
)

// TestMainHello drives main() with a valid command to assert that
// argument parsing and dispatch wire up correctly.
func TestMainHello(t *testing.T) {
	prevArgs := os.Args
	prevFlag := flag.CommandLine

	os.Args = []string{"hello-world", "hello", "world"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	t.Cleanup(func() {
		os.Args = prevArgs
		flag.CommandLine = prevFlag
	})

	main()
}

// TestMainErr re-invokes the test binary as a subprocess so the
// log.Fatal branch in main() can be exercised without terminating the
// test runner itself.
func TestMainErr(t *testing.T) {
	if os.Getenv("TIDYCLI_TEST_MAIN_ERR") == "1" {
		os.Args = []string{"hello-world"}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMainErr")
	cmd.Env = append(os.Environ(), "TIDYCLI_TEST_MAIN_ERR=1")

	if err := cmd.Run(); err == nil {
		t.Errorf("expected subprocess to exit non-zero, got nil")
	}
}
