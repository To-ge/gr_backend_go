package service

import (
	"github.com/To-ge/gr_backend_go/domain/entity"
	"github.com/To-ge/gr_backend_go/domain/repository"
)

type ITelemetryService interface {
	SendLocation(entity.Location) error
	Stop()
}

type telemetryService struct {
	telemetryRepo    repository.ITelemetryRepository
	telemetryLogRepo repository.ITelemetryLogRepository
}

func NewTelemetryService(telemetryRepo repository.ITelemetryRepository, telemetryLogRepo repository.ITelemetryLogRepository) ITelemetryService {
	return &telemetryService{
		telemetryRepo:    telemetryRepo,
		telemetryLogRepo: telemetryLogRepo,
	}
}

func (ts *telemetryService) SendLocation(location entity.Location) error {
	if err := ts.telemetryRepo.CreateLocation(location); err != nil {
		return err
	}
	entity.LiveLocationManager.Add(location)

	return nil
}

func (ts *telemetryService) Stop() {
	telemetryLog := entity.LiveLocationManager.StopTimer()
	ts.telemetryLogRepo.CreateTelemetryLog(*telemetryLog)
}
