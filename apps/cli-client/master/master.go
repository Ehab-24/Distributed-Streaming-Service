package master

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Ehab-24/eds-cli-client/args"
)

type Server struct {
	ID   int    `json:"id"`
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

type Chunk struct {
	ID     int64  `json:"chunk_id"`
	Server Server `json:"server"`
}

type VideoMetadata struct {
	ID     int64   `json:"video_id"`
	Chunks []Chunk `json:"chunks"`
}

type PostVideoMetadtaPayload struct {
	Title              string
	Description        string
	Replication_factor int
	Duration           float64
	Chunk_duration     int
}

func PostVideoMetadata(payload PostVideoMetadtaPayload) (*VideoMetadata, error) {
	data := map[string]any{
		"title":              payload.Title,
		"description":        payload.Description,
		"replication_factor": payload.Replication_factor,
		"duration":           payload.Duration,
		"chunk_duration":     payload.Chunk_duration,
	}
	jsonData, _ := json.Marshal(data)

  url := fmt.Sprintf("%s/create/", args.Args.MasterURL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, errors.New(string(body))
	}

	var result VideoMetadata
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
