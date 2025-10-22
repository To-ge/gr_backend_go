package entity

import "time"

type Status int

const (
	Unknown      Status = iota // 不明
	Disconnected               // 地上局が切断した
)

var (
	telemetryLog *TelemetryLog
)

type TimeSpan struct {
	StartTime time.Time
	EndTime   time.Time
}

type TelemetryLog struct {
	ID            uint
	StartTime     time.Time
	EndTime       time.Time
	LocationCount int
	IsPublic      bool
}

func NewTelemetryLog(time time.Time) *TelemetryLog {
	return &TelemetryLog{
		StartTime:     time,
		LocationCount: 0,
	}
}

func (tl *TelemetryLog) IncrementLocationCount() {
	tl.LocationCount++
}
