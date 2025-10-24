//go:generate mockgen -source=telemetry_repository.go -destination=../mock_repository/mock_telemetry_repository.go -package=mock_repository

package repository

import "github.com/To-ge/gr_backend_go/domain/entity"

type ITelemetryRepository interface {
	CreateLocation(entity.Location) error
}
