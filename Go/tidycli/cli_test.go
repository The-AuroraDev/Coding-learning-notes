// © Roscoe Skeens <web@rskeens.com>
// SPDX-License-Identifier: MIT

package tidycli

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

// input represents inputs.
type input struct {
	err error
	fn  Fn
}

// want represents wants.
type want struct {
	err  error
	fn   Fn
	help error
}

// errWriter represents [io.Writer] that always fails. Used for
// err return of [Sub.Fprint].
type errWriter struct{}

// resetRejVal represents [flag.Value] that accepts non-empty Set
// values and stores them, but returns err on Set with empty string
// (the DefValue for fs.Var registered with this type). Used to drive
// reset err path in [Sub.RouteArgs] after a prior parse has moved
// the flag value away from its default.
type resetRejVal struct {
	err error
	val string
}

// panicVal represents [flag.Value] to panic. Used for deferred
// [flag.FlagSet.SetOutput] restore in [Sub.RouteArgs].
type panicVal struct{}

var mock = input{
	err: fmt.Errorf("mock"),
	fn:  Fn(nil),
}

func (rej *resetRejVal) String() string {
	if rej == nil {
		return ""
	}

	return rej.val
}

func (rej *resetRejVal) Set(val string) error {
	if val == "" {
		return rej.err
	}

	rej.val = val

	return nil
}

func (panicVal) String() string {
	return ""
}

func (panicVal) Set(string) error {
	panic("panic")
}

// Write is the [errWriter] writer.
func (errWriter) Write([]byte) (int, error) {
	return 0, mock.err
}

// resetFlags replaces global [flag.CommandLine] and [os.Args] for
// tests without leaking state to other tests. Restores prev vals.
func resetFlags(t *testing.T, args []string) {
	t.Helper()

	var (
		prevArgs = os.Args
		prevFlag = flag.CommandLine
	)

	os.Args = append([]string{t.Name()}, args...)

	flag.CommandLine = flag.NewFlagSet(t.Name(), flag.ContinueOnError)

	t.Cleanup(func() {
		os.Args = prevArgs
		flag.CommandLine = prevFlag
	})

	flag.Parse()
}

// newCmd is constructor for route table tests.
func newCmd(min int, sub Sub) *Cmd {
	return &Cmd{
		Help: mock.err,
		Min:  min,
		Sub:  sub,
		Fn:   mock.fn,
	}
}

func TestRouteArgs(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		args []string
		sub  Sub
		want *want
	}{
		"none": {
			args: nil,

			sub: Sub{},

			want: &want{
				err: ErrArgNone,
			},
		},

		"single": {
			args: []string{"mock"},

			sub: Sub{
				"mock": {
					Help: mock.err,
					Fn:   mock.fn,
				},
			},

			want: &want{
				err:  nil,
				fn:   mock.fn,
				help: mock.err,
			},
		},

		"single-under": {
			args: []string{"mock"},

			sub: Sub{
				"mock": newCmd(
					1,

					Sub{},
				),
			},

			want: &want{
				err:  mock.err,
				fn:   mock.fn,
				help: mock.err,
			},
		},

		"nested": {
			args: []string{"mock-a", "mock-b"},

			sub: Sub{
				"mock-a": newCmd(
					1,

					Sub{
						"mock-b": {
							Help: mock.err,
							Min:  0,
							Fn:   mock.fn,
						},
					},
				),
			},

			want: &want{
				err:  nil,
				fn:   mock.fn,
				help: mock.err,
			},
		},

		"nested-under": {
			args: []string{"mock-a", "mock-b"},

			sub: Sub{
				"mock-a": newCmd(
					1,

					Sub{
						"mock-b": {
							Help: mock.err,
							Min:  1,
							Fn:   mock.fn,
						},
					},
				),
			},

			want: &want{
				err:  mock.err,
				fn:   mock.fn,
				help: mock.err,
			},
		},
	}

	for name, test := range tests {
		name, test := name, test

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, _, err := test.sub.RouteArgs(test.args)
			if !errors.Is(err, test.want.err) {
				t.Errorf("unexpected route error %v, want %v", err, test.want.err)
			}

			if got == nil {
				return
			}

			if fmt.Sprintf("%v", got.Fn) != fmt.Sprintf("%v", test.want.fn) {
				t.Errorf("unexpected route fn %v, want %v", got.Fn, test.want.fn)
			}

			if !errors.Is(got.Help, test.want.help) {
				t.Errorf("unexpected route help %v, want %v", got.Help, test.want.help)
			}
		})
	}
}

