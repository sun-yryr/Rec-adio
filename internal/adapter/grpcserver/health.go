// Package grpcserver は、gRPCサーバーの実装を提供します.
package grpcserver

import (
	"context"

	"go.uber.org/zap"

	"github.com/sun-yryr/recoto/internal/broker"
	"github.com/sun-yryr/recoto/internal/logger"
	pb "github.com/sun-yryr/recoto/pkg/api/health/v1"
)

type healthChecker interface {
	GetName() string
	Check(ctx context.Context) error
}

// HealthService は、HealthService Interfaceを実装した構造体.
type HealthService struct {
	pb.UnimplementedHealthServiceServer
	broker         broker.Broker
	healthCheckers []healthChecker
}

// NewHealthService は、HealthServiceのコンストラクタ.
func NewHealthService(
	broker broker.Broker,
	healthCheckers ...healthChecker,
) *HealthService {
	return &HealthService{
		broker:         broker,
		healthCheckers: healthCheckers,
	}
}

// Check は、Brokerのヘルスチェックを行う.
func (s *HealthService) Check(
	ctx context.Context,
	_ *pb.HealthCheckRequest,
) (*pb.HealthCheckResponse, error) {
	logger := logger.FromContext(ctx)
	results := make([]*pb.CheckResult, 0, len(s.healthCheckers))

	resultCh := make(chan *pb.CheckResult, len(s.healthCheckers))
	defer close(resultCh)

	// 全てのヘルスチェックを非同期で実行
	for _, checker := range s.healthCheckers {
		go func(checker healthChecker) {
			result := &pb.CheckResult{
				Name: checker.GetName(),
				Ok:   true,
			}

			if err := checker.Check(ctx); err != nil {
				result.Ok = false
				result.Error = err.Error()
				logger.Error(
					"Health check failed",
					zap.String("checker", checker.GetName()),
					zap.Error(err),
				)
			}

			resultCh <- result
		}(checker)
	}

	// 全ての結果を待機・収集
	for range s.healthCheckers {
		results = append(results, <-resultCh)
	}

	// 全体の状態を判断
	allOK := true

	for _, result := range results {
		if !result.GetOk() {
			allOK = false

			break
		}
	}

	status := "ok"
	if !allOK {
		status = "error"
	}

	return &pb.HealthCheckResponse{
		Status:  status,
		Results: results,
	}, nil
}
