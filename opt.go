package command

type LineCount int
type ByteCount int

type QuietFlag bool

const (
	Quiet   QuietFlag = true
	NoQuiet QuietFlag = false
)

type flags struct {
	Lines LineCount
	Bytes ByteCount
	Quiet QuietFlag
}

func (l LineCount) Configure(flags *flags) { flags.Lines = l }
func (b ByteCount) Configure(flags *flags) { flags.Bytes = b }
func (f QuietFlag) Configure(flags *flags) { flags.Quiet = f }
