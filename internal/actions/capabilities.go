package actions

import "fmt"

// Capability defines a git-remote-helper capability.
//
// See https://git-scm.com/docs/gitremote-helpers#_capabilities.
type Capability = string

const (
	CapFoo Capability = "foo"
	CapBar Capability = "bar"
	// CapPush   Capability = "push"
	// CapFetch             = "fetch"
	// CapOption            = "option"
)

// Command defines commands implemented to support capabilities.
//
// See https://git-scm.com/docs/gitremote-helpers#_commands.
type Command = string

const (
	CmdCapabilities Command = "capabilities"
	// CmdPush                 = "push"
	// CmdFetch                = "fetch"
	// CmdOption               = "option"
	// CmdList                 = "list"
	// CmdListForPush          = "for-push"
)

func (action *GitOCI) capabilities() error {
	// TODO: another method, we don't want to update this all the time...
	capabilities := []Capability{CapFoo, CapBar}
	if err := action.batcher.WriteBatch(capabilities...); err != nil {
		return fmt.Errorf("responding to capabilities command: %w", err)
	}
	return nil
}
