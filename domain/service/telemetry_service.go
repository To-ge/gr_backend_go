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
	repo repository.ITelemetryRepository
}

func NewTelemetryService(repo repository.ITelemetryRepository) ITelemetryService {
	return &telemetryService{
		repo: repo,
	}
}

func (ts *telemetryService) SendLocation(location entity.Location) error {
	if err := ts.repo.CreateLocation(location); err != nil {
		return err
	}
	entity.LiveLocationManager.Add(location)

	return nil
}

func (ts *telemetryService) Stop() {
	entity.LiveLocationManager.StopTimer()
}
