package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/muesli/termenv"
	"github.com/spf13/cobra"

	commands "github.com/act3-ai/go-common/pkg/cmd"
	"github.com/act3-ai/go-common/pkg/logger"
	"github.com/act3-ai/go-common/pkg/runner"
	vv "github.com/act3-ai/go-common/pkg/version"

	"github.com/act3-ai/gitoci/cmd/gitoci/cli"
	"github.com/act3-ai/gitoci/docs"
)

// Retrieves build info
func getVersionInfo() vv.Info {
	info := vv.Get()
	if version != "" {
		info.Version = version
	}
	return info
}

func main() {
	ctx := context.Background()

	info := getVersionInfo()         // Load the version info from the build
	root := cli.NewCLI(info.Version) // Create the root command
	root.SilenceUsage = true         // Silence usage when root is called

	// Layout of embedded documentation to surface in the help command
	// and generate in the gendocs command
	embeddedDocs := docs.Embedded(root)

	// Add common commands
	root.AddCommand(
		commands.NewVersionCmd(info),
		commands.NewGenschemaCmd(docs.Schemas(), docs.SchemaAssociations),
		commands.NewGendocsCmd(embeddedDocs),
		commands.NewInfoCmd(embeddedDocs),
	)

	// Restores the original ANSI processing state on Windows
	var restoreWindowsANSI func() error

	// Store persistent pre run function to avoid overwriting it
	persistentPreRun := root.PersistentPreRun

	// The pre run function logs build info and sets the default output writer
	root.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		slog.SetDefault(logger.FromContext(cmd.Context()))             // Set global slog.Logger
		slog.Info("Software", slog.String("version", info.Version))    // Log version info
		slog.Debug("Software details", slog.Any("info", info))         // Log build info
		termenv.SetDefaultOutput(termenv.NewOutput(cmd.OutOrStdout())) // Set termenv default output

		// Enable ANSI processing on Windows color output
		var err error
		restoreWindowsANSI, err = termenv.EnableVirtualTerminalProcessing(termenv.DefaultOutput())
		if err != nil {
			slog.Error("error enabling ANSI processing", slog.String("error", err.Error()))
		}

		if persistentPreRun != nil {
			persistentPreRun(cmd, args)
		}
	}

	// The post run function restores the terminal
	root.PersistentPostRun = func(cmd *cobra.Command, args []string) {
		// Restore original ANSI processing state on Windows
		if err := restoreWindowsANSI(); err != nil {
			slog.Error("error restoring ANSI processing state", slog.String("error", err.Error()))
		}
	}

	// Run the root command
	if err := runner.Run(ctx, root, "GITOCI_VERBOSITY"); err != nil {
		os.Exit(1)
	}
}
