package main

import (
	"flag"
	"fmt"
	"github.com/cheggaaa/pb/v3"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()
	fmt.Println("inputs: ", from, to, limit, offset)

	// start new bar
	bar := pb.StartNew(100)

	err := Copy(from, to, offset, limit, bar)
	if err != nil {
		fmt.Println("result: ", err)
	}

	// finish bar
	bar.Finish()

}