func TestRouteArgsErrInvalid(t *testing.T) {
	t.Parallel()

	input := Sub{
		"actions": newCmd(
			1,

			Sub{
				"get": {
					Help: mock.err,
					Min:  1,
					Fn:   mock.fn,
				},
			},
		),
	}

	if _, _, got := input.RouteArgs([]string{t.Name()}); !errors.Is(got, ErrArgInvalid) {
		t.Errorf("unexpected route error %v, want %v", got, ErrArgInvalid)
	}
}

func TestRunArgs(t *testing.T) {
	t.Parallel()

	var (
		called bool

		input = Sub{
			"install": {
				Help: mock.err,

				Fn: func(context.Context, []string) error {
					called = true
					return nil
				},
			},
		}

		got = RunArgs(context.Background(), input, []string{"install"})
	)

	if got != nil {
		t.Errorf("run error %v", got)
	}

	if !called {
		t.Errorf("expected fn invoke")
	}
}

func TestRunArgsRouteErr(t *testing.T) {
	t.Parallel()

	if got := RunArgs(context.Background(), Sub{}, nil); !errors.Is(got, ErrArgNone) {
		t.Errorf("unexpected run error %v, want %v", got, ErrArgNone)
	}
}

func TestPrint(t *testing.T) {
	input := Sub{
		"1": {
			Help: errors.New("desc-1"),
		},

		"2": {
			Help: errors.New("desc-2"),
		},
	}

	if err := input.Print(); err != nil {
		t.Errorf("print error %v", err)
	}
}

func TestFprint(t *testing.T) {
	t.Parallel()

	var (
		input = Sub{
			"2": {
				Help: errors.New("desc-2"),
			},

			"1": {
				Help: errors.New("desc-1"),
			},
		}

		buf bytes.Buffer

		want = "1\ndesc-1\n\n2\ndesc-2\n"
	)

	if err := input.Fprint(&buf); err != nil {
		t.Fatalf("fprint error %v", err)
	}

	if got := buf.String(); got != want {
		t.Errorf("unexpected fprint output\ngot:  %q\nwant: %q", got, want)
	}
}

// TestRouteArgsEmptySubLeaf asserts non nil empty Sub does not turn
// leaf cmd into non leaf. Fn must be reachable and Min checked with
// trailing args.
func TestRouteArgsEmptySubLeaf(t *testing.T) {
	t.Parallel()

	var (
		called bool

		input = Sub{
			"mock": {
				Help: mock.err,
				Sub:  Sub{},
				Fn: func(context.Context, []string) error {
					called = true
					return nil
				},
			},
		}
	)

	if err := RunArgs(context.Background(), input, []string{"mock"}); err != nil {
		t.Fatalf("unexpected run error %v", err)
	}

	if !called {
		t.Errorf("expected fn invoke")
	}
}

// TestRouteArgsEmptyArgs asserts empty arg return [ErrArgNone] without i/o.
func TestRouteArgsEmptyArgs(t *testing.T) {
	t.Parallel()

	input := Sub{
		"mock": {Help: mock.err},
	}

	if _, _, err := input.RouteArgs(nil); !errors.Is(err, ErrArgNone) {
		t.Errorf("unexpected route error %v, want %v", err, ErrArgNone)
	}
}

// TestRunArgsNilFn asserts Run does not panic on bad leaf with nil
// Fn and instead return [ErrCmdNoFn].
func TestRunArgsNilFn(t *testing.T) {
	t.Parallel()

	var (
		input = Sub{
			"broken": {
				Help: mock.err,
				Fn:   nil,
			},
		}

		got = RunArgs(context.Background(), input, []string{"broken"})
	)

	if !errors.Is(got, ErrCmdNoFn) {
		t.Errorf("unexpected run error %v, want %v", got, ErrCmdNoFn)
	}
}

// TestRouteArgsFlags asserts per cmd flags are parsed and removed
// from trailing args before returned to caller.
func TestRouteArgsFlags(t *testing.T) {
	t.Parallel()

	var (
		flags   = flag.NewFlagSet(t.Name(), flag.ContinueOnError)
		verbose = flags.Bool("v", false, "verbose")

		input = Sub{
			t.Name(): {
				Help:  mock.err,
				Flags: flags,
				Fn:    mock.fn,
			},
		}
	)

	got, args, err := input.RouteArgs([]string{t.Name(), "-v", "a"})
	if err != nil {
		t.Fatalf("route error %v", err)
	}

	if got == nil {
		t.Fatal("expected cmd, got nil")
	}

	if !*verbose {
		t.Errorf("expected -v to be parsed, verbose=%v", *verbose)
	}

	if len(args) != 1 || args[0] != "a" {
		t.Errorf("unexpected trailing args %v, want [a]", args)
	}
}

