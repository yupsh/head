package command_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/gloo-foo/testable/assertion"
	"github.com/gloo-foo/testable/run"
	command "github.com/yupsh/head"
)

// ==============================================================================
// Test Default Behavior (10 lines)
// ==============================================================================

func TestHead_DefaultTenLines(t *testing.T) {
	lines := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"}

	result := run.Command(command.Head()).
		WithStdinLines(lines...).
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
	})
}

func TestHead_LessThanDefault(t *testing.T) {
	result := run.Command(command.Head()).
		WithStdinLines("1", "2", "3", "4", "5").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"1", "2", "3", "4", "5",
	})
}

func TestHead_ExactlyTenLines(t *testing.T) {
	lines := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}

	result := run.Command(command.Head()).
		WithStdinLines(lines...).
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, lines)
}

func TestHead_EmptyInput(t *testing.T) {
	result := run.Quick(command.Head())

	assertion.NoError(t, result.Err)
	assertion.Empty(t, result.Stdout)
}

func TestHead_SingleLine(t *testing.T) {
	result := run.Command(command.Head()).
		WithStdinLines("only line").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"only line"})
}

// ==============================================================================
// Test Custom Line Counts
// ==============================================================================

func TestHead_ThreeLines(t *testing.T) {
	result := run.Command(command.Head(command.LineCount(3))).
		WithStdinLines("1", "2", "3", "4", "5").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"1", "2", "3"})
}

func TestHead_OneLine(t *testing.T) {
	result := run.Command(command.Head(command.LineCount(1))).
		WithStdinLines("first", "second", "third").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"first"})
}

func TestHead_FiveLines(t *testing.T) {
	result := run.Command(command.Head(command.LineCount(5))).
		WithStdinLines("a", "b", "c", "d", "e", "f", "g").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"a", "b", "c", "d", "e"})
}

func TestHead_LargeCount(t *testing.T) {
	// Request 100 lines, but only provide 5
	result := run.Command(command.Head(command.LineCount(100))).
		WithStdinLines("1", "2", "3", "4", "5").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"1", "2", "3", "4", "5"})
}

// ==============================================================================
// Test With Empty Lines
// ==============================================================================

func TestHead_EmptyLine(t *testing.T) {
	result := run.Command(command.Head(command.LineCount(3))).
		WithStdinLines("a", "", "b", "", "c").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"a", "", "b"})
}

func TestHead_AllEmptyLines(t *testing.T) {
	result := run.Command(command.Head(command.LineCount(3))).
		WithStdinLines("", "", "", "", "").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"", "", ""})
}

func TestHead_EmptyLinesAtStart(t *testing.T) {
	result := run.Command(command.Head(command.LineCount(3))).
		WithStdinLines("", "", "content", "more").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"", "", "content"})
}

// ==============================================================================
// Test With Whitespace
// ==============================================================================

func TestHead_Whitespace(t *testing.T) {
	result := run.Command(command.Head(command.LineCount(3))).
		WithStdinLines("  spaces", "\ttabs", "   both \t", "normal").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"  spaces",
		"\ttabs",
		"   both \t",
	})
}

func TestHead_WhitespaceOnly(t *testing.T) {
	result := run.Command(command.Head(command.LineCount(2))).
		WithStdinLines("   ", "\t\t", "content").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"   ", "\t\t"})
}

func TestHead_LeadingSpaces(t *testing.T) {
	result := run.Command(command.Head(command.LineCount(2))).
		WithStdinLines("    line1", "        line2", "line3").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"    line1",
		"        line2",
	})
}

func TestHead_TrailingSpaces(t *testing.T) {
	result := run.Command(command.Head(command.LineCount(2))).
		WithStdinLines("line1    ", "line2        ", "line3").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"line1    ",
		"line2        ",
	})
}

// ==============================================================================
// Test With Unicode
// ==============================================================================

func TestHead_Unicode_Japanese(t *testing.T) {
	result := run.Command(command.Head(command.LineCount(2))).
		WithStdinLines("こんにちは", "世界", "日本語").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"こんにちは", "世界"})
}

func TestHead_Unicode_Mixed(t *testing.T) {
	result := run.Command(command.Head(command.LineCount(3))).
		WithStdinLines("Hello", "世界", "123", "مرحبا", "test").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"Hello", "世界", "123"})
}

