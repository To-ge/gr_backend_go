package repository

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/To-ge/gr_backend_go/domain/entity"
	domainRepo "github.com/To-ge/gr_backend_go/domain/repository"
	"github.com/To-ge/gr_backend_go/infrastructure/database"
	"github.com/To-ge/gr_backend_go/infrastructure/database/model"
)

type locationRepository struct {
	dbc *database.DBConnector
}

func NewLocationRepository(dbc *database.DBConnector) domainRepo.ILocationRepository {
	return locationRepository{
		dbc: dbc,
	}
}

func (lr locationRepository) StreamArchiveLocation(span entity.TimeSpan) (*entity.LocationChannel, error) {
	// startTime := unixTimeToTime(span.StartTime)
	// endTime := unixTimeToTime(span.EndTime)
	locations := []model.Location{}

	if err := lr.dbc.Conn.Model(model.Location{}).Where("created_at >= ? AND created_at <= ?", span.StartTime, span.EndTime).Find(&locations).Error; err != nil {
		return nil, fmt.Errorf("archive location can't find, %s", err.Error())
	}

	ch := make(entity.LocationChannel)
	interval, _ := strconv.Atoi(os.Getenv("TX_INTERVAL"))

	go func() {
		defer close(ch)

		time.Sleep(2 * time.Second)
		for _, v := range locations {
			time.Sleep(time.Duration(interval) * time.Millisecond)
			ch <- entity.Location{
				Timestamp: v.Timestamp,
				Latitude:  v.Latitude,
				Longitude: v.Longitude,
				Altitude:  v.Altitude,
				Speed:     v.Speed,
			}
		}
	}()

	return &ch, nil
}

func unixTimeToTime(unixTime int64) time.Time {
	// 秒部分とナノ秒部分に分割
	sec := unixTime / int64(time.Second)
	nsec := unixTime % int64(time.Second)

	return time.Unix(sec, nsec).UTC()
}
