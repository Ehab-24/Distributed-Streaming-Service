package master

import chunkserver "gihu.bocm/Ehab-24/chunk-server/chunk_server"

func GetReplicationServers(videoID int64) ([]*chunkserver.ChunkServer, error) {
  servers := []*chunkserver.ChunkServer{
    // chunkserver.NewChunkServer(1, "127.0.0.1", 5000),
    chunkserver.NewChunkServer(2, "127.0.0.1", 5001),
    chunkserver.NewChunkServer(3, "127.0.0.1", 5002),
  }
  return servers, nil
}