func TestHead_Unicode_Emoji(t *testing.T) {
	result := run.Command(command.Head(command.LineCount(2))).
		WithStdinLines("😀", "👋", "🌍", "🚀").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"😀", "👋"})
}

func TestHead_Unicode_Arabic(t *testing.T) {
	result := run.Command(command.Head(command.LineCount(2))).
		WithStdinLines("مرحبا", "سلام", "أهلا").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 2)
}

// ==============================================================================
// Test With Special Characters
// ==============================================================================

func TestHead_SpecialCharacters(t *testing.T) {
	result := run.Command(command.Head(command.LineCount(3))).
		WithStdinLines("!@#$%", "^&*()", "{}[]", "<>?").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"!@#$%", "^&*()", "{}[]"})
}

func TestHead_Punctuation(t *testing.T) {
	result := run.Command(command.Head(command.LineCount(2))).
		WithStdinLines("Hello!", "How are you?", "Goodbye.").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"Hello!", "How are you?"})
}

func TestHead_Quotes(t *testing.T) {
	result := run.Command(command.Head(command.LineCount(2))).
		WithStdinLines(`"double"`, `'single'`, "`backtick`").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{`"double"`, `'single'`})
}

// ==============================================================================
// Test Edge Cases
// ==============================================================================

func TestHead_VeryLongLine(t *testing.T) {
	longLine := strings.Repeat("a", 10000)
	result := run.Command(command.Head(command.LineCount(1))).
		WithStdinLines(longLine, "short").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 1)
	assertion.Equal(t, result.Stdout[0], longLine, "long line")
}

func TestHead_ManyLines(t *testing.T) {
	lines := make([]string, 1000)
	for i := range lines {
		lines[i] = fmt.Sprintf("line %d", i+1)
	}

	result := run.Command(command.Head(command.LineCount(10))).
		WithStdinLines(lines...).
		Run()

	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 10)
	assertion.Equal(t, result.Stdout[0], "line 1", "first line")
	assertion.Equal(t, result.Stdout[9], "line 10", "tenth line")
}

func TestHead_ExactLineCount(t *testing.T) {
	// Request exactly as many lines as provided
	result := run.Command(command.Head(command.LineCount(5))).
		WithStdinLines("1", "2", "3", "4", "5").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"1", "2", "3", "4", "5"})
}

// ==============================================================================
// Test Error Handling
// ==============================================================================

func TestHead_InputError(t *testing.T) {
	result := run.Command(command.Head()).
		WithStdinError(errors.New("read failed")).
		Run()

	assertion.ErrorContains(t, result.Err, "read failed")
}

func TestHead_OutputError(t *testing.T) {
	result := run.Command(command.Head()).
		WithStdinLines("test").
		WithStdoutError(errors.New("write failed")).
		Run()

	assertion.ErrorContains(t, result.Err, "write failed")
}

// ==============================================================================
// Test Flags
// ==============================================================================

func TestHead_BytesFlag(t *testing.T) {
	// ByteCount flag is defined but not currently used in implementation
	result := run.Command(command.Head(command.ByteCount(10))).
		WithStdinLines("a", "b", "c").
		Run()

	assertion.NoError(t, result.Err)
	// Current implementation ignores bytes flag
}

func TestHead_QuietFlag(t *testing.T) {
	// Quiet flag is defined but not currently used in implementation
	result := run.Command(command.Head(command.Quiet)).
		WithStdinLines("a", "b").
		Run()

	assertion.NoError(t, result.Err)
	// Current implementation ignores quiet flag
}

// ==============================================================================
// Table-Driven Tests
// ==============================================================================

func TestHead_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		count    command.LineCount
		input    []string
		expected []string
	}{
		{
			name:     "three from five",
			count:    3,
			input:    []string{"a", "b", "c", "d", "e"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "one line",
			count:    1,
			input:    []string{"first", "second", "third"},
			expected: []string{"first"},
		},
		{
			name:     "all lines",
			count:    5,
			input:    []string{"a", "b"},
			expected: []string{"a", "b"},
		},
		{
			name:     "with empty lines",
			count:    3,
			input:    []string{"a", "", "b", "c"},
			expected: []string{"a", "", "b"},
		},
		{
			name:     "unicode",
			count:    2,
			input:    []string{"こんにちは", "世界", "日本語"},
			expected: []string{"こんにちは", "世界"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := run.Command(command.Head(tt.count)).
				WithStdinLines(tt.input...).
				Run()

			assertion.NoError(t, result.Err)
			assertion.Lines(t, result.Stdout, tt.expected)
		})
	}
}

