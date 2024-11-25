package model

import "time"

type Location struct {
	Timestamp time.Time `json:"timestamp"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Altitude  float64   `json:"altitude"`
	Speed     float64   `json:"speed"`
}

// Stream Live Location
type StreamLiveLocationInput struct {
}
type StreamLiveLocationOutput struct {
	LocationChannel <-chan Location
}

// Stream Archive Location
type StreamArchiveLocationInput struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}
type StreamArchiveLocationOutput struct {
	LocationChannel <-chan Location
}
