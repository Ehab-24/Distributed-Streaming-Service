package video

import (
	"log"
	"mime/multipart"
	"os/exec"
	"path/filepath"
	"strings"
)

func IsVideo(file *multipart.FileHeader) bool {
	// Add more video file extensions that are compatible with MPEG_DASH
	allowedExtensions := []string{".mp4"}
	fileExtension := strings.ToLower(filepath.Ext(file.Filename))
	for _, ext := range allowedExtensions {
		if ext == fileExtension {
			return true
		}
	}
	return false
}

func EncodeChunk(inputFile, outputFile, resolution, videoBitrate, audioBitrate string) error {
	cmd := exec.Command("ffmpeg",
		"-i", inputFile,
		"-c:v", "libx264",
		"-b:v", videoBitrate,
		"-s", resolution,
		"-profile:v", "main",
		"-level", "3.1",
		"-preset", "medium",
		"-x264-params", "scenecut=0:open_gop=0:min-keyint=72:keyint=72",
		"-c:a", "aac",
		"-b:a", audioBitrate,
		"-f", "mp4",
		outputFile,
	)
	log.Println(cmd.String())
	err := cmd.Run()
	return err
}

func SegmentChunk(inputFile, outputFile string) error {
	cmd := exec.Command("mp4fragment",
		inputFile,
		outputFile,
	)
	log.Println(cmd.String())
	err := cmd.Run()
	return err
}

func ToMPD(inputFile, outputFile string) error {
	cmd := exec.Command("ffmpeg",
		"-i", inputFile,
		"-c:v", "copy",
		"-c:a", "copy",
		"-f", "dash",
		outputFile,
	)
	log.Println(cmd.String())
	err := cmd.Run()
	return err
}
