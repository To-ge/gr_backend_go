package entity

import (
	"log"
	"sync"
	"time"
)

var (
	LiveLocationManager = NewLiveLocationManager()
)

type Location struct {
	Timestamp time.Time
	Latitude  float64
	Longitude float64
	Altitude  float64
	Speed     float64
}

type LocationChannel chan Location

type TimeSpan struct {
	StartTime int64
	EndTime   int64
}

type liveLocationManager struct {
	mu           sync.Mutex
	LocationList []Location
	ChannelList  []*LocationChannel
}

func NewLiveLocationManager() *liveLocationManager {
	return &liveLocationManager{
		LocationList: []Location{},
		ChannelList:  []*LocationChannel{},
	}
}

func (llm *liveLocationManager) Add(location Location) {
	llm.mu.Lock()
	llm.LocationList = append(llm.LocationList, location)
	llm.mu.Unlock()

	for _, ch := range llm.ChannelList {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Println("Recovered from panic:", r)
				}
			}()
			*ch <- location
		}()
	}
}
