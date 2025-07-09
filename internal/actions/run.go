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
		if err := action.capabilities(); err != nil {
			return fmt.Errorf("responding to capabilities command: %w", err)
		}
	}

	cmds, err := action.batcher.ReadBatch(ctx)
	if err != nil {
		return fmt.Errorf("reading bach input: %w", err)
	}

	for _, cmd := range cmds {
		switch cmd.CommandType {
		case comms.CmdEmpty:
			// TODO: revise: some commands aren't exactly terminated by the blank line
			// but we can still predict the end, i.e. the rest of the stream is all that's left
			slog.InfoContext(ctx, "batch complete")
			return nil
		default:
			return fmt.Errorf("")
		}
	}

	// TODO: Next command is 'list', can be read in a batch
	slog.InfoContext(ctx, "reading batch")
	_, err = action.batcher.ReadBatch(ctx)
	if err != nil {
		return fmt.Errorf("reading batch input: %w", err)
	}

	return fmt.Errorf("not implemented")
}
