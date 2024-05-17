package args

import (
	"log"
	"os"
	"strconv"
)

type Args struct {
    ChunkDuration int
    FilePath string
}


func Parse() Args {
    if len(os.Args) < 2 {
        panic("Please provide a file path")
    }
    args := Args {
        FilePath: os.Args[1],
        ChunkDuration: 10,
    }
    if len(os.Args) == 3 {
        duration, err := strconv.Atoi(os.Args[2])
        args.ChunkDuration = duration
        if err != nil {
            log.Fatal(err)
        }
    }
    return args
}
