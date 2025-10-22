package model

import "time"

type GetTelemetryLogsInput struct{}

type GetTelemetryLogsOutput struct {
	IsPublic bool           `json:"is_public"`
	Logs     []TelemetryLog `json:"logs"`
}

type TelemetryLog struct {
	ID            uint      `json:"id"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	LocationCount int       `json:"location_count"`
	IsPublic      bool      `json:"is_public"`
}

type ToggleTelemetryLogVisibilityInput struct {
	Id      uint `json:"id"`
	Visible bool `json:"visible"`
}
