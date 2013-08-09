package hexdump

import (
	"encoding/hex"
	"fmt"
	"strings"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Dump(bytes []byte) string {
	var out string
	const Seg = 4
	const Groups = 2
	const RowSize = Seg * Groups
	hexstr := hex.EncodeToString(bytes)
	N := len(hexstr) / (RowSize)
	for i := 0; i <= N; i++ {
		row := hexstr[i*RowSize : min(len(hexstr), (i+1)*RowSize)]
		out += fmt.Sprintf("  %2d: ", i*RowSize/2)
		for j := 0; j < Groups; j++ {
			if j*Seg < len(row) {
				seg := row[j*Seg : min(len(row), (j+1)*Seg)]
				out += " " + seg
				out += fmt.Sprintf("%*s", Seg-len(seg), "")
			} else {
				out += fmt.Sprintf(" %*s", Seg, "")
			}
		}
		rowstr, _ := hex.DecodeString(row)
		out += fmt.Sprintf("  %s\n", strings.TrimSpace(string(rowstr)))
	}
	return out
}