// TestRouteArgsNoMutation asserts that routing does not mutate cmd
// vals across successive RouteArgs calls.
func TestRouteArgsNoMutation(t *testing.T) {
	t.Parallel()

	var (
		leaf = &Cmd{
			Help: mock.err,
			Min:  1,
			Sub:  Sub{},
			Fn:   mock.fn,
		}

		sub = Sub{
			"leaf": leaf,
		}

		snap = *leaf
	)

	if _, _, err := sub.RouteArgs([]string{"leaf", "a"}); err != nil {
		t.Fatalf("unexpected route error %v", err)
	}

	if leaf.Min != snap.Min ||
		!errors.Is(leaf.Help, snap.Help) {
		t.Errorf("cmd changed: got %+v, want %+v", *leaf, snap)
	}

	if reflect.ValueOf(leaf.Sub).Pointer() != reflect.ValueOf(snap.Sub).Pointer() {
		t.Errorf("sub changed: got %v, want %v", leaf.Sub, snap.Sub)
	}
}

// TestRouteArgsSubMiss asserts that descent stops when next arg is
// unknown. Token return wrapped [ErrArgInvalid] so callers can
// distinguish typo from usage req.
func TestRouteArgsSubMiss(t *testing.T) {
	t.Parallel()

	input := Sub{
		"mock-a": newCmd(
			0,

			Sub{
				"mock-b": {
					Help: mock.err,
					Fn:   mock.fn,
				},
			},
		),
	}

	_, _, err := input.RouteArgs([]string{"mock-a", "mock-c"})

	if !errors.Is(err, ErrArgInvalid) {
		t.Errorf("unexpected route error %v, want %v", err, ErrArgInvalid)
	}

	if err == nil || !strings.Contains(err.Error(), "mock-c") {
		t.Errorf("expected error to mention offending token %q, got %v", "mock-c", err)
	}
}

// TestRouteArgsFlagsErr asserts err from [flag.FlagSet.Parse] is
// wrapped with [ErrFlagParse].
func TestRouteArgsFlagsErr(t *testing.T) {
	t.Parallel()

	flags := flag.NewFlagSet(t.Name(), flag.ContinueOnError)
	flags.SetOutput(&bytes.Buffer{})
	flags.Bool("v", false, "verbose")

	input := Sub{
		t.Name(): {
			Help:  mock.err,
			Flags: flags,
			Fn:    mock.fn,
		},
	}

	_, _, got := input.RouteArgs([]string{t.Name(), "-no"})

	if got == nil {
		t.Fatalf("unexpected flag parse success")
	}

	if !errors.Is(got, ErrFlagParse) {
		t.Errorf("unexpected route error %v, want %v", got, ErrFlagParse)
	}
}

// TestRouteArgsFlagsResetErr asserts err returned by custom [flag.Value]
// Set pre parse reset returns [ErrFlagReset] wrapped with flag name and
// underlying Set err.
func TestRouteArgsFlagsResetErr(t *testing.T) {
	t.Parallel()

	var (
		flags = flag.NewFlagSet(t.Name(), flag.ContinueOnError)
		rej   = &resetRejVal{err: mock.err}
	)

	flags.SetOutput(&bytes.Buffer{})
	flags.Var(rej, "tag", "tag")
	flags.String("zone", "", "zone")

	input := Sub{
		t.Name(): {
			Help:  mock.err,
			Flags: flags,
			Fn:    mock.fn,
		},
	}

	// Dirty flag successful first parse so subsequent route attempts
	// to reset back to DefValue.
	if _, _, err := input.RouteArgs([]string{t.Name(), "-tag=x"}); err != nil {
		t.Fatalf("unexpected first route error: %v", err)
	}

	_, _, got := input.RouteArgs([]string{t.Name(), "a"})
	if got == nil {
		t.Fatalf("unexpected flag reset success")
	}

	if !errors.Is(got, ErrFlagReset) {
		t.Errorf("unexpected route error %v, want %v", got, ErrFlagReset)
	}

	if !errors.Is(got, mock.err) {
		t.Errorf("expected underlying Set error to be wrapped, got %v", got)
	}

	if !strings.Contains(got.Error(), "tag") {
		t.Errorf("expected error to mention flag name %q, got %v", "tag", got)
	}
}

