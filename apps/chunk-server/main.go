package main

import (
	"fmt"
	"log"

	"gihu.bocm/Ehab-24/chunk-server/args"
	"gihu.bocm/Ehab-24/chunk-server/handlers"
	"gihu.bocm/Ehab-24/chunk-server/video"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	args.Parse()

	r := gin.Default()
	applyMiddleware(r)
	setupRoutes(r)

	r.Run(fmt.Sprintf(":%d", args.Args.Port))
}

func applyMiddleware(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
	}))
}

func setupRoutes(r *gin.Engine) {
	r.POST("/video/upload", handlers.UploadVideohandler)
	r.GET("/video/serve", handlers.ServeMPDHandler)
  r.GET("/health", handlers.HealthCheckHandler)

  processedDir := video.GetProcessDir()
	r.Static(processedDir, processedDir)
}
