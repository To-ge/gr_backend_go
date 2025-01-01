package repository

import "github.com/To-ge/gr_backend_go/domain/entity"

type ITelemetryLogRepository interface {
	CreateTelemetryLog(entity.TelemetryLog) error
	GetTelemetryLogs() ([]entity.TelemetryLog, error)
}
