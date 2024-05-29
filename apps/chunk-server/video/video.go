package video

import (
	"fmt"
	"os"
	"time"
)

const (
	DATA_DIR    = "data"
	MEDIA_DIR   = "media"
	TMP_DIR     = MEDIA_DIR + "/tmp"
	PROCESS_DIR = MEDIA_DIR + "/processed"
	UPLOAD_DIR  = MEDIA_DIR + "/uploads"
	ARCHIVE_DIR = MEDIA_DIR + "/archives"
)

type VideoQuality struct {
	Label        string
	Resolution   string
	VideoBitrate string
	AudioBitrate string
}

var VideoQualities = []VideoQuality{
	{"720p", "1280x720", "5000k", "192k"},
	// {"480p", "854x480", "1500k", "128k"},
	// {"360p", "640x360", "800k", "96k"},
	// {"240p", "426x240", "400k", "64k"},
}

type VideoMetaData struct {
	Resolution   string    `json:"resolution"`
	VideoBitrate string    `json:"video_bitrate"`
	AudioBitrate string    `json:"audio_bitrate"`
	FilePath     string    `json:"file_path"`
	UploadedAt   time.Time `json:"uploaded_at"`
	Title        string    `json:"title"`
}

func (v *VideoMetaData) String() string {
	return fmt.Sprintf("%s - %s - %s - %s - %s - %s", v.Title, v.Resolution, v.VideoBitrate, v.AudioBitrate, v.UploadedAt, v.FilePath)
}

func CleanTmpDir() error {
	return os.RemoveAll(GetTmpDir())
}

func CleanUploadDir() error {
	return os.RemoveAll(GetUploadDir())
}

func CleanArchiveDir(videoID int64, chunkID int64) error {
	return os.RemoveAll(GetArchiveDir(videoID, chunkID))
}
