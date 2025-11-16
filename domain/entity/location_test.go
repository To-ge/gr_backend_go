package entity

import (
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

var (
	location = Location{
		Timestamp: time.Now().Unix(),
		Latitude:  35.3535,
		Longitude: 137.137,
		Altitude:  100.0,
		Speed:     10.0,
	}
)

func isChannelClosed(ch LocationChannel) bool {
	for { // チャネルのバッファを吐き出すまでループ
		select {
		case _, isExisted := <-ch:
			if !isExisted {
				return true
			}
			// もう一度確認
		default:
			// チャネルによるブロッキングを避ける
			return false
		}
	}
}

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../.env"); err != nil {
		println("Warning: Error loading .env file:", err.Error())
	}

	code := m.Run()

	os.Exit(code)
}

func TestNewLiveLocationManager(t *testing.T) {
	llm := NewLiveLocationManager()

	if llm == nil {
		t.Error("NewLiveLocationManager() returned nil")
	}
	if llm.mu == nil {
		t.Error("mutex is nil")
	}
	if llm.LocationList == nil {
		t.Error("LocationList is nil")
	}
	if llm.ChannelList == nil {
		t.Error("ChannelList is nil")
	}
	if llm.ctx == nil {
		t.Error("context is nil")
	}
	if llm.cancel == nil {
		t.Error("cancel function is nil")
	}
}

func TestLiveLocationManager_Add(t *testing.T) {
	llm := NewLiveLocationManager()
	defer llm.StopTimer()

	data := location

	// Addする前の状態を確認
	if len(llm.LocationList) != 0 {
		t.Errorf("Initial LocationList length = %d, want 0", len(llm.LocationList))
	}

	// Locationを追加
	llm.Add(data)

	// goroutineがタイマーを開始するのを待つ
	time.Sleep(10 * time.Millisecond)

	// 追加後の状態を確認
	llm.mu.Lock()
	locationCount := len(llm.LocationList)
	timerExists := llm.timer != nil
	llm.mu.Unlock()

	if locationCount != 1 {
		t.Errorf("LocationList length after Add = %d, want 1", locationCount)
	}
	if !reflect.DeepEqual(llm.LocationList[0], location) {
		t.Errorf("Added location = %v, want %v", llm.LocationList[0], location)
	}
	if !timerExists {
		t.Error("Timer should be started after first Add")
	}
}

func TestLiveLocationManager_AddChannel(t *testing.T) {
	llm := NewLiveLocationManager()
	defer llm.StopTimer()

	ch := make(LocationChannel, 1)

	// チャンネル追加前の状態を確認
	if len(llm.ChannelList) != 0 {
		t.Errorf("Initial ChannelList length = %d, want 0", len(llm.ChannelList))
	}

	// チャンネルを追加
	llm.AddChannel(&ch)

	// 追加後の状態を確認
	llm.mu.Lock()
	if len(llm.ChannelList) != 1 {
		t.Errorf("ChannelList length after AddChannel = %d, want 1", len(llm.ChannelList))
	}
	llm.mu.Unlock()
}

func TestLiveLocationManager_ChannelReceivesLocation(t *testing.T) {
	llm := NewLiveLocationManager()
	defer llm.StopTimer()

	ch := make(LocationChannel, 1)
	llm.AddChannel(&ch)

	data := location

	// Locationを追加
	llm.Add(data)

	// チャンネルからLocationを受信
	select {
	case receivedLocation := <-ch:
		if !reflect.DeepEqual(receivedLocation, data) {
			t.Errorf("Received location = %v, want %v", receivedLocation, data)
		}
	case <-time.After(1 * time.Second):
		t.Error("Timeout: no location received on channel")
	}
}

func TestLiveLocationManager_MultipleChannels(t *testing.T) {
	llm := NewLiveLocationManager()
	defer llm.StopTimer()

	// 複数のチャンネルを作成
	ch1 := make(LocationChannel, 1)
	ch2 := make(LocationChannel, 1)
	ch3 := make(LocationChannel, 1)

	llm.AddChannel(&ch1)
	llm.AddChannel(&ch2)
	llm.AddChannel(&ch3)

	data := location

	// Locationを追加
	llm.Add(data)

	// 全てのチャンネルがLocationを受信するかチェック
	channels := []*LocationChannel{&ch1, &ch2, &ch3}
	var wg sync.WaitGroup
	wg.Add(len(channels))

	for i, ch := range channels {
		go func(index int, channel *LocationChannel) {
			defer wg.Done()
			select {
			case receivedLocation := <-*channel:
				if !reflect.DeepEqual(receivedLocation, data) {
					t.Errorf("Channel %d received location = %v, want %v", index, receivedLocation, data)
				}
			case <-time.After(1 * time.Second):
				t.Errorf("Channel %d timeout: no location received", index)
			}
		}(i, ch)
	}

	wg.Wait()
}

func TestLiveLocationManager_StopTimer(t *testing.T) {
	llm := NewLiveLocationManager()

	// Locationを追加してタイマーを開始
	data := location
	llm.Add(data)

	// 少し待つ（goroutineがタイマーを開始するのを待つ）
	time.Sleep(10 * time.Millisecond)

	// タイマーが開始されていることを確認
	llm.mu.Lock()
	timerStarted := llm.timer != nil
	llm.mu.Unlock()

	if !timerStarted {
		t.Error("Timer should be started after Add")
	}

	// StopTimerを呼び出し
	telemetryLog := llm.StopTimer()

	if telemetryLog == nil {
		t.Error("StopTimer() returned nil")
	}

	// 新しいコンテキストが作成されていることを確認
	if llm.ctx == nil {
		t.Error("Context should be recreated after StopTimer")
	}
	if llm.cancel == nil {
		t.Error("Cancel function should be recreated after StopTimer")
	}
}

