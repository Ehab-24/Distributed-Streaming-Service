package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

func HealthCheck(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{
		"status":       "ok",
		"cpu":          gin.H{"usage": cpuUsage},
    "memory":      gin.H{"total": memInfo.Total, "used": memInfo.Used, "free": memInfo.Free, "active": memInfo.Active, "inactive": memInfo.Inactive, "cached": memInfo.Cached, "buffered": memInfo.Buffers, "shared": memInfo.Shared, "slab": memInfo.Slab},
	})
}
