// Hexdumper is a utility to dump binary files as hex (and ASCII).
package main

import (
	"io"
	"log"
	"os"

	"github.com/augustoroman/hexdump"

	"github.com/alecthomas/kingpin"
)

func main() {
	kingpin.CommandLine.Author("Augusto Roman")
	kingpin.CommandLine.Help =
		"hexdumper is a utility to dump binary files as hex.  " +
			"Now with color!"
	files := kingpin.Arg("file", "File(s) to dump").ExistingFiles()
	start := kingpin.Flag("start", "Offset to start dump from.").Short('s').Default("0").Bytes()
	num := kingpin.Flag("num", "Number of bytes to dump.  Defaults to the entire file.").Short('n').Bytes()
	width := kingpin.Flag("width", "Number of bytes in each row.").Short('w').Default("32").Int()
	kingpin.Parse()

	dumper := hexdump.Config{Width: *width}

	for _, filename := range *files {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("Cannot open %q: %v", filename, err)
		}
		stat, err := file.Stat()
		if err != nil {
			log.Fatalf("Cannot stat %q: %v", filename, err)
		}
		if _, err := file.Seek(int64(*start), os.SEEK_SET); err != nil {
			log.Fatalf("Cannot seek to %d in %q: %v", *start, filename, err)
		}

		size := stat.Size() - int64(*start)
		if *num != 0 && int64(*num) < size {
			size = int64(*num)
		}

		var in io.Reader = io.LimitReader(file, size)
		if err := dumper.Stream(in, os.Stdout); err != nil {
			log.Fatalf("Error reading %q: %v", filename, err)
		}

		file.Close()
	}

	if len(*files) == 0 {
		var in io.Reader = os.Stdin
		if *start != 0 {
			io.Copy(io.Discard, io.LimitReader(in, int64(*start)))
		}

		if *start != 0 {
			in = io.LimitReader(os.Stdin, int64(*num))
		}
		err := dumper.Stream(in, os.Stdout)
		if err != nil {
			log.Fatal(err)
		}
	}
}
