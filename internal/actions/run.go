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

// func cleanPrefix(address string) string {
// 	return strings.TrimPrefix(address, "oci://")
// }

// Runs the Hello action
func (action *GitOCI) Run(ctx context.Context) error {
	// first command is always "capabilities"
	cmd, err := action.batcher.Read()
	switch {
	case err != nil:
		return fmt.Errorf("reading initial command: %w", err)
	case cmd.CommandType != comms.CmdCapabilities:
	default:
		if err := action.capabilities(); err != nil {
			return fmt.Errorf("responding to capabilities command: %w", err)
		}
	}

	slog.InfoContext(ctx, "reading batch")
	batch, err := action.batcher.ReadBatch()
	if err != nil {
		return fmt.Errorf("reading input batch: %w", err)
	}

	for _, cmd := range batch {
		slog.InfoContext(ctx, "executing command", "command", cmd)
		switch cmd.CommandType {
		case comms.CmdCapabilities:
			action.capabilities()
		default:
			return fmt.Errorf("unsupported command %s", cmd)
		}
	}

	return nil
}
