// Hexdumper is a utility to dump binary files as hex (and ASCII).
package main

import (
	"fmt"
	"github.com/augustoroman/hexdump"
	"log"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	kingpin.CommandLine.Author("Augusto Roman")
	kingpin.CommandLine.Help =
		"hexdumper is a utility to dump binary files as hex.  " +
			"Now with color!"
	files := kingpin.Arg("file", "File(s) to dump").Required().ExistingFiles()
	start := kingpin.Flag("start", "Offset to start dump from.").Short('s').Default("0").Bytes()
	num := kingpin.Flag("num", "Number of bytes to dump.  Defaults to the entire file.").Short('n').Bytes()
	width := kingpin.Flag("width", "Number of bytes in each row.").Short('w').Default("25").Int()
	kingpin.Parse()

	dumper := hexdump.Config{*width}

	for _, filename := range *files {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("Cannot open %q: %v", filename, err)
		}
		stat, err := file.Stat()
		if err != nil {
			log.Fatalf("Cannot stat %q: %v", filename, err)
		}
		size := stat.Size() - int64(*start)
		if *num != 0 && int64(*num) < size {
			size = int64(*num)
		}
		buf := make([]byte, size)
		if _, err := file.ReadAt(buf, int64(*start)); err != nil {
			log.Fatalf("Cannot read %q: %v", filename, err)
		}
		fmt.Println(dumper.Dump(buf))
		file.Close()
	}
}
