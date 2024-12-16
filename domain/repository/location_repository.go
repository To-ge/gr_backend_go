package repository

import "github.com/To-ge/gr_backend_go/domain/entity"

type ILocationRepository interface {
	StreamArchiveLocation(entity.TimeSpan) (*entity.LocationChannel, error)
}
