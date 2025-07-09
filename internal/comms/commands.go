package comms

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

// Error types.
var (
	ErrUnsupportedCommand = errors.New("unsupported git-remote-helper command")
)

// CommandType is an implemented git-remote-helper command provided by Git.
//
// See https://git-scm.com/docs/gitremote-helpers#_commands.
type CommandType string

const (
	// Git conventions
	CmdCapabilities CommandType = "capabilities"

	// CmdPush                 = "push"
	// CmdFetch                = "fetch"
	// CmdList                 = "list"
	// CmdListForPush       = "for-push"

	// not a Git convention, marks end of input
	CmdEmpty CommandType = "empty"
)

const (
	CmdOption          CommandType = "option"
	CmdOptionVerbosity CommandType = "verbosity"
)

// Command represents a parsed command received from Git.
type Command struct {
	CommandType
	Data []string
}

// SupportedOption returns true if an option is supported.
func SupportedOption(name string) bool {
	option := CommandType(name)

	return (option == CmdOptionVerbosity)
}

// ParseCommand reads a single line received from Git, turning it into a Command
// easily identified by CommandType.
func ParseCommand(ctx context.Context, line string) (Command, error) {
	slog.DebugContext(ctx, "parsing command")
	fields := strings.Fields(line)
	if len(fields) < 1 {
		return Command{
			CommandType: CmdEmpty,
		}, nil
	}

	cmd := CommandType(fields[0])
	switch cmd {
	case CmdCapabilities:
		return Command{
			CommandType: CmdCapabilities,
		}, nil
	case CmdOption:
		// TODO: we should try to not make options fatal, but we may have to
		// make an exception for force (or others).
		if len(fields) != 3 {
			slog.ErrorContext(ctx, "invalid number of arguments to option command", "got", fmt.Sprintf("%d", len(fields)), "want", "3")
			return Command{}, fmt.Errorf("invalid number of args to option command")
		} else {
			return Command{
				CommandType: CmdOption,
				Data:        fields[1:],
			}, nil
		}
	default:
		return Command{}, fmt.Errorf("%w: %s", ErrUnsupportedCommand, cmd)
	}
}
