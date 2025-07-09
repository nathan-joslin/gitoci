package actions

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/act3-ai/gitoci/internal/comms"
)

type Option struct {
	Verbosity int // TODO: this likely isnt needed, also progress is the option for Git's -v
}

// option handles and responds to the option subcommands.
func (action *GitOCI) option(ctx context.Context, cmd comms.Command) error {
	const (
		ok          = "ok"
		unsupported = "unsupported"
	)

	// sanity
	if cmd.CommandType == "" {
		return fmt.Errorf("received empty option command")
	}
	if len(cmd.Data) != 2 {
		return fmt.Errorf("received invalid option value: %v", cmd.Data)
	}

	name := cmd.Data[0]
	value := cmd.Data[1]

	// https://git-scm.com/docs/gitremote-helpers#Documentation/gitremote-helpers.txt-optionnamevalue
	var result string
	err := action.handleOption(ctx, name, value)
	switch {
	case errors.Is(err, comms.ErrUnsupportedCommand):
		slog.DebugContext(ctx, "received unsupported option command", "command", name)
		result = unsupported
	case err != nil:
		slog.ErrorContext(ctx, "failed to handle option command", "command", name)
		result = err.Error()
	default:
		slog.DebugContext(ctx, "successfully handled option command", "command", name)
		result = ok
	}

	if err := action.batcher.Write(result); err != nil {
		return fmt.Errorf("writing option response %s: %w", name, err)
	}
	if err := action.batcher.Flush(false); err != nil {
		return fmt.Errorf("flushing option writes: %w", err)
	}
	return nil
}

// handleOption fulfills and option command if supported.
func (action *GitOCI) handleOption(ctx context.Context, name string, value string) error {
	if !comms.SupportedOption(name) {
		return errors.Join(comms.ErrUnsupportedCommand, fmt.Errorf("unsupported option %s", name))
	}

	cmd := comms.CommandType(name)
	switch cmd {
	case comms.CmdOptionVerbosity:
		return action.verbosity(value)
	default:
		// sanity, should never happen
		slog.DebugContext(ctx, "handleOption not able to fulfill  supposedly supported option command", "command", cmd)
	}

	return nil
}

// verbosity handles the 'option verbosity' command.
func (action *GitOCI) verbosity(value string) error {
	var err error
	action.Verbosity, err = strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("converting verbosity value to int: %w", err)
	}
	return nil
}
