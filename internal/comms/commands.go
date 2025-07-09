package comms

import "errors"

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
