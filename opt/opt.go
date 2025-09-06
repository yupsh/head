package opt

// Custom types for parameters
type LineCount int
type ByteCount int

// Boolean flag types with constants
type QuietFlag bool

const (
	Quiet   QuietFlag = true
	NoQuiet QuietFlag = false
)

// Flags represents the configuration options for the head command
type Flags struct {
	Lines LineCount // Number of lines to show (default 10)
	Bytes ByteCount // Number of bytes to show
	Quiet QuietFlag // Suppress headers when multiple files
}

// Configure methods for the opt system
func (l LineCount) Configure(flags *Flags) { flags.Lines = l }
func (b ByteCount) Configure(flags *Flags) { flags.Bytes = b }
func (f QuietFlag) Configure(flags *Flags) { flags.Quiet = f }
