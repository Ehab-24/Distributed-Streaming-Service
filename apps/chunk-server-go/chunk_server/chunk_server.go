package chunkserver

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type ChunkServer struct {
	ID   int64
	IP   string
	Port int
}

func (cs *ChunkServer) Host() string {
	return fmt.Sprintf("%s:%d", cs.IP, cs.Port)
}

func (cs *ChunkServer) ReplicateURL() string {
	return fmt.Sprintf("http://%s/video/upload", cs.Host())
}

func NewChunkServer(id int64, ip string, port int) *ChunkServer {
	return &ChunkServer{ID: id, IP: ip, Port: port}
}

func (cs *ChunkServer) Replicate(videoID int64, chunkID int64, filePath string) (*http.Response, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	part, err := writer.CreateFormFile("file", fmt.Sprintf("%d.zip", videoID))
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	if err = writer.Close(); err != nil {
		return nil, err
	}

	req, err := cs.newReplicateRequest(writer, videoID, chunkID, &requestBody)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (cs *ChunkServer) newReplicateRequest(writer *multipart.Writer, videoID int64, chunkID int64, body *bytes.Buffer) (*http.Request, error) {
	url := cs.ReplicateURL() + fmt.Sprintf("?id=%d&chunk_id=%d", videoID, chunkID)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}
