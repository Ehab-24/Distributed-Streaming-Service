package args

import (
	"flag"
	"fmt"
	"os"
)

type TArgs struct {
	VideoTitle        string
	VideoDescription  string
	ChunkDuration     int
	FilePath          string
	ReplicationFactor int
}

var Args TArgs

func Parse() {
	if len(os.Args) < 4 {
		printUsage()
	}
	flag.IntVar(&Args.ReplicationFactor, "replicas", 1, "The number of replicas of the video chunks")
	flag.IntVar(&Args.ChunkDuration, "chunk-duration", 10, "Duration of each video chunk in seconds")
	flag.StringVar(&Args.FilePath, "file", "", "Path to the video file to upload")
	flag.StringVar(&Args.VideoTitle, "title", "", "The title of video")
	flag.StringVar(&Args.VideoDescription, "description", "", "The description of video")
	flag.Parse()
}

func printUsage() {
	fmt.Println("Usage: eds-cli-client <file-path> <duration> <port>")
	os.Exit(1)
}
