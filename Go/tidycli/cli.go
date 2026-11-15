// © Roscoe Skeens <web@rskeens.com>
// SPDX-License-Identifier: MIT

// Package tidycli is a tiny and dependency free router for clis that need
// subcommands without adopting a framework. A minimal positional subcommand
// router rather than bloat.
//
// tidycli routes positional subcommands in under a couple hundred lines of
// code and stays out of the way of the standard library's [flag] package.
//
// Prettier than stdlib, simpler than urfave/cli or kong and less heavy than
// Cobra. A core focus on simplicity, excellent defaults, clean help output,
// tiny api surface and fast setup.
//
// A cmd is described as a tree of [Cmd] values keyed by arg in a [Sub].
// Each leaf cmd is invoked through a [Fn], which receives [context.Context]
// and trailing arg slice. Every value in a [Sub] must be non nil [*Cmd].
//
// Typical use:
//
//	import (
//	    "context"
//	    "errors"
//	    "flag"
//	    "fmt"
//	    "log"
//
//	    "github.com/rskeens/tidycli"
//	)
//
//	cmds := tidycli.Sub{
//	    "hello": {
//	        Help: errors.New("hello NAME"),
//	        Min:  1,
//	        Fn: func(ctx context.Context, args []string) error {
//	            fmt.Println("hello,", args[0])
//	            return nil
//	        },
//	    },
//	}
//
//	flag.Parse()
//	if err := tidycli.Run(context.Background(), cmds); err != nil {
//	    log.Fatal(err)
//	}
//
// Per cmd flags are supported by attaching [flag.FlagSet] to a leaf [Cmd]
// via [Cmd.Flags]. That is parsed against trailing args before invoking
// [Cmd.Fn]. [Cmd.Flags] on non leaf commands are ignored.
//
// tidycli does not bootstrap, configure logging or parse global flags.
// Callers wire those concerns to fit their needs.
package tidycli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
)

// Fn is cmd func. Optional [context.Context] for cancel and deadline.
// Args follows cmd keyword per cmd flag parse.
type Fn func(ctx context.Context, args []string) error

// Cmd represents single cmd in the cmd tree.
//
// Help is returned as err if cmd is invoked with fewer than min args or when
// cmd is non leaf (cmd with sub holds at least one entry) and descent stops
// there because no further args. Unknown subcmd at returns [ErrArgInvalid]
// instead of Help. The err is intentional, it lets callers propagate usage
// string.
//
// Sub, when non empty, makes this cmd a parent. Routing descends it when next
// arg matches key. Nil sub and non nil but empty sub are both treated as a
// leaf. Empty Sub{} on leaf cmd does not turn into parent.
//
// Fn, Min, and Flags only effect leaf cmds. If Sub is non empty then routing
// either descends into child or returns err (Help when descent stops at cmd,
// [ErrArgInvalid] when the next arg does not match child key), so parent's Fn
// is never invoked, its Min is never checked and the Flags are never parsed.
// Setting Fn, Min, or Flags alongside a non empty Sub is silently ignored,
// attach them to leaf instead.
//
// Flags, when non nil and set on leaf are parsed against trailing args when
// routing has resolved cmd. Remaining positional args are passed to Fn and
// counted against Min. Cmd vals are not mutated by routing. Flags is owned by
// the caller and may carry parsed vals after Run / Route returns.
//
// Flags must use [flag.ContinueOnError]. Routing relies on Parse returning
// errs to convert to Help. [Sub.RouteArgs] re applies [flag.ContinueOnError]
// via [flag.FlagSet.Init] before parsing, then restores prior error
// handling after Parse. FlagSet with [flag.ExitOnError] or
// [flag.PanicOnError] are downgraded for the duration of Parse rather than
// allowed to terminate proc.
//
// Fn invoked for leaf cmds after routing succeeds and arg count satisfies Min.
type Cmd struct {
	Help  error
	Min   int
	Sub   Sub
	Flags *flag.FlagSet
	Fn    Fn
}

// Sub represents arg keyword to cmd. Vals may not be nil. Routing treats nil
// [*Cmd] as [ErrCmdNil] without panic.
type Sub map[string]*Cmd

