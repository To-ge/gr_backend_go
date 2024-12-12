package usecase

import (
	"github.com/To-ge/gr_backend_go/domain/entity"
	"github.com/To-ge/gr_backend_go/domain/service"
	"github.com/To-ge/gr_backend_go/usecase/model"
)

type ITelemetryUsecase interface {
	SendLocation(*model.SendLocationInput) (*model.SendLocationOutput, error)
	Stop()
}

type telemetryUsecase struct {
	srv service.ITelemetryService
}

func NewTelemetryUsecase(srv service.ITelemetryService) ITelemetryUsecase {
	return &telemetryUsecase{
		srv: srv,
	}
}

func (tu *telemetryUsecase) SendLocation(input *model.SendLocationInput) (*model.SendLocationOutput, error) {
	location := entity.Location{
		Timestamp: input.Timestamp,
		Latitude:  input.Latitude,
		Longitude: input.Longitude,
		Altitude:  input.Altitude,
		Speed:     input.Speed,
	}

	if err := tu.srv.SendLocation(location); err != nil {
		return nil, err
	}

	return &model.SendLocationOutput{}, nil
}

func (tu *telemetryUsecase) Stop() {
	tu.srv.Stop()
}
