package usecase

import (
	"github.com/To-ge/gr_backend_go/domain/entity"
	"github.com/To-ge/gr_backend_go/domain/service"
	"github.com/To-ge/gr_backend_go/usecase/model"
)

type ILocationUsecase interface {
	StreamLiveLocation(*model.StreamLiveLocationInput) (*model.StreamLiveLocationOutput, error)
	StreamArchiveLocation(*model.StreamArchiveLocationInput) (*model.StreamArchiveLocationOutput, error)
}

type locationUsecase struct {
	srv service.ILocationService
}

func NewLocationUsecase(srv service.ILocationService) ILocationUsecase {
	return &locationUsecase{
		srv: srv,
	}
}

func (lu *locationUsecase) StreamLiveLocation(*model.StreamLiveLocationInput) (*model.StreamLiveLocationOutput, error) {
	ch, err := lu.srv.StreamLiveLocation()
	if err != nil {
		return nil, err
	}
	locationChannel := make(chan model.Location)

	go func(locCh chan<- model.Location) {
		defer close(locationChannel)
		for {
			location, ok := <-ch

			if !ok {
				break
			}

			locCh <- model.Location{
				Timestamp: location.Timestamp,
				Latitude:  location.Latitude,
				Longitude: location.Longitude,
				Altitude:  location.Altitude,
				Speed:     location.Speed,
			}
		}
	}(locationChannel)

	return &model.StreamLiveLocationOutput{LocationChannel: locationChannel}, nil
}

func (lu *locationUsecase) StreamArchiveLocation(input *model.StreamArchiveLocationInput) (*model.StreamArchiveLocationOutput, error) {
	span := entity.TimeSpan{
		StartTime: input.StartTime,
		EndTime:   input.EndTime,
	}

	ch, err := lu.srv.StreamArchiveLocation(span)
	if err != nil {
		return nil, err
	}
	locationChannel := make(chan model.Location)

	go func(locCh chan<- model.Location) {
		defer close(locationChannel)
		for {
			location, ok := <-ch

			if !ok {
				break
			}

			locCh <- model.Location{
				Timestamp: location.Timestamp,
				Latitude:  location.Latitude,
				Longitude: location.Longitude,
				Altitude:  location.Altitude,
				Speed:     location.Speed,
			}
		}
	}(locationChannel)

	return &model.StreamArchiveLocationOutput{LocationChannel: locationChannel}, nil
}