// ==============================================================================
// Test Real-World Scenarios
// ==============================================================================

func TestHead_LogFileTop(t *testing.T) {
	result := run.Command(command.Head(command.LineCount(3))).
		WithStdinLines(
			"[2024-01-01] First entry",
			"[2024-01-02] Second entry",
			"[2024-01-03] Third entry",
			"[2024-01-04] Fourth entry",
			"[2024-01-05] Fifth entry",
		).Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"[2024-01-01] First entry",
		"[2024-01-02] Second entry",
		"[2024-01-03] Third entry",
	})
}

func TestHead_CodeSnippet(t *testing.T) {
	result := run.Command(command.Head(command.LineCount(3))).
		WithStdinLines(
			"package main",
			"",
			"import \"fmt\"",
			"",
			"func main() {",
			"    fmt.Println(\"Hello\")",
			"}",
		).Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"package main",
		"",
		"import \"fmt\"",
	})
}

func TestHead_CSVHeader(t *testing.T) {
	result := run.Command(command.Head(command.LineCount(1))).
		WithStdinLines(
			"Name,Age,City",
			"Alice,30,NYC",
			"Bob,25,LA",
			"Carol,35,SF",
		).Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"Name,Age,City"})
}

func TestHead_DataPreview(t *testing.T) {
	// Preview first few lines of data
	result := run.Command(command.Head(command.LineCount(5))).
		WithStdinLines(
			"Row 1",
			"Row 2",
			"Row 3",
			"Row 4",
			"Row 5",
			"Row 6",
			"Row 7",
			"Row 8",
			"Row 9",
			"Row 10",
		).Run()

	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 5)
	assertion.Equal(t, result.Stdout[0], "Row 1", "first row")
	assertion.Equal(t, result.Stdout[4], "Row 5", "fifth row")
}

// ==============================================================================
// Test Line Number Boundaries
// ==============================================================================

func TestHead_OneFromMany(t *testing.T) {
	result := run.Command(command.Head(command.LineCount(1))).
		WithStdinLines("first", "second", "third", "fourth").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"first"})
}

func TestHead_TwoFromThree(t *testing.T) {
	result := run.Command(command.Head(command.LineCount(2))).
		WithStdinLines("a", "b", "c").
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"a", "b"})
}

func TestHead_NineFromTen(t *testing.T) {
	lines := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	expected := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}

	result := run.Command(command.Head(command.LineCount(9))).
		WithStdinLines(lines...).
		Run()

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, expected)
}

// ==============================================================================
// Test Mixed Content
// ==============================================================================

func TestHead_MixedContent(t *testing.T) {
	result := run.Command(command.Head(command.LineCount(5))).
		WithStdinLines(
			"normal text",
			"",
			"line with\ttabs",
			"  spaces  ",
			"unicode: 日本語",
			"special: !@#$",
			"more content",
		).Run()

	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 5)
	assertion.Equal(t, result.Stdout[0], "normal text", "first line")
	assertion.Equal(t, result.Stdout[1], "", "empty line")
	assertion.Equal(t, result.Stdout[4], "unicode: 日本語", "unicode line")
}

// ==============================================================================
// Test Stream Behavior
// ==============================================================================

func TestHead_StopsReading(t *testing.T) {
	// head should stop reading after N lines
	// Create 100 lines but only request 5
	lines := make([]string, 100)
	for i := range lines {
		lines[i] = fmt.Sprintf("line %d", i+1)
	}

	result := run.Command(command.Head(command.LineCount(5))).
		WithStdinLines(lines...).
		Run()

	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 5)
	// Verify it got the correct first 5 lines
	assertion.Equal(t, result.Stdout[0], "line 1", "first line")
	assertion.Equal(t, result.Stdout[4], "line 5", "fifth line")
}

