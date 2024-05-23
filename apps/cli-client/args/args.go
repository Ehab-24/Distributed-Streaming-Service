package args

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Args struct {
	ChunkDuration int
	FilePath      string
  Port          int
}

func Parse() Args {
  if len(os.Args) < 4 {
    printUsage()
  }
  args := Args{
    FilePath:      os.Args[1],
    ChunkDuration: 10,
    Port:          5000,
  }
  if len(os.Args) == 5 {
    duration, err := strconv.Atoi(os.Args[2])
    args.ChunkDuration = duration
    if err != nil {
      log.Fatal(err)
    }
    port, err := strconv.Atoi(os.Args[4])
    args.Port = port
    if err != nil {
      log.Fatal(err)
    }
  }
  return args
}

func printUsage() {
  fmt.Println("Usage: eds-cli-client <file-path> <duration> <port>")
  os.Exit(1)
}