// Run routes parsed args by stdlib [flag] package and on success invokes resolved
// cmd [Cmd.Fn] with ctx and trailing args.
//
// Thin wrapper around [RunArgs] that supplies [flag.Args], use [RunArgs] when
// integrating with alternative flag parser, embedded shell or test harness for
// custom arg slice.
//
// [Cmd] with nil [Cmd.Fn] returns [Cmd.Help] err without panic wrapped with [ErrCmdNoFn].
func Run(ctx context.Context, sub Sub) error {
	return RunArgs(ctx, sub, flag.Args())
}

// RunArgs is [Run] entry for given arg slice instead of [flag.Args].
// Identical routing and [Cmd.Fn] invocation as [Run].
func RunArgs(ctx context.Context, sub Sub, args []string) error {
	cmd, rest, err := sub.RouteArgs(args)
	if err != nil {
		return err
	}

	if cmd.Fn == nil {
		return errors.Join(ErrCmdNoFn, help(cmd))
	}

	return cmd.Fn(ctx, rest)
}

// Route resolves [flag.Args] against sub and returns leaf cmd with trailing
// positional args. Thin wrapper around [Sub.RouteArgs].
func (sub Sub) Route() (*Cmd, []string, error) {
	return sub.RouteArgs(flag.Args())
}

// RouteArgs resolves args against sub and returns leaf cmd with trailing
// positional arguments. Returns [Cmd.Help] err when cmd is non leaf with no further
// args to consume or when fewer than [Cmd.Min] trailing args. Returns [ErrArgNone]
// if args is empty and [ErrArgInvalid] when arg at any depth does not match
// cmd. If cmd has non nil [Cmd.Flags] then RouteArgs resets each flag to
// [flag.Flag.DefValue] and parses against trailing args.
//
// Reset relies on clean Value.Set(DefValue), which is reliable for stdlib
// [flag] scalar flag types but may not fully reset custom [flag.Value] whose
// Set accumulates state. Such flags should construct a new [flag.FlagSet]
// per invocation rather than rely on reset. Reset is skipped for any flag
// whose current [flag.Value.String] already equals [flag.Flag.DefValue], so
// custom Values that reject their own DefValue do not fail on the first
// invocation before user input is parsed. If [flag.Value.Set] returns err
// during reset then RouteArgs returns [ErrFlagReset] rather than silently
// leaking stale state.
//
// A [flag.ErrHelp] result from the user supplying -h or --help is reported
// as cmd [Cmd.Help] err. Other errs from [flag.FlagSet.Parse] are wrapped with
// [ErrFlagParse].
//
// Before parsing, RouteArgs re applies [flag.ContinueOnError] to FlagSet via
// [flag.FlagSet.Init]: a FlagSet constructed with [flag.ExitOnError] or
// [flag.PanicOnError] would otherwise terminate the process from inside Parse
// on parse err or on -h / --help, bypassing RouteArgs's err reporting. The
// prior error handling is captured and restored after Parse so RouteArgs
// does not leave a silent side effect on the caller-owned FlagSet.
//
// RouteArgs does not mutate [Cmd] other than vals stored inside [Cmd.Flags]
// (caller owned) and does not read from global [flag.CommandLine]. Safe for
// alongside other flag parsers, embedded shells and parallel tests that don't
// share [Cmd.Flags] across goroutines. Performs no i/o of its own so callers
// are responsible for rendering errs (including [ErrArgNone]).
//
// Malformed cmd trees report as errs without panic. Nil [*Cmd] in sub or
// descendant [Sub] returns [ErrCmdNil] and routed cmd with [Cmd.Help] as nil when
// [Cmd.Help] would be returned returns [ErrCmdNoHelp].
func (sub Sub) RouteArgs(args []string) (*Cmd, []string, error) {
	if len(args) == 0 {
		return nil, nil, ErrArgNone
	}

	cmd, ok := sub[args[0]]
	if !ok {
		return nil, nil, fmt.Errorf("%w: %v", ErrArgInvalid, args[0])
	}

	if cmd == nil {
		return nil, nil, fmt.Errorf("%w: %v", ErrCmdNil, args[0])
	}

	rest := args[1:]

	for len(cmd.Sub) > 0 && len(rest) > 0 {
		next, ok := cmd.Sub[rest[0]]
		if !ok {
			// Distinguish typo from -h or empty descent.
			return nil, nil, fmt.Errorf("%w: %v", ErrArgInvalid, rest[0])
		}

		if next == nil {
			return nil, nil, fmt.Errorf("%w: %v", ErrCmdNil, rest[0])
		}

		cmd = next
		rest = rest[1:]
	}

	if len(cmd.Sub) > 0 {
		return nil, nil, help(cmd)
	}

	if cmd.Flags != nil {
		// Reset each registered flag to default. Flag set on prev invocation
		// does not silently leak into next one if reused. Round trip through
		// Value.Set is expected to succeed for scalar flag types provided by
		// stdlib. Custom flag.Value with Set rejects DefValue or accumulates
		// state can still fail as [ErrFlagReset].
		//
		// Skip flags if current Value already equals DefValue. On first
		// invocation no prior parse has dirty state so reset is noop for well
		// behaved Values, and custom flag.Value whose Set rejects its DefValue
		// would fail before user input parse.
		var (
			errReset  error
			resetName string
		)

		cmd.Flags.VisitAll(func(f *flag.Flag) {
			if errReset != nil {
				return
			}

			if f.Value.String() == f.DefValue {
				return
			}

			if err := f.Value.Set(f.DefValue); err != nil {
				errReset = err
				resetName = f.Name
			}
		})

		if errReset != nil {
			return nil, nil, fmt.Errorf("%w: %s: %w", ErrFlagReset, resetName, errReset)
		}

		// flag.FlagSet.Parse writes to fs.Output() (stderr by default) on parse
		// err and -h / --help, emitting "Usage of <name>:" plus flag defaults.
		// RouteArgs has no own i/o. Err rendering to callers so silence FlagSet
		// for Parse then restore writer.
		//
		// Restore is deferred for panic from custom flag.Value.Set.
		prevOut := cmd.Flags.Output()
		cmd.Flags.SetOutput(io.Discard)
		defer cmd.Flags.SetOutput(prevOut)

		// Also capture prior error handling and restore it.
		prevErrHandling := cmd.Flags.ErrorHandling()
		cmd.Flags.Init(cmd.Flags.Name(), flag.ContinueOnError)
		defer cmd.Flags.Init(cmd.Flags.Name(), prevErrHandling)

		if err := cmd.Flags.Parse(rest); err != nil {
			// usage request, not a parse failure. Return cmd [Cmd.Help].
			if errors.Is(err, flag.ErrHelp) {
				return nil, nil, help(cmd)
			}

			return nil, nil, fmt.Errorf("%w: %w", ErrFlagParse, err)
		}

		rest = cmd.Flags.Args()
	}

	if len(rest) < cmd.Min {
		return nil, nil, help(cmd)
	}

	return cmd, rest, nil
}

