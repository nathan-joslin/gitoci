package actions

import "fmt"

// Capability defines a git-remote-helper capability.
//
// See https://git-scm.com/docs/gitremote-helpers#_capabilities.
type Capability = string

const (
	CapOption Capability = "option"
	CapFoo    Capability = "foo"
	CapBar    Capability = "bar"
	// CapPush   Capability = "push"
	// CapFetch             = "fetch"
)

func (action *GitOCI) capabilities() error {
	// TODO: another method, we don't want to update this all the time...
	capabilities := []Capability{CapOption, CapFoo, CapBar}
	if err := action.batcher.WriteBatch(capabilities...); err != nil {
		return fmt.Errorf("writing capabilities: %w", err)
	}
	return nil
}
