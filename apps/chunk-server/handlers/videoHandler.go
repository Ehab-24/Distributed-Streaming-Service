package handlers

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"gihu.bocm/Ehab-24/chunk-server/args"
	chunkserver "gihu.bocm/Ehab-24/chunk-server/chunk_server"
	"gihu.bocm/Ehab-24/chunk-server/db"
	"gihu.bocm/Ehab-24/chunk-server/master"
	"gihu.bocm/Ehab-24/chunk-server/video"
	"github.com/gin-gonic/gin"
)

func ServeMPDHandler(c *gin.Context) {
	videoID := c.Query("id")
	if videoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No video id was provided."})
		return
	}

	videoMetadata, err := db.GetVideoMetaData(videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Video metadata not found"})
		return
	}

	fileContent, err := os.ReadFile(videoMetadata.FilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading file"})
		return
	}

	c.Data(http.StatusOK, "application/dash+xml", fileContent)
}

type UploadVideoRequest struct {
	VideoID   int64
	ChunkID   int64
	Replicate bool
}

func ParseUploadVideoRequest(c *gin.Context) *UploadVideoRequest {
  videoID, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Query param: Invalid video_id"})
    return nil
  }
  chunkID, err := strconv.ParseInt(c.Query("chunk_id"), 10, 64)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Query param: Invalid chunk_id"})
    return nil
  }
  replicate, err := strconv.ParseBool(c.Query("replicate"))
  if err != nil {
    replicate = false
  }
	return &UploadVideoRequest{
    VideoID:   videoID,
    ChunkID:   chunkID,
    Replicate: replicate,
  }
}

func UploadVideohandler(c *gin.Context) {
  payload := ParseUploadVideoRequest(c)
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing file in request"})
		return
	}
	if payload.Replicate {
		handleClientUpload(c, payload.VideoID, payload.ChunkID, file)
	} else {
		handleReplicationRequest(c, args.Args.ID, payload.VideoID, payload.ChunkID, file)
	}

 //  video.CleanArchiveDir(payload.VideoID, payload.ChunkID)
	// video.CleanTmpDir()
	// video.CleanUploadDir()
}

func handleClientUpload(c *gin.Context, videoID int64, chunkID int64, file *multipart.FileHeader) {

	chunkDir := video.GetChunkDir(videoID, chunkID)
	if !video.IsVideo(file) {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Only video files are allowed. Found %s", file.Filename)})
		return
	}

	uploadedFilePath := filepath.Join(video.GetUploadDir(), fmt.Sprintf("%d_%d", int(videoID), int(chunkID)))
	if err := c.SaveUploadedFile(file, uploadedFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save uploaded file"})
		return
	}

	ext := ".mp4"
	for _, q := range video.VideoQualities {
		fragmentLabel := fmt.Sprintf("%d_%d_%s", videoID, chunkID, q.Label)
		outFile := fmt.Sprintf("%s/%s%s", video.GetTmpDir(), fragmentLabel, ext)

		if err := video.EncodeChunk(uploadedFilePath, outFile, q.Resolution, q.VideoBitrate, q.AudioBitrate); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode video: " + err.Error()})
			return
		} else {
      log.Println("Command executed successfully!")
    }

		inFile := outFile
		outFile = strings.Replace(outFile, ext, "_frag"+ext, 1)
		if err := video.SegmentChunk(inFile, outFile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to segment video: " + err.Error()})
			return
		} else {
      log.Println("Command executed successfully!")
    }

		inFile = outFile
		outFile = fmt.Sprintf("%s/%s.mpd", chunkDir, q.Label)
		if err := video.ToMPD(inFile, outFile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate MPEG-DASH files."})
			return
		} else {
      log.Println("Command executed successfully!")
    }

		// Handle VideoMetaData creation here
		fmt.Println("VIDEO METADATA:", fragmentLabel)
	}

	// TODO: Save chunk metadata to master

	c.JSON(http.StatusCreated, gin.H{"id": videoID, "chunk_id": chunkID})

	servers, err := master.GetReplicationServers(videoID, chunkID)
	if err != nil {
		log.Println(err)
		return
	}

	// TODO:
	ownID := int64(1)
	archivePath, err := video.CreateArchive(videoID, chunkID)
	if err != nil {
		log.Println(err)
		return
	}

	wg := sync.WaitGroup{}
	for _, server := range servers {
		if server.ID == ownID {
			log.Println("ERROR: Cannot replicate to self!")
			continue
		}
		wg.Add(1)
		go func(server *chunkserver.ChunkServer, videoID int64, chunkID int64, archivePath string) {
			defer wg.Done()
			replicate(server, videoID, chunkID, archivePath)
		}(server, videoID, chunkID, archivePath)
	}
	wg.Wait()

	if err := os.Remove(archivePath); err != nil {
		log.Println(err)
		return
	}
}

func replicate(server *chunkserver.ChunkServer, videoID int64, chunkID int64, archivePath string) {
	resp, err := server.Replicate(videoID, chunkID, archivePath)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		log.Println("Failed to replicate to", server.Host())
	}
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		return
	}
	log.Println("Replication response:", result)
	return
}

func handleReplicationRequest(c *gin.Context, serverID int64, videoID int64, chunkID int64, file *multipart.FileHeader) {

	archivePath := video.GetReplicationArchivePath(serverID, videoID, chunkID)
	log.Println("Saving archive to: ", archivePath)
	if err := c.SaveUploadedFile(file, archivePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save uploaded archive"})
		return
	}

	if _, err := os.Stat(archivePath); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Archive does not exist: %s", archivePath)})
		return
	}

	archive, err := zip.OpenReader(archivePath)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error while opening archive: %s", err.Error())})
		return
	}
	outDir := video.GetChunkDir(videoID, chunkID)
	video.UnzipArchive(archive, outDir)
	archive.Close()
  // TODO: video.CleanArchiveDir(videoID, chunkID)

	c.JSON(http.StatusCreated, gin.H{"message": "Replication successful!"})
}
