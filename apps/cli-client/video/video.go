package video

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/Ehab-24/eds-cli-client/args"
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

func NewChunkServerClient(scheme string, host string, port int) VideoClient {
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
	url := vc.UploadURL() + fmt.Sprintf("?video_id=%d&chunk_id=%d&replicate=%t", videoID, chunkID, true)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func (vc *VideoClient) Upload(videoID int64, chunkID int64, fileName string, filePath string, ext string) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("file", fmt.Sprintf("%d_%d.%s", int(videoID), int(chunkID), ext))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	if err = writeField(writer, "title", args.Args.VideoTitle); err != nil {
		return err
	}
	if err = writeField(writer, "descriptionn", args.Args.VideoDescription); err != nil {
		return err
	}

	err = writer.Close()
  if err != nil {
		return err
	}

	req, err := vc.newUploadRequest(writer, videoID, chunkID, &requestBody)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
  log.Println(string(body))
  return nil
}

func GetFileNameAndExt(filePath string) (string, string) {
	splits := strings.Split(filePath, "/")
	fileNameWithExt := splits[len(splits)-1]
	splits = strings.Split(fileNameWithExt, ".")
	fileName := splits[0]
	ext := splits[len(splits)-1]

	return fileName, ext
}

func writeField(writer *multipart.Writer, name string, value string) error {
	err := writer.WriteField(name, value)
	if err != nil {
		return err
	}
	return nil
}

func GetDurationRange(index int, totalDuration float64) (Duration, Duration) {
  duration := int(math.Ceil(totalDuration))
  startSec := index * args.Args.ChunkDuration
  endSec := min(startSec+args.Args.ChunkDuration, duration)
  startDur := Duration {
    Hours: 0,
    Minutes: startSec/60,
    Seconds: startSec%60,
  }
  endDur := Duration {
    Hours: 0,
    Minutes: endSec/60,
    Seconds: endSec%60,
  }
  return startDur, endDur
}
