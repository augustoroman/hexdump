// Package hexdump provides utility functions to display binary slices as
// hex and printable ASCII.
package hexdump

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

// Dump the byte slice to a human-readable hex dump using the default
// configuration.
func Dump(buf []byte) string { return defaultConfig.Dump(buf) }

// Config allows customizing the dump configuration.
type Config struct {
	// Number of bytes from the input buffer to print in a single row. The default
	// is 32.
	Width int
}

type dumpState struct {
	Config
	rowIndex    int
	maxRowWidth int
}

func (s *dumpState) dump(out io.Writer, buf []byte) {
	N := s.Width
	for i := 0; i*N < len(buf); i++ {
		a, b := i*N, (i+1)*N
		if b > len(buf) {
			b = len(buf)
		}
		row := buf[a:b]
		hex, ascii := printable(row)

		if len(row) < s.maxRowWidth {
			padding := s.maxRowWidth*2 + s.maxRowWidth/4 - len(row)*2 - len(row)/4
			hex += strings.Repeat(" ", padding)
		}
		s.maxRowWidth = len(row)

		fmt.Fprintf(out, "%5d: %s | %s\n", s.rowIndex*N, hex, ascii)
		s.rowIndex++
	}
}

func (c Config) newDumpState() *dumpState {
	s := &dumpState{Config: c}
	if s.Width == 0 {
		s.Width = kDefaultWidth
	}
	return s
}

// Dump converts the byte slice to a human-readable hex dump.
func (c Config) Dump(data []byte) string {
	var out bytes.Buffer
	c.newDumpState().dump(&out, data)
	return out.String()

}

// Read will read from the input io.Reader and write human-readable, formatted
// hexdumps (with color annotations) to the output. The entire input reader is
// consumed. Any errors other than io.EOF are returned.
func (c Config) Stream(in io.Reader, out io.Writer) error {
	s := c.newDumpState()
	buf := make([]byte, 1*c.Width)
	for {
		n, err := io.ReadFull(in, buf)
		s.dump(out, buf[:n])
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil
		} else if err != nil {
			return err
		}
	}
}

const kDefaultWidth = 32

var defaultConfig = Config{kDefaultWidth}

const (
	kESC        = "\033["
	kRESET      = "0"
	kTERM       = "m"
	kWHITE      = "37"
	kGRAY       = "90"
	kDARK_GREEN = "32"
	kGREEN      = "92"
	kDIM        = "2;"
	kNO_ATTR    = ""
)

func color(str, attr, color string) string {
	return kESC + attr + color + kTERM +
		str +
		kESC + kRESET + kTERM
}

func printable(data []byte) (hex, ascii string) {
	s := string(data)
	var hi, lo string
	for i := 0; i < len(s); i++ {
		if (i/8)%2 == 0 {
			hi = kWHITE
			lo = kGRAY
		} else {
			hi = kGREEN
			lo = kDARK_GREEN
		}

		if s[i] < 32 || s[i] >= 127 {
			ascii += color("░", kDIM, lo)
			hex += color(fmt.Sprintf("%02x", s[i]), kDIM, hi)
		} else {
			ascii += color(string(s[i]), kNO_ATTR, hi)
			hex += color(fmt.Sprintf("%02x", s[i]), kNO_ATTR, hi)
		}
		if i%4 == 3 {
			hex += " "
		}
	}
	return hex, ascii
}
