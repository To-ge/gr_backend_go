package entity

import (
	"context"
	"log"
	"sync"
	"time"
)

var (
	LiveLocationManager = NewLiveLocationManager()
	timerMinutes        = 1
)

type Location struct {
	Timestamp int64
	Latitude  float64
	Longitude float64
	Altitude  float64
	Speed     float64
}

type LocationChannel chan Location

type liveLocationManager struct {
	mu           *sync.Mutex
	LocationList []Location
	ChannelList  []*LocationChannel
	timer        *time.Timer
	timerReset   chan struct{}
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewLiveLocationManager() *liveLocationManager {
	ctx, cancel := context.WithCancel(context.Background())
	llm := &liveLocationManager{
		mu:           &sync.Mutex{},
		LocationList: []Location{},
		ChannelList:  []*LocationChannel{},
		timer:        nil,
		timerReset:   make(chan struct{}),
		ctx:          ctx,
		cancel:       cancel,
	}
	// go llm.startTimer()

	return llm
}

func (llm *liveLocationManager) Add(location Location) {
	llm.mu.Lock()
	if llm.timer == nil {
		jst, _ := time.LoadLocation("Asia/Tokyo")
		telemetryLog = NewTelemetryLog(time.Now().In(jst))

		llm.timer = time.NewTimer(time.Duration(timerMinutes) * time.Minute)
		go llm.startTimer()
	}
	telemetryLog.IncrementLocationCount()
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
	// タイマーをリセット
	llm.timerReset <- struct{}{}
}

func (llm *liveLocationManager) startTimer() {
	for {
		llm.mu.Lock()
		timer := llm.timer // タイマーのローカルコピーを取得
		llm.mu.Unlock()

		select {
		case <-llm.timerReset:
			// タイマーをリセット
			llm.mu.Lock()
			if llm.timer != nil {
				if !llm.timer.Stop() {
					<-llm.timer.C // 残っているイベントを消費
				}
			}
			llm.timer = time.NewTimer(time.Duration(timerMinutes) * time.Minute)
			log.Printf("Timer reset to %d minutes\n", timerMinutes)
			llm.mu.Unlock()
		case <-timer.C:
			// タイムアウト処理
			llm.mu.Lock()
			if llm.timer != nil {
				log.Printf("Timer expired: no Add call for %d minutes\n", timerMinutes)
				llm.timer = nil // タイマーをクリア
				llm.emptyChannelList()
			}
			llm.mu.Unlock()
			return
		case <-llm.ctx.Done():
			// キャンセル処理
			llm.mu.Lock()
			if llm.timer != nil {
				if !llm.timer.Stop() {
					<-llm.timer.C // 残っているイベントを消費
				}
			}
			llm.timer = nil // タイマーをクリア
			llm.emptyChannelList()
			log.Println("Timer have been stopped!")
			llm.mu.Unlock()
			return
		}
	}
}

func (llm *liveLocationManager) AddChannel(ch *LocationChannel) {
	llm.mu.Lock()
	defer llm.mu.Unlock()

	llm.ChannelList = append(llm.ChannelList, ch)
}

func (llm *liveLocationManager) emptyChannelList() {
	for _, v := range llm.ChannelList {
		close(*v)
	}
	llm.ChannelList = []*LocationChannel{}
	llm.LocationList = []Location{}
}

func (llm *liveLocationManager) StopTimer() *TelemetryLog {
	llm.mu.Lock()
	defer llm.mu.Unlock()

	if llm.cancel != nil {
		llm.cancel() // ゴルーチンを終了
	}
	// 新しいコンテキストを作成
	ctx, cancel := context.WithCancel(context.Background())
	llm.ctx = ctx
	llm.cancel = cancel

	jst, _ := time.LoadLocation("Asia/Tokyo")
	telemetryLog.EndTime = time.Now().In(jst)
	return telemetryLog
}
