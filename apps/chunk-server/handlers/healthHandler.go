package handlers

import (
	"net/http"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type DiskUsage struct {
	Total uint64 `json:"total"`
	Free  uint64 `json:"free"`
	Used  uint64 `json:"used"`
}

func HealthCheckHandler(c *gin.Context) {
	cpuUsage, err := cpu.Percent(0, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve CPU usage"})
		return
	}
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve memory usage"})
		return
	}
	diskUsage := getDiskUsgae()

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"cpu":    gin.H{"usage": cpuUsage},
		"memory": gin.H{"total": memInfo.Total, "used": memInfo.Used, "free": memInfo.Free, "active": memInfo.Active, "inactive": memInfo.Inactive, "cached": memInfo.Cached, "buffered": memInfo.Buffers, "shared": memInfo.Shared, "slab": memInfo.Slab},
		"disk":   diskUsage,
	})
}

func getDiskUsgae() DiskUsage {
	var stat syscall.Statfs_t
	path := "/"

	syscall.Statfs(path, &stat)
	free := stat.Bavail * uint64(stat.Bsize)
	total := stat.Blocks * uint64(stat.Bsize)
	used := total - free

	return DiskUsage{
		Total: total,
		Free:  free,
		Used:  used,
	}
}
