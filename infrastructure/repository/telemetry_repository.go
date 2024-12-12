package repository

import (
	"fmt"

	"github.com/To-ge/gr_backend_go/domain/entity"
	domainRepo "github.com/To-ge/gr_backend_go/domain/repository"
	"github.com/To-ge/gr_backend_go/infrastructure/database"
	"github.com/To-ge/gr_backend_go/infrastructure/database/model"
)

type telemetryRepository struct {
	dbc *database.DBConnector
}

func NewtelemetryRepository(dbc *database.DBConnector) domainRepo.ITelemetryRepository {
	return &telemetryRepository{
		dbc: dbc,
	}
}

func (tr *telemetryRepository) CreateLocation(location entity.Location) error {
	record := &model.Location{
		Timestamp: location.Timestamp,
		Latitude:  location.Latitude,
		Longitude: location.Longitude,
		Altitude:  location.Altitude,
		Speed:     location.Speed,
	}

	if err := tr.dbc.Conn.Create(record).Error; err != nil {
		return fmt.Errorf("new location can't create, %s", err.Error())
	}
	return nil

}
