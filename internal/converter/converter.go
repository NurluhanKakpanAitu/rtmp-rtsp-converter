package converter

import (
	"context"
	"fmt"
	"os/exec"
	"sync"
	"time"

	"rtmp-rtsp-converter/internal/logger"
	"rtmp-rtsp-converter/pkg/models"
)

type Converter struct {
	streams    map[string]*models.Stream
	processes  map[string]*exec.Cmd
	mu         sync.RWMutex
	rtspHost   string
	rtspPort   int
	maxStreams int
}

func NewConverter(rtspHost string, rtspPort, maxStreams int) *Converter {
	return &Converter{
		streams:    make(map[string]*models.Stream),
		processes:  make(map[string]*exec.Cmd),
		rtspHost:   rtspHost,
		rtspPort:   rtspPort,
		maxStreams: maxStreams,
	}
}

func (c *Converter) StartStream(ctx context.Context, streamID, rtmpUrl string) (*models.Stream, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.streams) >= c.maxStreams {
		return nil, fmt.Errorf("maximum number of streams reached: %d", c.maxStreams)
	}

	if _, exists := c.streams[streamID]; exists {
		return nil, fmt.Errorf("stream with ID %s already exists", streamID)
	}

	rtspUrl := fmt.Sprintf("rtsp://%s:%d/%s", c.rtspHost, c.rtspPort, streamID)

	stream := &models.Stream{
		ID:        streamID,
		RTMPUrl:   rtmpUrl,
		RTSPUrl:   rtspUrl,
		Status:    models.StatusStarting,
		StartedAt: time.Now(),
	}

	c.streams[streamID] = stream

	// Start FFmpeg process
	go c.runFFmpeg(ctx, stream)

	return stream, nil
}

func (c *Converter) StopStream(streamID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	stream, exists := c.streams[streamID]
	if !exists {
		return fmt.Errorf("stream with ID %s not found", streamID)
	}

	if process, exists := c.processes[streamID]; exists {
		if err := process.Process.Kill(); err != nil {
			logger.GetLogger().WithError(err).Errorf("Failed to kill process for stream %s", streamID)
		}
		delete(c.processes, streamID)
	}

	now := time.Now()
	stream.StoppedAt = &now
	stream.Status = models.StatusStopped

	delete(c.streams, streamID)

	return nil
}

func (c *Converter) GetStream(streamID string) (*models.Stream, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	stream, exists := c.streams[streamID]
	if !exists {
		return nil, fmt.Errorf("stream with ID %s not found", streamID)
	}

	return stream, nil
}

func (c *Converter) ListStreams() []*models.Stream {
	c.mu.RLock()
	defer c.mu.RUnlock()

	streams := make([]*models.Stream, 0, len(c.streams))
	for _, stream := range c.streams {
		streams = append(streams, stream)
	}

	return streams
}

func (c *Converter) runFFmpeg(ctx context.Context, stream *models.Stream) {
	log := logger.GetLogger()

	// FFmpeg command for RTMP to RTSP conversion
	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-i", stream.RTMPUrl,
		"-c", "copy",
		"-f", "rtsp",
		"-rtsp_transport", "tcp",
		stream.RTSPUrl,
	)

	c.mu.Lock()
	c.processes[stream.ID] = cmd
	c.mu.Unlock()

	log.Infof("Starting FFmpeg for stream %s: %s -> %s", stream.ID, stream.RTMPUrl, stream.RTSPUrl)

	if err := cmd.Start(); err != nil {
		log.WithError(err).Errorf("Failed to start FFmpeg for stream %s", stream.ID)
		c.mu.Lock()
		stream.Status = models.StatusError
		stream.ErrorMsg = err.Error()
		c.mu.Unlock()
		return
	}

	c.mu.Lock()
	stream.Status = models.StatusRunning
	c.mu.Unlock()

	if err := cmd.Wait(); err != nil {
		log.WithError(err).Errorf("FFmpeg process failed for stream %s", stream.ID)
		c.mu.Lock()
		stream.Status = models.StatusError
		stream.ErrorMsg = err.Error()
		now := time.Now()
		stream.StoppedAt = &now
		c.mu.Unlock()
	} else {
		log.Infof("FFmpeg process completed for stream %s", stream.ID)
		c.mu.Lock()
		stream.Status = models.StatusStopped
		now := time.Now()
		stream.StoppedAt = &now
		c.mu.Unlock()
	}

	c.mu.Lock()
	delete(c.processes, stream.ID)
	c.mu.Unlock()
}
