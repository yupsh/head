package command

import (
	"bufio"
	"context"
	"fmt"
	"io"

	yup "github.com/gloo-foo/framework"
)

type command yup.Inputs[yup.File, flags]

func Head(parameters ...any) yup.Command {
	return command(yup.Initialize[yup.File, flags](parameters...))
}

func (p command) Executor() yup.CommandExecutor {
	return yup.Inputs[yup.File, flags](p).Wrap(func(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) error {
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

		// Return scanner error if any (but not EOF)
		if err := scanner.Err(); err != nil && err != io.EOF {
			return err
		}

		return nil
	})
}
