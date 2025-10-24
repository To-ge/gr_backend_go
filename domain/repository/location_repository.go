//go:generate mockgen -source=location_repository.go -destination=../mock_repository/mock_location_repository.go -package=mock_repository

package repository

import "github.com/To-ge/gr_backend_go/domain/entity"

type ILocationRepository interface {
	StreamArchiveLocation(entity.TimeSpan) (*entity.LocationChannel, error)
}
