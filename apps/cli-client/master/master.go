package master

import (
	"bytes"
	"encoding/json"
	"errors"
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

func PostVideoMetadta(duration float64) (*VideoMetadata, error) {
	data := map[string]any{
		"title":              "Test Video 1",
		"description":        "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat.Your Description",
		"replication_factor": args.Args.ReplicationFactor,
		"duration":           duration,
	}
	jsonData, _ := json.Marshal(data)

	// TODO: master server ip and port
	resp, err := http.Post("http://127.0.0.1:8000/create/", "application/json", bytes.NewBuffer(jsonData))
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
