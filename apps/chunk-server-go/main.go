package main

import (
	"fmt"

	"gihu.bocm/Ehab-24/chunk-server/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
	}))

	r.POST("/video/upload", handlers.UploadVideo)
	r.GET("/video/serve", handlers.ServeMPD)
	r.Static("media/processed", "./media/processed")
	r.GET("/health", handlers.HealthCheck)

	port := 5000
	r.Run(fmt.Sprintf(":%d", port))
}
