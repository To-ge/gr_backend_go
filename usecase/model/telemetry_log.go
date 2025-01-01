package model

import "time"

type GetTelemetryLogsInput struct{}

type GetTelemetryLogsOutput struct {
	Logs []TelemetryLog `json:"logs"`
}

type TelemetryLog struct {
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	LocationCount int       `json:"location_count"`
}
