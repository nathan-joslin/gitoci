package cli

import (
	"github.com/spf13/cobra"

	"github.com/act3-ai/gitoci/internal/actions"
)

// NewHelloCmd creates a new "hello" command
func NewHelloCmd(tool *actions.Tool) *cobra.Command {
	action := &actions.Hello{Tool: tool}

	cmd := &cobra.Command{
		Use:   "hello",
		Short: "Say hello",
		RunE: func(cmd *cobra.Command, args []string) error {
			return action.Run(cmd.Context(), cmd.OutOrStdout())
		},
	}

	return cmd
}