// TestRouteArgsFlagsReset asserts flag vals left over from prev
// RouteArgs are reset to defaults before next Parse.
func TestRouteArgsFlagsReset(t *testing.T) {
	t.Parallel()

	var (
		flags   = flag.NewFlagSet(t.Name(), flag.ContinueOnError)
		verbose = flags.Bool("v", false, "verbose")
		name    = flags.String("name", "default", "name")

		input = Sub{
			t.Name(): {
				Help:  mock.err,
				Flags: flags,
				Fn:    mock.fn,
			},
		}
	)

	if _, _, got := input.RouteArgs([]string{t.Name(), "-v", "-name=alt", "a"}); got != nil {
		t.Fatalf("unexpected route error %v", got)
	}

	if !*verbose || *name != "alt" {
		t.Fatalf("expected set flags, got verbose=%v name=%q", *verbose, *name)
	}

	if _, _, got := input.RouteArgs([]string{t.Name(), "a"}); got != nil {
		t.Fatalf("unexpected route error %v", got)
	}

	if *verbose {
		t.Errorf("expected -v to reset to false, got true")
	}

	if *name != "default" {
		t.Errorf("expected -name to reset to %q, got %q", "default", *name)
	}
}

// TestRouteArgsFlagsHelp asserts -h / --help on per cmd
// FlagSet returns cmd Help err instead of [flag.ErrHelp].
func TestRouteArgsFlagsHelp(t *testing.T) {
	t.Parallel()

	flags := flag.NewFlagSet(t.Name(), flag.ContinueOnError)
	flags.SetOutput(&bytes.Buffer{})
	flags.Bool("v", false, "verbose")

	input := Sub{
		t.Name(): {
			Help:  mock.err,
			Flags: flags,
			Fn:    mock.fn,
		},
	}

	_, _, got := input.RouteArgs([]string{t.Name(), "-h"})

	if !errors.Is(got, mock.err) {
		t.Errorf("unexpected route error %v, want %v", got, mock.err)
	}

	if errors.Is(got, ErrFlagParse) {
		t.Errorf("flag.ErrHelp should not be wrapped with ErrFlagParse, got %v", got)
	}
}

// TestRouteArgsFlagsHelpNoHelp asserts -h on cmd with nil Help
// returns [ErrCmdNoHelp] instead of [flag.ErrHelp].
func TestRouteArgsFlagsHelpNoHelp(t *testing.T) {
	t.Parallel()

	flags := flag.NewFlagSet(t.Name(), flag.ContinueOnError)
	flags.SetOutput(&bytes.Buffer{})
	flags.Bool("v", false, "verbose")

	input := Sub{
		t.Name(): {
			Flags: flags,
			Fn:    mock.fn,
		},
	}

	if _, _, got := input.RouteArgs([]string{t.Name(), "-h"}); !errors.Is(got, ErrCmdNoHelp) {
		t.Errorf("unexpected route error %v, want %v", got, ErrCmdNoHelp)
	}
}

// TestRouteArgsFlagsSilent asserts RouteArgs performs no i/o when
// [flag.FlagSet.Parse] would otherwise emit usage to fs.Output().
func TestRouteArgsFlagsSilent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
	}{
		{
			"parse error", []string{t.Name(), "-err"},
		},

		{
			"help flag", []string{t.Name(), "-h"},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var (
				flags = flag.NewFlagSet(t.Name(), flag.ContinueOnError)
				out   = &bytes.Buffer{}
			)

			flags.SetOutput(out)
			flags.Bool("v", false, "verbose")

			input := Sub{
				t.Name(): {
					Help:  mock.err,
					Flags: flags,
					Fn:    mock.fn,
				},
			}

			if _, _, err := input.RouteArgs(test.args); err == nil {
				t.Fatalf("unexpected route success")
			}

			if out.Len() != 0 {
				t.Errorf("wrote %q output", out.String())
			}

			if flags.Output() != out {
				t.Errorf("left flagset output as %v, want caller's writer restored", flags.Output())
			}
		})
	}
}

