package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gihu.bocm/Ehab-24/chunk-server/db"
	"gihu.bocm/Ehab-24/chunk-server/video"
	"github.com/gin-gonic/gin"
)

const (
	TMP_DIR = "./media/tmp/"
	PROCESS_DIR = "./media/processed"
	UPLOAD_DIR = "./media/uploads"
)

func ServeMPD(c *gin.Context) {
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

func UploadVideo(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file was uploaded."})
		return
	}

	if !video.IsVideo(file) {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Only video files are allowed. Found %s", file.Filename)})
		return
	}

  // TODO:
	videoID := 1
	_ = "Test Title 1"

	uploadedFilePath := filepath.Join(UPLOAD_DIR, file.Filename)
	if err := c.SaveUploadedFile(file, uploadedFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save uploaded file"})
		return
	}

	ext := ".mp4"
	for _, q := range video.VideoQualities {
		fragmentLabel := fmt.Sprintf("%d_%s", videoID, q.Label)
		outFile := fmt.Sprintf("%s/%s%s", TMP_DIR, fragmentLabel, ext)

		if err := video.EncodeChunk(uploadedFilePath, outFile, q.Resolution, q.VideoBitrate, q.AudioBitrate); err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode video: " + err.Error()})
			return
		}

		inFile := outFile
		outFile = strings.Replace(outFile, ext, "_frag"+ext, 1)
		if err := video.SegmentChunk(inFile, outFile); err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to segment video: " + err.Error()})
			return
		}

		inFile = outFile
		outFile = fmt.Sprintf("%s/%s", PROCESS_DIR, fragmentLabel)
		if err := video.ToMPD(inFile, outFile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate MPEG-DASH files."})
			return
		}

		// Handle VideoMetaData creation here
		fmt.Println("VIDEO METADATA:", fragmentLabel)
	}

	// Save chunk metadata to master

	c.JSON(http.StatusCreated, gin.H{"id": videoID})
}
