package cli

import (
	"context"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/act3-ai/gitoci/internal/actions"
	"github.com/act3-ai/gitoci/pkg/apis/gitoci.act3-ai.io/v1alpha1"

	"github.com/act3-ai/go-common/pkg/config"
)

// NewCLI creates the base gitoci command
func NewCLI(version string) *cobra.Command {
	// Create Tool action with scheme initialized
	action := actions.NewTool(version)

	// Add environment variable configuration overrides
	action.AddConfigOverrideFunction(func(ctx context.Context, c *v1alpha1.Configuration) error {
		c.ExampleOption = config.EnvBoolOr("GITOCI_EXAMPLE_OPTION", c.ExampleOption)
		c.Name = config.EnvOr("GITOCI_NAME", c.Name)
		return nil
	})

	// cmd represents the base command when called without any subcommands
	cmd := &cobra.Command{
		Use:   "gitoci",
		Short: "A Git remote helper for syncing Git repositories in OCI Registries.",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		SilenceUsage: true,
	}

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	cmd.PersistentFlags().StringArrayVar(&action.ConfigFiles, "config",
		config.EnvArrayOr("GITOCI_CONFIG", action.ConfigFiles, ":"),
		"Config file search locations (env: GITOCI_CONFIG)")

	// Configuration override flag
	var nameFlag string
	cmd.PersistentFlags().StringVar(&nameFlag, "name", "", "Your name (overrides config)")
	action.AddConfigOverrideFunction(func(ctx context.Context, c *v1alpha1.Configuration) error {
		if nameFlag != "" {
			slog.Info("overriding name with flag value", "nameFlag", nameFlag)
			c.Name = nameFlag
		}
		return nil
	})

	cmd.AddCommand(
		NewHelloCmd(action),
	)

	return cmd
}
