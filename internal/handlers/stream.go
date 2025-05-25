package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"rtmp-rtsp-converter/internal/converter"
	"rtmp-rtsp-converter/internal/logger"
	"rtmp-rtsp-converter/pkg/models"
)

type StreamHandler struct {
	converter *converter.Converter
}

func NewStreamHandler(converter *converter.Converter) *StreamHandler {
	return &StreamHandler{
		converter: converter,
	}
}

func (h *StreamHandler) CreateStream(c *gin.Context) {
	var req models.CreateStreamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	streamID := req.StreamID
	if streamID == "" {
		streamID = uuid.New().String()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stream, err := h.converter.StartStream(ctx, streamID, req.RTMPUrl)
	if err != nil {
		logger.GetLogger().WithError(err).Errorf("Failed to start stream %s", streamID)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.StreamResponse{
		Stream:  stream,
		Message: "Stream created successfully",
	})
}

func (h *StreamHandler) GetStream(c *gin.Context) {
	streamID := c.Param("id")

	stream, err := h.converter.GetStream(streamID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.StreamResponse{
		Stream: stream,
	})
}

func (h *StreamHandler) ListStreams(c *gin.Context) {
	streams := h.converter.ListStreams()
	c.JSON(http.StatusOK, gin.H{"streams": streams})
}

func (h *StreamHandler) StopStream(c *gin.Context) {
	streamID := c.Param("id")

	if err := h.converter.StopStream(streamID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Stream %s stopped successfully", streamID)})
}

func (h *StreamHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"time":   time.Now(),
	})
}
