package master

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	chunkserver "gihu.bocm/Ehab-24/chunk-server/chunk_server"
)

func GetReplicationServers(videoID int64, chunkID int64) ([]*chunkserver.ChunkServer, error) {

  url := fmt.Sprintf("http://127.0.0.1:8000/replica/servers/?video_id=%d&chunk_id=%d", videoID, chunkID)
  resp, err := http.Get(url)
  if err !=nil {
    return nil, err
  }
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
  if err != nil {
    return nil, err
  }
  if resp.StatusCode != http.StatusOK {
    return nil, fmt.Errorf(string(body))
  }

  var servers chunkserver.ChunkServers
  if err := json.Unmarshal(body, &servers); err != nil {
    return nil, err
  }
  return servers.Servers, nil
}

func NotifyReplicationSuccess(videoID int64, chunkID int64) error {
  return nil
}
