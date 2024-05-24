package args

import (
	"flag"
	"fmt"
	"os"
)

type TArgs struct {
	ID    int64
	Port  int
}

var Args TArgs

func Parse() {
  var id int
  flag.IntVar(&id, "id", 0, "Unique integer identifier for this chunk server")
	flag.IntVar(&Args.Port, "port", 8080, "Port to listen on")
	flag.Usage = printUsage

	flag.Parse()
  Args.ID = int64(id)
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage: chunk-server [clean]")
	flag.PrintDefaults()
}
