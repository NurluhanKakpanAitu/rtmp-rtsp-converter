package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"rtmp-rtsp-converter/internal/config"
	"rtmp-rtsp-converter/internal/converter"
	"rtmp-rtsp-converter/internal/handlers"
	"rtmp-rtsp-converter/internal/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger.Init(cfg.Logging.Level, cfg.Logging.Format)
	log := logger.GetLogger()

	log.Info("Starting RTMP to RTSP converter service")

	// Initialize converter
	conv := converter.NewConverter(cfg.RTSP.Host, cfg.RTSP.Port, cfg.Streams.MaxConcurrent)

	// Initialize handlers
	streamHandler := handlers.NewStreamHandler(conv)

	// Setup Gin router
	if cfg.Logging.Level != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Routes
	api := router.Group("/api/v1")
	{
		api.GET("/health", streamHandler.Health)
		api.POST("/streams", streamHandler.CreateStream)
		api.GET("/streams", streamHandler.ListStreams)
		api.GET("/streams/:id", streamHandler.GetStream)
		api.DELETE("/streams/:id", streamHandler.StopStream)
	}

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Infof("Server starting on %s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