// TestRouteArgsFlagsExitOnErrorDowngrade asserts FlagSet with
// [flag.ExitOnError] is downgraded to [flag.ContinueOnError] before
// Parse, so parse err or -h / --help returns routing err instead of
// terminate proc, and the prior error handling is restored after
// Parse so the caller-owned FlagSet is not silently mutated.
func TestRouteArgsFlagsExitOnErrorDowngrade(t *testing.T) {
	t.Parallel()

	flags := flag.NewFlagSet(t.Name(), flag.ExitOnError)
	flags.SetOutput(&bytes.Buffer{})
	flags.Bool("v", false, "verbose")

	input := Sub{
		t.Name(): {
			Help:  mock.err,
			Flags: flags,
			Fn:    mock.fn,
		},
	}

	if _, _, got := input.RouteArgs([]string{t.Name(), "-nope"}); !errors.Is(got, ErrFlagParse) {
		t.Errorf("unexpected route error %v, want %v", got, ErrFlagParse)
	}

	if eh := flags.ErrorHandling(); eh != flag.ExitOnError {
		t.Errorf("flagset ErrorHandling = %v, want %v restored after parse err", eh, flag.ExitOnError)
	}

	if _, _, got := input.RouteArgs([]string{t.Name(), "-h"}); !errors.Is(got, mock.err) {
		t.Errorf("unexpected route error %v, want %v", got, mock.err)
	}

	if eh := flags.ErrorHandling(); eh != flag.ExitOnError {
		t.Errorf("flagset ErrorHandling = %v, want %v restored after help", eh, flag.ExitOnError)
	}
}

// TestRouteArgsFlagsResetSkipDefault asserts that the pre parse reset
// is skipped for any flag whose Value already equals DefValue, so a
// custom flag.Value whose Set rejects its own DefValue does not fail
// on the first invocation before user input is parsed.
func TestRouteArgsFlagsResetSkipDefault(t *testing.T) {
	t.Parallel()

	var (
		flags = flag.NewFlagSet(t.Name(), flag.ContinueOnError)
		rej   = &resetRejVal{err: mock.err}
	)

	flags.SetOutput(&bytes.Buffer{})
	flags.Var(rej, "tag", "tag")

	input := Sub{
		t.Name(): {
			Help:  mock.err,
			Flags: flags,
			Fn:    mock.fn,
		},
	}

	if _, _, err := input.RouteArgs([]string{t.Name(), "-tag=x"}); err != nil {
		t.Fatalf("unexpected first route error: %v", err)
	}
}

// TestRouteArgsFlagsRestoreOnPanic asserts that caller FlagSet
// output writer is restored even if Parse panics, so to not
// silently pin FlagSet output to [io.Discard].
func TestRouteArgsFlagsRestoreOnPanic(t *testing.T) {
	t.Parallel()

	var (
		out   = &bytes.Buffer{}
		flags = flag.NewFlagSet(t.Name(), flag.ContinueOnError)
	)

	flags.SetOutput(out)
	flags.Var(panicVal{}, "panic", "panics on Set")

	input := Sub{
		t.Name(): {
			Help:  mock.err,
			Flags: flags,
			Fn:    mock.fn,
		},
	}

	defer func() {
		if ok := recover(); ok == nil {
			t.Fatalf("expected panic")
		}

		if flags.Output() != out {
			t.Errorf("router left flagset output %v after panic, want writer restored", flags.Output())
		}
	}()

	_, _, _ = input.RouteArgs([]string{t.Name(), "-panic=x"})
}

// TestFprintErr asserts write err propagates to caller.
func TestFprintErr(t *testing.T) {
	t.Parallel()

	input := Sub{
		t.Name(): {
			Help: mock.err,
		},
	}

	if got := input.Fprint(errWriter{}); !errors.Is(got, mock.err) {
		t.Errorf("unexpected fprint error %v, want %v", got, mock.err)
	}
}

// TestRouteArgsNilCmd asserts nil [*Cmd] in sub return [ErrCmdNil]
// instead of panic.
func TestRouteArgsNilCmd(t *testing.T) {
	t.Parallel()

	input := Sub{t.Name(): nil}

	if _, _, got := input.RouteArgs([]string{t.Name()}); !errors.Is(got, ErrCmdNil) {
		t.Errorf("unexpected route error %v, want %v", got, ErrCmdNil)
	}
}

