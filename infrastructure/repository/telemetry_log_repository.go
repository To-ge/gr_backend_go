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
