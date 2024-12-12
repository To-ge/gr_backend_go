package repository

import "github.com/To-ge/gr_backend_go/domain/entity"

type ITelemetryRepository interface {
	CreateLocation(entity.Location) error
}
