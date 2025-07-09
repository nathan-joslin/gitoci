package actions

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/act3-ai/gitoci/internal/comms"
)

// GitOCI represents the base action
type GitOCI struct {
	// TODO: Could be dangerous when storing in struct like this... mutex?
	batcher comms.BatchReadWriter

	// local repository
	gitDir string

	// OCI remote
	name   string // may have same value as address
	addess string

	Option

	version string
}

// NewGitOCI creates a new Tool with default values
func NewGitOCI(in io.Reader, out io.Writer, gitDir, shortname, address, version string) *GitOCI {
	return &GitOCI{
		batcher: comms.NewBatcher(in, out),
		gitDir:  gitDir,
		name:    shortname,
		addess:  address,
		version: version,
	}
}

// Runs the Hello action
func (action *GitOCI) Run(ctx context.Context) error {
	// first command is always "capabilities"
	cmd, err := action.batcher.Read(ctx)
	switch {
	case err != nil:
		return fmt.Errorf("reading initial command: %w", err)
	case cmd.CommandType != comms.CmdCapabilities:
		return fmt.Errorf("unexpected first command %s, expected 'capabilities'", cmd.CommandType)
	default:
		slog.InfoContext(ctx, "writing capabilities")
		if err := action.capabilities(ctx); err != nil {
			return err
		}
	}

	var done bool
	for !done {
		cmd, err := action.batcher.Read(ctx)
		if err != nil {
			return fmt.Errorf("reading next line: %w", err)
		}
		slog.InfoContext(ctx, "read command from Git", "command", cmd.CommandType, "data", fmt.Sprintf("%v", cmd.Data))
		switch cmd.CommandType {
		case comms.CmdEmpty:
			done = true
		case comms.CmdCapabilities:
			// Git shouldn't need to do this again, but let's be safe
			if err := action.capabilities(ctx); err != nil {
				return err
			}
		case comms.CmdOption:
			if err := action.option(ctx, cmd); err != nil {
				return err
			}
		}
	}

	// // TODO: Next command is 'list', can be read in a batch
	// slog.InfoContext(ctx, "reading batch")
	// _, err = action.batcher.ReadBatch(ctx)
	// if err != nil {
	// 	return fmt.Errorf("reading batch input: %w", err)
	// }

	return fmt.Errorf("not implemented")
}
