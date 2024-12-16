package model

import "time"

type Location struct {
	Timestamp int64   `json:"timestamp"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
	Speed     float64 `json:"speed"`
}

// Stream Live Location
type StreamLiveLocationInput struct {
}
type StreamLiveLocationOutput struct {
	LocationChannel <-chan Location
}

// Stream Archive Location
type StreamArchiveLocationInput struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}
type StreamArchiveLocationOutput struct {
	LocationChannel <-chan Location
}
