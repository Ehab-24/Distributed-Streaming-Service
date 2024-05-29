package video

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"gihu.bocm/Ehab-24/chunk-server/args"
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
  logPrefix := "EncodeChunk: "
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
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
  log.Println(cmd.String())
  err := cmd.Run()
	if outBuf.Len() > 0 {
    outBuf.WriteString(logPrefix)
		log.Println("Standard Output:", outBuf.String())
	}
	if errBuf.Len() > 0 {
    errBuf.WriteString(logPrefix)
		log.Println("Standard Error:", errBuf.String())
	}
  return err
}

func SegmentChunk(inputFile, outputFile string) error {
  time.Sleep(time.Second)
  logPrefix := "SegmentChunk: "
	cmd := exec.Command("mp4fragment",
		inputFile,
		outputFile,
	)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
  log.Println(cmd.String())
  err := cmd.Run()
	if outBuf.Len() > 0 {
    outBuf.WriteString(logPrefix)
		log.Println("Standard Output:", outBuf.String())
	}
	if errBuf.Len() > 0 {
    errBuf.WriteString(logPrefix)
		log.Println("Standard Error:", errBuf.String())
	}
  return err
}

func ToMPD(inputFile, outputFile string) error {
  logPrefix := "ToMPD: "
	cmd := exec.Command("ffmpeg",
		"-i", inputFile,
		"-c:v", "copy",
		"-c:a", "copy",
		"-f", "dash",
		outputFile,
	)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
  log.Println(cmd.String())
  err := cmd.Run()
	if outBuf.Len() > 0 {
    outBuf.WriteString(logPrefix)
		log.Println("Standard Output:", outBuf.String())
	}
	if errBuf.Len() > 0 {
    errBuf.WriteString(logPrefix)
		log.Println("Standard Error:", errBuf.String())
	}
  return err
}

func CreateArchive(videoID int64, chunkID int64) (string, error) {
	fmt.Println("creating zip archive...")
	archivePath := GetArchivePath(videoID, chunkID)
	archive, err := os.Create(archivePath)
	if err != nil {
		return "", err
	}
	defer archive.Close()
	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	dataDir := GetChunkDir(videoID, chunkID)
	files, err := os.ReadDir(dataDir)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		f, err := os.Open(dataDir + "/" + file.Name())
		if err != nil {
			return "", err
		}
		defer f.Close()
		fmt.Println("writing file to archive...")
		w, err := zipWriter.Create(file.Name())
		if err != nil {
			return "", err
		}
		if _, err := io.Copy(w, f); err != nil {
			return "", err
		}
	}
	return archivePath, nil
}

func UnzipArchive(archive *zip.ReadCloser, outDir string) error {
	for _, f := range archive.File {
		log.Println("Unzipping file: ", f.Name)
		filePath := filepath.Join(outDir, f.Name)

		if !strings.HasPrefix(filePath, filepath.Clean(outDir)+string(os.PathSeparator)) {
			return errors.New("Invalid file path")
		}
		if f.FileInfo().IsDir() {
			fmt.Println("creating directory...")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return err
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return err
		}

		dstFile.Close()
		fileInArchive.Close()
	}
	return nil
}

func maybeCreateDir(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
		return true
	}
	return false
}

func GetArchivePath(videoID int64, chunkID int64) string {
	fileName := fmt.Sprintf("%d_%d.zip", videoID, chunkID)
	return filepath.Join(GetArchiveDir(videoID, chunkID), fileName)
}

func GetReplicationArchivePath(serverID int64, videoID int64, chunkID int64) string {
	dirPath := fmt.Sprintf("%s/%d/%s/", DATA_DIR, serverID, ARCHIVE_DIR)
	fileName := fmt.Sprintf("%d_%d.zip", videoID, chunkID)
	maybeCreateDir(dirPath)
	return filepath.Join(dirPath, fileName)
}

func GetUploadDir() string {
	dirPath := fmt.Sprintf("%s/%d/%s/", DATA_DIR, args.Args.ID, UPLOAD_DIR)
	maybeCreateDir(dirPath)
	return dirPath
}

func GetTmpDir() string {
	dirPath := fmt.Sprintf("%s/%d/%s/", DATA_DIR, args.Args.ID, TMP_DIR)
	maybeCreateDir(dirPath)
	return dirPath
}

func GetProcessDir() string {
	dirPath := fmt.Sprintf("%s/%d/%s/", DATA_DIR, args.Args.ID, PROCESS_DIR)
	maybeCreateDir(dirPath)
	return dirPath
}

func GetArchiveDir(videoID int64, chunkID int64) string {
	dirPath := fmt.Sprintf("%s/%d/%s/", DATA_DIR, args.Args.ID, ARCHIVE_DIR)
	maybeCreateDir(dirPath)
	return dirPath
}

func GetChunkDir(videoID int64, chunkID int64) string {
	dirPath := fmt.Sprintf("%s/%d_%d", GetProcessDir(), videoID, chunkID)
	maybeCreateDir(dirPath)
	return dirPath
}
