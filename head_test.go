package head_test

import (
	"context"
	"os"
	"strings"

	"github.com/yupsh/head"
	"github.com/yupsh/head/opt"
)

func ExampleHead() {
	ctx := context.Background()
	input := strings.NewReader("Line 1\nLine 2\nLine 3\nLine 4\nLine 5\nLine 6\nLine 7\nLine 8\nLine 9\nLine 10\nLine 11\nLine 12\n")

	cmd := head.Head() // Default 10 lines
	cmd.Execute(ctx, input, os.Stdout, os.Stderr)
	// Output:
	// Line 1
	// Line 2
	// Line 3
	// Line 4
	// Line 5
	// Line 6
	// Line 7
	// Line 8
	// Line 9
	// Line 10
}

func ExampleHead_lineCount() {
	ctx := context.Background()
	input := strings.NewReader("Line 1\nLine 2\nLine 3\nLine 4\nLine 5\n")

	cmd := head.Head(opt.LineCount(3))
	cmd.Execute(ctx, input, os.Stdout, os.Stderr)
	// Output:
	// Line 1
	// Line 2
	// Line 3
}

func ExampleHead_byteCount() {
	ctx := context.Background()
	input := strings.NewReader("Hello World! This is a test string.")

	cmd := head.Head(opt.ByteCount(12))
	cmd.Execute(ctx, input, os.Stdout, os.Stderr)
	// Output:
	// Hello World!
}
