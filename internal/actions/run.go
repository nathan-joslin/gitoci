package actions

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/act3-ai/gitoci/internal/interpreter"
)

// GitOCI represents the base action
type GitOCI struct {
	// TODO: Could be dangerous when storing in struct like this... mutex?
	batcher interpreter.Batcher

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
		batcher: interpreter.NewBatcher(in, out),
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
	slog.InfoContext(ctx, "reading batch")
	batch := action.batcher.ReadBatch()

	for _, cmd := range batch {
		slog.InfoContext(ctx, "executing command", "command", cmd)
		switch cmd {
		case CmdCapabilities:
			action.capabilities()
		default:
			return fmt.Errorf("unsupported command %s", cmd)
		}
	}

	return nil
}
