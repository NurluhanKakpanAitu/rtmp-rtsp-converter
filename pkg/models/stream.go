package models

import (
    "time"
)

type StreamStatus string

const (
    StatusStarting StreamStatus = "starting"
    StatusRunning  StreamStatus = "running"
    StatusStopped  StreamStatus = "stopped"
    StatusError    StreamStatus = "error"
)

type Stream struct {
    ID          string       `json:"id"`
    RTMPUrl     string       `json:"rtmp_url"`
    RTSPUrl     string       `json:"rtsp_url"`
    Status      StreamStatus `json:"status"`
    StartedAt   time.Time    `json:"started_at"`
    StoppedAt   *time.Time   `json:"stopped_at,omitempty"`
    ErrorMsg    string       `json:"error_msg,omitempty"`
}

type CreateStreamRequest struct {
    RTMPUrl string `json:"rtmp_url" binding:"required"`
    StreamID string `json:"stream_id,omitempty"`
}

type StreamResponse struct {
    Stream *Stream `json:"stream"`
    Message string `json:"message"`
}