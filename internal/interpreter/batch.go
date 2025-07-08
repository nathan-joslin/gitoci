// Package interpreter implements utilities for reading from and writing to Git.
package interpreter

import (
	"bufio"
	"fmt"
	"io"
)

// Batcher provides methods for reading lines from and writing
// lines to Git in batches.
type Batcher interface {
	// Read reads lines from Git until an empty line is encountered.
	ReadBatch() []string

	// WriteBatch writes line(s) to Git, completing the set with a blank line, and flushing the buffered writes to Git.
	WriteBatch(lines ...string) error

	// Write buffers a single line write to Git, must be followed up with a flush.
	// Write(line string) error

	// Flush writes buffered Write(s) to Git, followed up with a blank line.
	// Flush() error
}

// batcher implements BatchInterpreter.
type batcher struct {
	in  *bufio.Scanner
	out *bufio.Writer
}

// NewBatcher returns a buffered Batcher.
func NewBatcher(in io.Reader, out io.Writer) Batcher {
	return &batcher{
		in:  bufio.NewScanner(in),
		out: bufio.NewWriter(out),
	}
}

// Read reads lines from Git until an empty line is encountered.
func (i *batcher) ReadBatch() []string {
	result := make([]string, 0, 2)
	for i.in.Scan() {
		line := i.in.Text()
		if line == "" {
			break
		}
		result = append(result, line)
	}
	return result
}

// WriteBatch writes line(s) to Git, completing the set with a blank line, and flushing the buffered writes to Git.
func (i *batcher) WriteBatch(lines ...string) error {
	for _, line := range lines {
		if _, err := fmt.Fprintln(i.out, line); err != nil {
			return fmt.Errorf("writing to Git, line: %s: %w", line, err)
		}
	}

	return i.flush()
}

// write buffers a single line write to Git, must be followed up with a flush.
// func (i *batcher) write(line string) error {
// 	if _, err := fmt.Fprintln(i.out, line); err != nil {
// 		return fmt.Errorf("writing to Git, line: %s: %w", line, err)
// 	}

// 	return nil
// }

// flush writes buffered Write(s) to Git, followed up with a blank line.
func (i *batcher) flush() error {
	if _, err := fmt.Fprintln(i.out); err != nil {
		return fmt.Errorf("writing blank line to Git: %w", err)
	}

	if err := i.out.Flush(); err != nil {
		return fmt.Errorf("flushing writes to Git: %w", err)
	}

	return nil
}