func TestLiveLocationManager_EmptyChannelList(t *testing.T) {
	llm := NewLiveLocationManager()
	defer llm.StopTimer()

	// チャンネルを追加
	ch1 := make(LocationChannel, 1)
	ch2 := make(LocationChannel, 1)
	llm.AddChannel(&ch1)
	llm.AddChannel(&ch2)

	// Locationを追加
	data := location
	llm.Add(data)

	time.Sleep(10 * time.Millisecond)

	// チャンネルリストとロケーションリストをクリア
	llm.emptyChannelList()

	// チャンネルがクローズされていることを確認
	if !isChannelClosed(ch1) {
		t.Error("Channel 1 should be closed")
	}

	if !isChannelClosed(ch2) {
		t.Error("Channel 2 should be closed")
	}

	// リストがクリアされていることを確認
	llm.mu.Lock()
	if len(llm.ChannelList) != 0 {
		t.Errorf("ChannelList length after empty = %d, want 0", len(llm.ChannelList))
	}
	if len(llm.LocationList) != 0 {
		t.Errorf("LocationList length after empty = %d, want 0", len(llm.LocationList))
	}
	llm.mu.Unlock()
}

func TestLiveLocationManager_ConcurrentAccess(t *testing.T) {
	llm := NewLiveLocationManager()
	defer llm.StopTimer()

	var wg sync.WaitGroup
	numGoroutines := 10
	numLocationsPerGoroutine := 5

	// 複数のゴルーチンから同時にLocationを追加
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(routineID int) {
			defer wg.Done()
			for j := 0; j < numLocationsPerGoroutine; j++ {
				location := Location{
					Timestamp: time.Now().Unix(),
					Latitude:  float64(routineID),
					Longitude: float64(j),
				}
				llm.Add(location)
			}
		}(i)
	}

	wg.Wait()

	// 期待される総数を確認
	expectedCount := numGoroutines * numLocationsPerGoroutine
	llm.mu.Lock()
	actualCount := len(llm.LocationList)
	llm.mu.Unlock()

	if actualCount != expectedCount {
		t.Errorf("LocationList length = %d, want %d", actualCount, expectedCount)
	}
}

// テスト用のヘルパー関数：短いタイマー期間でテストする場合
func TestLiveLocationManager_TimerReset(t *testing.T) {
	// 注意: このテストは実際のタイマー機能をテストするため、
	// config.LoadConfig().DomainInfo.TimerMinutesが短い値に設定されている必要があります
	// または、テスト用の設定を使用してください

	llm := NewLiveLocationManager()
	defer llm.StopTimer()

	location := Location{
		Timestamp: time.Now().Unix(),
		Latitude:  35.6762,
		Longitude: 139.6503,
	}

	// 最初のLocationを追加
	llm.Add(location)

	// 少し待つ（goroutineがタイマーを開始するのを待つ）
	time.Sleep(10 * time.Millisecond)

	// タイマーが開始されていることを確認
	llm.mu.Lock()
	timerExists := llm.timer != nil
	llm.mu.Unlock()

	if !timerExists {
		t.Error("Timer should be started after first Add")
	}

	// 少し待ってから2つ目のLocationを追加（タイマーリセットをテスト）
	time.Sleep(100 * time.Millisecond)
	llm.Add(location)

	// タイマーがまだ存在することを確認
	llm.mu.Lock()
	timerStillExists := llm.timer != nil
	llm.mu.Unlock()

	if !timerStillExists {
		t.Error("Timer should still exist after reset")
	}
}

// テスト用の関数：チャンネルクローズの動作を確認
func TestChannelCloseWithData(t *testing.T) {
	// バッファ付きチャンネルを作成
	ch := make(LocationChannel, 2)

	// データを2つ送信
	loc1 := Location{Timestamp: 1, Latitude: 35.0, Longitude: 139.0}
	loc2 := Location{Timestamp: 2, Latitude: 36.0, Longitude: 140.0}
	ch <- loc1
	ch <- loc2

	// データがある状態でチャンネルをクローズ
	close(ch)

	t.Log("チャンネルをクローズしました（データ2つがバッファ内にある）")

	// 1回目の読み取り - まだデータがある
	data1, ok1 := <-ch
	t.Logf("1回目: data=%v, ok=%v", data1, ok1)
	if !ok1 {
		t.Error("1回目: okがfalseになるべきではない")
	}
	if !reflect.DeepEqual(data1, loc1) {
		t.Errorf("1回目: 期待値と異なる data=%v, want=%v", data1, loc1)
	}

	// 2回目の読み取り - まだデータがある
	data2, ok2 := <-ch
	t.Logf("2回目: data=%v, ok=%v", data2, ok2)
	if !ok2 {
		t.Error("2回目: okがfalseになるべきではない")
	}
	if !reflect.DeepEqual(data2, loc2) {
		t.Errorf("2回目: 期待値と異なる data=%v, want=%v", data2, loc2)
	}

	// 3回目の読み取り - バッファが空でクローズされている
	data3, ok3 := <-ch
	t.Logf("3回目: data=%v, ok=%v", data3, ok3)
	if ok3 {
		t.Error("3回目: okがtrueになるべきではない")
	}
	if !reflect.DeepEqual(data3, Location{}) {
		t.Errorf("3回目: ゼロ値でない data=%v", data3)
	}

	// 4回目の読み取り - 何度読み取ってもゼロ値とfalse
	data4, ok4 := <-ch
	t.Logf("4回目: data=%v, ok=%v", data4, ok4)
	if ok4 {
		t.Error("4回目: okがtrueになるべきではない")
	}
}