// help returns cmd.Help. Falls back to [ErrCmdNoHelp] if no [Cmd.Help]. Used wherever
// routing would return nil err in place of usage string.
func help(cmd *Cmd) error {
	if cmd.Help == nil {
		return ErrCmdNoHelp
	}

	return cmd.Help
}

// Print writes cmd keywords and Helps to stdout alphabetically. Thin
// wrapper around [Sub.Fprint] targets [os.Stdout]. Redirect of usage view
// should call [Sub.Fprint]. Any logo or trailing footer should print around it.
func (sub Sub) Print() error {
	return sub.Fprint(os.Stdout)
}

// Fprint writes cmd keywords and Helps to wr alphabetically. Each entry
// renders keyword on one line followed by Help on the next then separated
// from following entry by blank line.
//
// Nil [*Cmd] sub returns [ErrCmdNil] without panic. [Cmd] with nil [Cmd.Help] is
// rendered with placeholder line so keyword appears in usage view.
func (sub Sub) Fprint(wr io.Writer) error {
	args := make([]string, 0, len(sub))

	for arg := range sub {
		args = append(args, arg)
	}

	sort.Strings(args)

	for idx, arg := range args {
		var (
			cmd = sub[arg]
			sep = ""
		)

		if cmd == nil {
			return fmt.Errorf("%w: %v", ErrCmdNil, arg)
		}

		if idx > 0 {
			sep = "\n"
		}

		desc := any(cmd.Help)
		if cmd.Help == nil {
			desc = "(no help)"
		}

		if _, err := fmt.Fprintf(wr, "%s%s\n%v\n", sep, arg, desc); err != nil {
			return err
		}
	}

	return nil
}
