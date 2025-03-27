package event

import (
	"sync"
	"time"
)

// Event はシステムイベントを表す構造体です
type Event struct {
	// Type はイベントのタイプです（例: "recording.started"）
	Type string `json:"type"`

	// Timestamp はイベントの発生時刻です
	Timestamp time.Time `json:"timestamp"`

	// Data はイベントのデータです
	Data map[string]interface{} `json:"data"`
}

// EventListener はイベントリスナーのインターフェースです
type EventListener interface {
	// OnEvent はイベント発生時に呼び出されるメソッドです
	OnEvent(event Event) error

	// GetSupportedEvents はリスナーがサポートするイベントタイプのリストを返します
	GetSupportedEvents() []string
}

// EventEmitter はイベントエミッターのインターフェースです
type EventEmitter interface {
	// AddListener はイベントリスナーを追加します
	AddListener(listener EventListener) error

	// RemoveListener はイベントリスナーを削除します
	RemoveListener(listener EventListener) error

	// EmitEvent はイベントを発行します
	EmitEvent(event Event) error
}

// EventManager はイベントの管理を行うマネージャーです
type EventManager struct {
	listeners []EventListener
	mu        sync.RWMutex
}

// NewEventManager は新しいEventManagerインスタンスを作成します
func NewEventManager() *EventManager {
	return &EventManager{
		listeners: make([]EventListener, 0),
	}
}

// AddListener はイベントリスナーを追加します
func (m *EventManager) AddListener(listener EventListener) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.listeners = append(m.listeners, listener)
	return nil
}

// RemoveListener はイベントリスナーを削除します
func (m *EventManager) RemoveListener(listener EventListener) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, l := range m.listeners {
		if l == listener {
			m.listeners = append(m.listeners[:i], m.listeners[i+1:]...)
			break
		}
	}
	return nil
}

// EmitEvent はイベントを発行します
func (m *EventManager) EmitEvent(event Event) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// イベント発生時刻が設定されていない場合は現在時刻を設定
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// 各リスナーにイベントを通知
	for _, listener := range m.listeners {
		// リスナーがサポートするイベントタイプかどうかを確認
		supported := false
		for _, supportedType := range listener.GetSupportedEvents() {
			if supportedType == event.Type || supportedType == "*" {
				supported = true
				break
			}
		}

		// サポートするイベントタイプの場合のみ通知
		if supported {
			// エラーが発生しても処理を継続
			_ = listener.OnEvent(event)
		}
	}

	return nil
}