// TestRouteArgsNilSubCmd asserts nil [*Cmd] during descent into
// parent Sub return [ErrCmdNil] instead of panic.
func TestRouteArgsNilSubCmd(t *testing.T) {
	t.Parallel()

	input := Sub{
		"parent": newCmd(
			0,

			Sub{
				"child": nil,
			},
		),
	}

	if _, _, got := input.RouteArgs([]string{"parent", "child"}); !errors.Is(got, ErrCmdNil) {
		t.Errorf("unexpected route error %v, want %v", got, ErrCmdNil)
	}
}

// TestRouteArgsNilHelpNonLeaf asserts non leaf with nil Help
// returns [ErrCmdNoHelp] instead of nil.
func TestRouteArgsNilHelpNonLeaf(t *testing.T) {
	t.Parallel()

	input := Sub{
		"parent": {
			Sub: Sub{
				"child": {
					Help: mock.err,
					Fn:   mock.fn,
				},
			},
		},
	}

	if _, _, got := input.RouteArgs([]string{"parent"}); !errors.Is(got, ErrCmdNoHelp) {
		t.Errorf("unexpected route error %v, want %v", got, ErrCmdNoHelp)
	}
}

// TestRouteArgsNilHelpUnderMin asserts leaf with nil Help with
// too few args return [ErrCmdNoHelp] instead of nil.
func TestRouteArgsNilHelpUnderMin(t *testing.T) {
	t.Parallel()

	input := Sub{
		t.Name(): {
			Min: 1,
			Fn:  mock.fn,
		},
	}

	if _, _, got := input.RouteArgs([]string{t.Name()}); !errors.Is(got, ErrCmdNoHelp) {
		t.Errorf("unexpected route error %v, want %v", got, ErrCmdNoHelp)
	}
}

// TestRunArgsNilFnNilHelp asserts leaf with both Fn and Help
// nil return [ErrCmdNoFn] joined with [ErrCmdNoHelp] instead of panic.
func TestRunArgsNilFnNilHelp(t *testing.T) {
	t.Parallel()

	var (
		input = Sub{t.Name(): {}}

		got = RunArgs(
			context.Background(),
			input,
			[]string{t.Name()},
		)
	)

	if !errors.Is(got, ErrCmdNoFn) {
		t.Errorf("unexpected run error %v, want %v", got, ErrCmdNoFn)
	}

	if !errors.Is(got, ErrCmdNoHelp) {
		t.Errorf("unexpected run error %v, want %v", got, ErrCmdNoHelp)
	}
}

// TestFprintNilCmd asserts nil [*Cmd] in sub causes Fprint to
// return [ErrCmdNil] instead of panic.
func TestFprintNilCmd(t *testing.T) {
	t.Parallel()

	input := Sub{t.Name(): nil}

	if err := input.Fprint(&bytes.Buffer{}); !errors.Is(err, ErrCmdNil) {
		t.Errorf("unexpected fprint error %v, want %v", err, ErrCmdNil)
	}
}

// TestFprintNilHelp asserts cmd with nil Help shows placeholder
// instead of "<nil>".
func TestFprintNilHelp(t *testing.T) {
	t.Parallel()

	var (
		input = Sub{
			"mock": {},
		}

		buf  bytes.Buffer
		want = "mock\n(no help)\n"
	)

	if err := input.Fprint(&buf); err != nil {
		t.Fatalf("fprint error %v", err)
	}

	if got := buf.String(); got != want {
		t.Errorf("unexpected fprint output\ngot:  %q\nwant: %q", got, want)
	}
}

// TestRun asserts [Run] dispatches to [RunArgs] using left over
// [flag.Parse] args.
func TestRun(t *testing.T) {
	resetFlags(t, []string{"install"})

	var (
		called bool

		input = Sub{
			"install": {
				Help: mock.err,

				Fn: func(context.Context, []string) error {
					called = true
					return nil
				},
			},
		}
	)

	if err := Run(context.Background(), input); err != nil {
		t.Errorf("run error %v", err)
	}

	if !called {
		t.Errorf("expected fn invoke")
	}
}

// TestRoute asserts [Sub.Route] resolves [flag.Args] to matched
// leaf cmd.
func TestRoute(t *testing.T) {
	resetFlags(t, []string{t.Name()})

	input := Sub{
		t.Name(): {
			Help: mock.err,
			Fn:   mock.fn,
		},
	}

	got, _, err := input.Route()
	if err != nil {
		t.Fatalf("route error %v", err)
	}

	if got == nil {
		t.Fatal("expected resolved cmd, got nil")
	}

	if !errors.Is(got.Help, mock.err) {
		t.Errorf("unexpected route help %v, want %v", got.Help, mock.err)
	}
}
