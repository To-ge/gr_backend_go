package service

import (
	"github.com/To-ge/gr_backend_go/domain/entity"
	"github.com/To-ge/gr_backend_go/domain/repository"
)

type ILocationService interface {
	StreamLiveLocation() (*entity.LocationChannel, error)
	StreamArchiveLocation(entity.TimeSpan) (entity.LocationChannel, error)
}

type locationService struct {
	repo repository.ILocationRepository
}

func NewLocationService(lr repository.ILocationRepository) ILocationService {
	return &locationService{
		repo: lr,
	}
}

func (ls *locationService) StreamLiveLocation() (*entity.LocationChannel, error) {
	ch := new(entity.LocationChannel)
	*ch = make(entity.LocationChannel, 10)
	locationList := entity.LiveLocationManager.LocationList

	for _, v := range locationList {
		*ch <- v
	}
	entity.LiveLocationManager.AddChannel(ch)
	return ch, nil
}

func (ls *locationService) StreamArchiveLocation(span entity.TimeSpan) (entity.LocationChannel, error) {
	ch, err := ls.repo.StreamArchiveLocation(span)
	if err != nil {
		return nil, err
	}
	return ch, nil
}
