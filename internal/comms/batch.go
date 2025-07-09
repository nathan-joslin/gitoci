// Package comms implements utilities for communicating with Git, i.e. reading from and writing to Git.
package comms

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
)

// BatchReadWriter supports both reading from and writing to Git in batches.
type BatchReadWriter interface {
	BatchReader
	BatchWriter
}

// BatchReader extends Reader to support reading sets of commands
// provided by Git.
type BatchReader interface {
	Reader

	// ReadBatch reads lines from Git until an empty line is encountered.
	ReadBatch(context.Context) ([]Command, error)
}

type BatchWriter interface {
	Writer

	// WriteBatch writes a batch of messages to Git, which
	// MAY need to be followed up with a flush.
	WriteBatch(...string) error
}

// type ReadWriter interface {
// 	Reader
// 	Writer
// }

// Reader reads a single command from Git.
type Reader interface {
	// Read reads a single line from Git.
	Read(context.Context) (Command, error)
}

// Writer is used to Write single lines to Git, completed with a Flush.
type Writer interface {
	// Write buffers a single line write to Git. One or more
	// calls MAY need to be followed up with a flush.
	Write(string) error

	// Flush writes buffered Write(s) to Git, optionally followed up with a blank line.
	Flush(bool) error
}

// batcher implements BatchReadWriter.
type batcher struct {
	in  *bufio.Scanner
	out *bufio.Writer
}

// NewBatcher returns a buffered BatchReadWriter.
func NewBatcher(in io.Reader, out io.Writer) BatchReadWriter {
	return &batcher{
		in:  bufio.NewScanner(in),
		out: bufio.NewWriter(out),
	}
}

// Read parses a single command received by Git.
func (b *batcher) Read(ctx context.Context) (Command, error) {
	ok := b.in.Scan()
	switch {
	case !ok && b.in.Err() != nil:
		return Command{}, fmt.Errorf("reading single command from Git: %w", b.in.Err())
	case !ok:
		// EOF
		return Command{CommandType: CmdEmpty}, nil
	default:
		txt := b.in.Text()
		slog.DebugContext(ctx, "read line from Git", "text", txt)
		cmd, err := b.parseCommand(ctx, txt)
		if err != nil {
			return Command{}, fmt.Errorf("parsing Git command: %w", err)
		}
		return cmd, nil
	}
}

// ReadBatch reads lines from Git until an empty line is encountered.
func (b *batcher) ReadBatch(ctx context.Context) ([]Command, error) {
	result := make([]Command, 0, 2)
	for b.in.Scan() {
		line := b.in.Text()
		if line == "" {
			break
		}
		cmd, err := b.parseCommand(ctx, line)
		if err != nil {
			return nil, fmt.Errorf("parsing Git command: %w", err)
		}
		result = append(result, cmd)
	}
	if b.in.Err() != nil {
		return result, fmt.Errorf("scanning input: %w", b.in.Err())
	}
	return result, nil
}

// WriteBatch writes Message(s) to Git, completing the batch with a blank line, and flushing the buffered writes to Git.
func (b *batcher) WriteBatch(lines ...string) error {
	for _, line := range lines {
		if _, err := fmt.Fprintln(b.out, line); err != nil {
			return fmt.Errorf("writing to Git, line: %s: %w", line, err)
		}
	}

	return b.Flush(true)
}

// Write buffers a single line write to Git, must be followed up with a flush.
func (b *batcher) Write(line string) error {
	if _, err := fmt.Fprintln(b.out, line); err != nil {
		return fmt.Errorf("writing to Git, line: %s: %w", line, err)
	}

	return nil
}

// Flush writes buffered Write(s) to Git, followed up with a blank line.
func (b *batcher) Flush(blankLine bool) error {
	if blankLine {
		if _, err := fmt.Fprintln(b.out); err != nil {
			return fmt.Errorf("writing blank line to Git: %w", err)
		}
	}

	if err := b.out.Flush(); err != nil {
		return fmt.Errorf("flushing writes to Git: %w", err)
	}

	return nil
}

// parseCommand reads a single line received from Git, turning it into a Command
// easily identified by CommandType.
func (b *batcher) parseCommand(ctx context.Context, line string) (Command, error) {
	slog.DebugContext(ctx, "parsing command")
	fields := strings.Fields(line)
	if len(fields) < 1 {
		return Command{
			CommandType: CmdEmpty,
		}, nil
	}

	cmd := fields[0]
	switch cmd {
	case string(CmdCapabilities):
		return Command{
			CommandType: CmdCapabilities,
		}, nil
	case string(CmdOption):
		// TODO: we should try to not make options fatal, but we may have to
		// make an exception for force (or others).
		if len(fields) != 3 {
			slog.ErrorContext(ctx, "invalid number of arguments to option command", "got", fmt.Sprintf("%d", len(fields)), "want", "3")
			return Command{}, fmt.Errorf("invalid number of args to option command")
		} else {
			return Command{
				CommandType: CmdOption,
				Data:        fields[1:],
			}, nil
		}
	default:
		return Command{}, fmt.Errorf("%w: %s", ErrUnsupportedCommand, cmd)
	}
}
