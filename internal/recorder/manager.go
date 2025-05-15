// Package recorder は録音の管理を行うパッケージ
package recorder

import (
	"context"
	"slices"
	"sync"
	"time"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"

	"github.com/sun-yryr/recoto/internal/domain"
	"github.com/sun-yryr/recoto/internal/event/recording"
	"github.com/sun-yryr/recoto/internal/eventutil"
)

// Recorder は録音を行うインターフェース。
type Recorder interface {
	GetSupportSource() []domain.SourceKind
	Rec(
		ctx context.Context,
		event *recording.RequestedEvent,
	) error
	GetName() string
	CheckAvailable() error
}

type recorderWithStatus struct {
	recorder  Recorder
	available bool
	error     error
}

// RecordingManager は録音の管理を行うマネージャー。
type RecordingManager struct {
	requestedService *eventutil.EventService[recording.RequestedEvent]
	startedService   *eventutil.EventService[recording.StartedEvent]
	finishedService  *eventutil.EventService[recording.FinishedEvent]
	logger           *zap.Logger
	mu               sync.Mutex
	recorders        []recorderWithStatus
	cancels          map[string]context.CancelFunc
}

// NewRecordingManager は録音の管理を行うマネージャーを生成する。
func NewRecordingManager(
	requestedService *eventutil.EventService[recording.RequestedEvent],
	startedService *eventutil.EventService[recording.StartedEvent],
	finishedService *eventutil.EventService[recording.FinishedEvent],
	logger *zap.Logger,
) *RecordingManager {
	return &RecordingManager{
		requestedService: requestedService,
		startedService:   startedService,
		finishedService:  finishedService,
		logger:           logger,
		mu:               sync.Mutex{},
		recorders:        []recorderWithStatus{},
		cancels:          make(map[string]context.CancelFunc),
	}
}

// AddRecorder はRecorderに状態を付与してRecordingManagerに追加する。
func (m *RecordingManager) AddRecorder(recorder Recorder) error {
	rws := recorderWithStatus{
		recorder:  recorder,
		available: true,
		error:     nil,
	}

	// 有効でない場合はエラーを格納し、無効な状態にする
	if err := recorder.CheckAvailable(); err != nil {
		rws.available = false
		rws.error = err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.recorders = append(m.recorders, rws)

	return rws.error
}

// Start は録音マネージャーを起動し、イベント駆動で録音を管理する。
func (m *RecordingManager) Start(ctx context.Context) error {
	unsubscribeFunc, err := m.requestedService.Subscribe(
		ctx,
		m.handleRequestedEvent,
	)
	if err != nil {
		return errors.Wrap(err, "failed to subscribe to requested events")
	}

	// コンテキストが終了したら録音を停止する。
	go func() {
		<-ctx.Done()
		m.logger.Info("recording manager is stopped")
		m.mu.Lock()
		defer m.mu.Unlock()

		for _, cancel := range m.cancels {
			cancel()
		}

		if err := unsubscribeFunc(); err != nil {
			m.logger.Error("failed to unsubscribe from requested events", zap.Error(err))
		}
	}()

	return nil
}

func (m *RecordingManager) handleRequestedEvent(
	ctx context.Context,
	event *recording.RequestedEvent,
) {
	// サポートしているrecorderを探す
	for _, recorderWithStatus := range m.recorders {
		if !recorderWithStatus.available {
			continue
		}

		recorder := recorderWithStatus.recorder
		if len(recorder.GetSupportSource()) == 0 {
			continue
		}

		if !slices.Contains(recorder.GetSupportSource(), event.Source.Kind) {
			continue
		}

		// 録音処理は別のgoroutineで行う
		go func() {
			recCtx, cancel := context.WithCancel(ctx)

			m.mu.Lock()
			m.cancels[event.RecordingID] = cancel
			m.mu.Unlock()

			defer func() {
				m.mu.Lock()
				delete(m.cancels, event.RecordingID)
				m.mu.Unlock()
			}()

			if err := m.beforeRec(recCtx, event); err != nil {
				m.logger.Error("failed to process before recording", zap.Error(err))

				return
			}

			recErr := recorder.Rec(recCtx, event)
			if recErr != nil {
				m.logger.Error("failed to record", zap.Error(recErr))
			}

			if err := m.afterRec(recCtx, event, recErr); err != nil {
				m.logger.Error("failed to process after recording", zap.Error(err))
			}
		}()

		return
	}

	// サポートしているrecorderが見つからなかった場合はエラーを出力する
	m.logger.Error(
		"no recorder found",
		zap.String("recordingId", event.RecordingID),
		zap.String("source", string(event.Source.Kind)),
		zap.String("sourceId", event.Source.ID),
	)
}

func (m *RecordingManager) beforeRec(ctx context.Context, event *recording.RequestedEvent) error {
	if err := m.startedService.Publish(ctx, recording.StartedEvent{
		RecordingID: event.RecordingID,
		Timestamp:   time.Now(),
	}); err != nil {
		return errors.Wrap(err, "failed to publish started event")
	}

	return nil
}

func (m *RecordingManager) afterRec(
	ctx context.Context,
	event *recording.RequestedEvent,
	recErr error,
) error {
	var errorMsg *string

	if recErr != nil {
		errStr := recErr.Error()
		errorMsg = &errStr
	}

	if err := m.finishedService.Publish(ctx, recording.FinishedEvent{
		RecordingID: event.RecordingID,
		Success:     recErr == nil,
		Error:       errorMsg,
		Timestamp:   time.Now(),
	}); err != nil {
		return errors.Wrap(err, "failed to publish finished event")
	}

	return nil
}
