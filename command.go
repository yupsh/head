package command

import (
	"bufio"
	"context"
	"fmt"
	"io"

	gloo "github.com/gloo-foo/framework"
)

type command gloo.Inputs[gloo.File, flags]

func Head(parameters ...any) gloo.Command {
	return command(gloo.Initialize[gloo.File, flags](parameters...))
}

func (p command) Executor() gloo.CommandExecutor {
	return gloo.Inputs[gloo.File, flags](p).Wrap(func(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) error {
		lineCount := int64(10)
		if p.Flags.Lines > 0 {
			lineCount = int64(p.Flags.Lines)
		}

		scanner := bufio.NewScanner(stdin)
		lineNum := int64(0)

		for scanner.Scan() && lineNum < lineCount {
			lineNum++
			line := scanner.Text()
			if _, err := fmt.Fprintln(stdout, line); err != nil {
				return err
			}
		}

		// Return scanner error if any
		return scanner.Err()
	})
}
