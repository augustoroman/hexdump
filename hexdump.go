// Package hexdump provides utility functions to display binary slices as
// hex and printable ASCII.
package hexdump

import (
	"bytes"
	"fmt"
)

// Dump the byte slice to a human-readable hex dump using the default
// configuration.
func Dump(buf []byte) string { return defaultConfig.Dump(buf) }

// Config allows customizing the dump configuration.
type Config struct {
	// Defaults to 32.
	Width int
}

// Dump converts the byte slice to a human-readable hex dump.
func (c Config) Dump(buf []byte) string {
	N := c.Width
	var out bytes.Buffer
	for rowIndex := 0; rowIndex < len(buf)/N; rowIndex++ {
		a, b := rowIndex*N, (rowIndex+1)*N
		row := buf[a:b]
		hex, ascii := printable(row)
		fmt.Fprintf(&out, "%5d: %s | %s\n", a, hex, ascii)
	}
	return out.String()
}

var defaultConfig = Config{32}

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
