package writerutil

import (
	"fmt"
	"io"
	"strings"
)

// LeveledOutput :
type LeveledOutput struct {
	W      io.Writer
	i      int
	prefix string
}

// Indent :
func (lw *LeveledOutput) Indent() {
	lw.i++
	lw.prefix = strings.Repeat("	", lw.i)
}

// UnIndent :
func (lw *LeveledOutput) UnIndent() {
	lw.i--
	lw.prefix = strings.Repeat("	", lw.i)
}

// WithBlock :
func (lw *LeveledOutput) WithBlock(prefix string, callback func()) {
	lw.Println(prefix + " {")
	lw.Indent()
	callback()
	lw.UnIndent()
	lw.Println("}")
}

// Newline :
func (lw *LeveledOutput) Newline() (int, error) {
	return io.WriteString(lw.W, "\n")
}

// Println :
func (lw *LeveledOutput) Println(s string) (int, error) {
	return fmt.Fprintf(lw.W, "%s%s\n", lw.prefix, s)
}

// Printf :
func (lw *LeveledOutput) Printf(format string, args ...interface{}) (int, error) {
	return fmt.Fprintf(lw.W, (lw.prefix + format), args...)
}
