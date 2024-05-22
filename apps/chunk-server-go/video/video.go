package video

import (
	"fmt"
	"time"
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
