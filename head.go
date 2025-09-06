package head

import (
	"bufio"
	"context"
	"fmt"
	"io"

	yup "github.com/yupsh/framework"
	"github.com/yupsh/framework/opt"

	localopt "github.com/yupsh/head/opt"
)

// Flags represents the configuration options for the head command
type Flags = localopt.Flags

// Command implementation using StandardCommand abstraction
type command struct {
	yup.StandardCommand[Flags]
}

// Head creates a new head command with the given parameters
func Head(parameters ...any) yup.Command {
	args := opt.Args[string, Flags](parameters...)
	cmd := command{
		StandardCommand: yup.StandardCommand[Flags]{
			Positional: args.Positional,
			Flags:      args.Flags,
			Name:       "head",
		},
	}
	// Set default if no lines/bytes specified
	if cmd.Flags.Lines == 0 && cmd.Flags.Bytes == 0 {
		cmd.Flags.Lines = 10
	}
	return cmd
}

func (c command) Execute(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) error {
	return yup.ProcessFilesWithContext(
		ctx, c.Positional, stdin, stdout, stderr,
		yup.FileProcessorOptions{
			CommandName:     "head",
			ShowHeaders:     !bool(c.Flags.Quiet),
			BlankBetween:    true,
			ContinueOnError: true,
		},
		func(ctx context.Context, source yup.InputSource, output io.Writer) error {
			return c.processReader(ctx, source.Reader, output, source.Filename, false)
		},
	)
}

func (c command) processReader(ctx context.Context, reader io.Reader, output io.Writer, filename string, showHeader bool) error {
	if showHeader {
		fmt.Fprintf(output, "==> %s <==\n", filename)
	}

	if c.Flags.Bytes > 0 {
		return c.processBytes(ctx, reader, output)
	}

	return c.processLines(ctx, reader, output)
}

func (c command) processLines(ctx context.Context, reader io.Reader, output io.Writer) error {
	scanner := bufio.NewScanner(reader)
	count := 0

	for yup.ScanWithContext(ctx, scanner) && count < int(c.Flags.Lines) {
		fmt.Fprintln(output, scanner.Text())
		count++
	}

	// Check if context was cancelled
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return err
	}

	return scanner.Err()
}

func (c command) processBytes(ctx context.Context, reader io.Reader, output io.Writer) error {
	// Check for cancellation before starting
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return err
	}

	buf := make([]byte, c.Flags.Bytes)
	n, err := io.ReadFull(reader, buf)

	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return err
	}

	_, writeErr := output.Write(buf[:n])
	return writeErr
}
