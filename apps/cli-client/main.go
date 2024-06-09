package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/Ehab-24/eds-cli-client/args"
	"github.com/Ehab-24/eds-cli-client/master"
	"github.com/Ehab-24/eds-cli-client/video"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	args.Parse()

	fileName, ext := video.GetFileNameAndExt(args.Args.FilePath)
	totalDuration, err := video.GetDuration(args.Args.FilePath)
	check(err)

	payload := master.PostVideoMetadtaPayload{
		Title:              fileName,
		Description:        "This is a video",
		Replication_factor: args.Args.ReplicationFactor,
		Duration:           totalDuration,
		Chunk_duration:     args.Args.ChunkDuration,
	}
	metadata, err := master.PostVideoMetadata(payload)
	check(err)

	var wg sync.WaitGroup
	for i, chunk := range metadata.Chunks {
		wg.Add(1)
		go func(index int, chunk master.Chunk) {
			defer wg.Done()

			log.Printf("⟳ [Chunk:%d] Creating split...\n", chunk.ID)
			chunkServer := video.NewChunkServerClient("http", chunk.Server.IP, chunk.Server.Port)
			startDuration, endDuration := video.GetDurationRange(index, totalDuration)
			chunkFile := fmt.Sprintf("tmp/chunks/%s_%d.mp4", fileName, chunk.ID)
			video.Split(args.Args.FilePath, chunkFile, startDuration, endDuration)

			log.Printf("⟳ [Chunk:%d] Uploading...\n", chunk.ID)
			if err := chunkServer.Upload(metadata.ID, chunk.ID, fileName+"."+ext, chunkFile, ext); err != nil {
				log.Printf(" [Chunk:%d] Error while uploading: %s", chunk.ID, err.Error())
			} else {
				log.Printf(" [Chunk:%d] Upload complete.\n", chunk.ID)
			}
		}(i, chunk)
	}
	wg.Wait()
	log.Println(" Uploaded all chunks")
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
