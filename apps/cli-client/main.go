package main

import (
	"fmt"
	"log"

	"github.com/Ehab-24/eds-cli-client/args"
	"github.com/Ehab-24/eds-cli-client/video"
)

func main() {

	args := args.Parse()

	fileName, ext := video.GetFileNameAndExt(args.FilePath)
	_, err := video.GetDuration(args.FilePath)
	check(err)

	startDuration := video.Duration{
		Hours:   0,
		Minutes: 0,
		Seconds: 0,
	}
	endDuration := video.Duration{
		Hours:   0,
		Minutes: 0,
		Seconds: args.ChunkDuration,
	}

	chunkFile := fmt.Sprintf("tmp/chunks/%s_0.mp4", fileName)
	video.Split(args.FilePath, chunkFile, startDuration, endDuration)

  // TODO
  videoID := int64(1)
  chunkID := int64(1)
  videoClient := video.NewClient("http", "127.0.0.1", args.Port)
	videoClient.Upload(videoID, chunkID, fileName+"."+ext, chunkFile, "Test Video 1")
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
