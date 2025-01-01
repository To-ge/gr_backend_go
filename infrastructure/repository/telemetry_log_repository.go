package repository

import (
	"fmt"

	"github.com/To-ge/gr_backend_go/domain/entity"
	domainRepo "github.com/To-ge/gr_backend_go/domain/repository"
	"github.com/To-ge/gr_backend_go/infrastructure/database"
	"github.com/To-ge/gr_backend_go/infrastructure/database/model"
)

type telemetryLogRepository struct {
	dbc *database.DBConnector
}

func NewTelemetryLogRepository(dbc *database.DBConnector) domainRepo.ITelemetryLogRepository {
	return &telemetryLogRepository{
		dbc: dbc,
	}
}

func (tlr *telemetryLogRepository) CreateTelemetryLog(tl entity.TelemetryLog) error {
	record := &model.TelemetryLog{
		StartTime:     tl.StartTime,
		EndTime:       tl.EndTime,
		LocationCount: tl.LocationCount,
	}

	if err := tlr.dbc.Conn.Create(record).Error; err != nil {
		return fmt.Errorf("new telemetry_log can't create, %s", err.Error())
	}
	return nil
}

func (tlr *telemetryLogRepository) GetTelemetryLogs() ([]entity.TelemetryLog, error) {
	var telemetryLogs []model.TelemetryLog
	if err := tlr.dbc.Conn.Find(&telemetryLogs).Error; err != nil {
		return nil, fmt.Errorf("archive location can't find, %s", err.Error())
	}

	var result []entity.TelemetryLog
	for _, v := range telemetryLogs {
		result = append(result, entity.TelemetryLog{
			StartTime:     v.StartTime,
			EndTime:       v.EndTime,
			LocationCount: v.LocationCount,
		})
	}
	return result, nil
}
