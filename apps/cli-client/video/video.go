package video

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func GetDuration(filePath string) (float64, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", filePath)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return 0, err
	}
	durationStr := strings.TrimSpace(out.String())
	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return 0, err
	}
	return duration, nil
}

func (d *Duration) String() string {
	return fmt.Sprintf("%d:%d:%d", d.Hours, d.Minutes, d.Seconds)
}

func Split(inputFile string, outputFile string, startDuration Duration, endDuration Duration) {
	cmd := exec.Command("ffmpeg", "-ss", startDuration.String(), "-to", endDuration.String(), "-i", inputFile, "-codec", "copy", outputFile)
	cmd.Run()

	log.Printf("%s %s %s %s %s %s %s %s %s %s\n", "ffmpeg", "-ss", startDuration.String(), "-to", endDuration.String(), "-i", inputFile, "-codec", "copy", outputFile)

	verbose := true
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
}

func Quality(inputFile string) (VideoQaulity, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height,bit_rate", "-of", "csv=s=x:p=0", inputFile)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return VideoQaulity{}, err
	}

	splits := strings.Split(out.String(), "x")
	width, err := strconv.Atoi(splits[0])
	if err != nil {
		return VideoQaulity{}, err
	}
	height, err := strconv.Atoi(splits[1])
	if err != nil {
		return VideoQaulity{}, err
	}
	bitrate, err := strconv.Atoi(splits[2])
	if err != nil {
		return VideoQaulity{}, err
	}

	return VideoQaulity{
		bitrate: Bitrate(bitrate),
		resolution: Resolution{
			width:  width,
			height: height,
		},
	}, nil
}

type VideoClient struct {
	Scheme string
	Host   string
	Port   int
}

func NewClient(scheme string, host string, port int) VideoClient {
	return VideoClient{
		Scheme: scheme,
		Host:   host,
		Port:   port,
	}
}

func (vc *VideoClient) URL() string {
	return fmt.Sprintf("http://%s:%d", vc.Host, vc.Port)
}

func (vc *VideoClient) UploadURL() string {
	return fmt.Sprintf("%s/video/upload", vc.URL())
}

func (vc *VideoClient) newUploadRequest(writer *multipart.Writer, videoID int64, chunkID int64, body *bytes.Buffer) (*http.Request, error) {
	url := vc.UploadURL() + fmt.Sprintf("?id=%d&chunk_id=%d&replicate=%t", videoID, chunkID, true)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
    return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
  return req, nil
}

func (vc *VideoClient) Upload(videoID int64, chunkID int64, fileName string, filePath string, title string) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a buffer to write our multipart form data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add the file to the form
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return
	}
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Error copying file:", err)
		return
	}

	// Add the title field to the form
	err = writer.WriteField("title", title)
	if err != nil {
		fmt.Println("Error writing title field:", err)
		return
	}

	// Close the writer to finalize the form data
	err = writer.Close()
	if err != nil {
		fmt.Println("Error closing writer:", err)
		return
	}

  req, err := vc.newUploadRequest(writer, videoID, chunkID, &requestBody)
  if err != nil {
    fmt.Println("Error creating upload request:", err)
    return
  }

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Print the response status and body
	fmt.Println("Response Status:", resp.Status)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	fmt.Println("Response Body:", string(body))
}

func GetFileNameAndExt(filePath string) (string, string) {
	splits := strings.Split(filePath, "/")
	fileNameWithExt := splits[len(splits)-1]
	splits = strings.Split(fileNameWithExt, ".")
	fileName := splits[0]
	ext := splits[len(splits)-1]

	return fileName, ext
}
