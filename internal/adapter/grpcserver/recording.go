package grpcserver

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"

	config "github.com/sun-yryr/recoto/internal/config/server"
	"github.com/sun-yryr/recoto/internal/domain"
	"github.com/sun-yryr/recoto/internal/event/recording"
	"github.com/sun-yryr/recoto/internal/eventutil"
	"github.com/sun-yryr/recoto/internal/fileutil"
	"github.com/sun-yryr/recoto/internal/logger"
	recordingv1 "github.com/sun-yryr/recoto/pkg/api/recoto/recording/v1"
)

// RecordingService は、RecordingService Interfaceを実装した構造体。
type RecordingService struct {
	recordingv1.UnimplementedRecordingServiceServer
	requestedService *eventutil.EventService[recording.RequestedEvent]
	saveDir          string
}

// NewRecordingService は、RecordingServiceのコンストラクタ。
func NewRecordingService(
	cfg *config.Config,
	requestedService *eventutil.EventService[recording.RequestedEvent],
) *RecordingService {
	return &RecordingService{
		requestedService: requestedService,
		saveDir:          cfg.Recording.SaveDir,
	}
}

// StartFromURL は、URLから録音を開始する。
func (s *RecordingService) StartFromURL(
	ctx context.Context,
	req *recordingv1.StartFromURLRequest,
) (*recordingv1.StartFromURLResponse, error) {
	log := logger.FromContextWithTrace(ctx)

	log.Info("starting recording from URL",
		zap.String("url", req.GetUrl()),
		zap.String("title", req.GetTitle()),
	)

	// タイトルをサニタイズ
	safeTitle := fileutil.SanitizeFilename(req.GetTitle())

	filename := filepath.Join(s.saveDir, safeTitle+".m4a")
	// すでに存在する場合はunixtimeを付与
	if _, err := os.Stat(filename); err == nil {
		filename = filepath.Join(
			s.saveDir,
			safeTitle+"_"+time.Now().Format("20060102150405")+".m4a",
		)
	}

	source, err := domain.NewURLSource(req.GetUrl())
	if err != nil {
		log.Error("failed to create url source", zap.Error(err))

		return nil, errors.Wrap(err, "failed to create url source")
	}

	newRecording, err := domain.NewRecording(source, filename, req.GetDuration().AsDuration())
	if err != nil {
		log.Error("failed to create recording model", zap.Error(err))

		return nil, errors.Wrap(err, "failed to create recording model")
	}

	// イベント発行 (TraceID が自動的に伝播される)
	if err := s.requestedService.Publish(ctx, *recording.NewRequestedEvent(*newRecording)); err != nil {
		log.Error("failed to publish requested event", zap.Error(err))

		return nil, errors.Wrap(err, "failed to publish requested event")
	}

	log.Info("recording request published successfully",
		zap.String("recording_id", string(newRecording.ID)),
	)

	return &recordingv1.StartFromURLResponse{
		RecordingId: string(newRecording.ID),
	}, nil
}
