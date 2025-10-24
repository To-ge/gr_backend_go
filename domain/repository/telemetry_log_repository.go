//go:generate mockgen -source=telemetry_log_repository.go -destination=../mock_repository/mock_telemetry_log_repository.go -package=mock_repository

package repository

import "github.com/To-ge/gr_backend_go/domain/entity"

type ITelemetryLogRepository interface {
	CreateTelemetryLog(entity.TelemetryLog) error
	GetTelemetryLogs() ([]entity.TelemetryLog, error)
	GetPublicTelemetryLogs() ([]entity.TelemetryLog, error)
	ToggleTelemetryLogVisibility(id uint, visible bool) error
}
