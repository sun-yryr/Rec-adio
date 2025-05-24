package broker

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"

	"github.com/sun-yryr/recoto/internal/logger"
)

const timeout = 3 * time.Second

// HealthCheck は、Brokerの健康状態を確認するための構造体。
type HealthCheck struct {
	Broker Broker
}

// NewHealthCheck は、HealthCheckのコンストラクタ。
func NewHealthCheck(broker Broker) *HealthCheck {
	return &HealthCheck{
		Broker: broker,
	}
}

// GetName は、HealthCheckの名前を返す関数。
func (h *HealthCheck) GetName() string {
	return "BrokerHealthCheck"
}

// Check は、BrokerのPublish, Subscribeが正常に動作することを確認する関数。
func (h *HealthCheck) Check(ctx context.Context) error {
	cctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	log := logger.FromContextWithTrace(cctx)

	healthCheckSubject := "recoto.health.check.v1"

	resultCh := make(chan struct{}, 1)
	defer close(resultCh)

	unsubscribe, err := h.Broker.Subscribe(
		cctx,
		healthCheckSubject,
		func(_ context.Context, _ []byte) {
			resultCh <- struct{}{}
		},
	)
	if err != nil {
		return errors.Wrap(err, "failed to subscribe health check message")
	}

	defer func() {
		if unsubscribe != nil {
			if err := unsubscribe(); err != nil {
				log.Error("failed to unsubscribe health check message", zap.Error(err))
			}
		}
	}()

	err = h.Broker.Publish(cctx, healthCheckSubject, "ping")
	if err != nil {
		return errors.Wrap(err, "failed to publish health check message")
	}

	select {
	case <-resultCh:
		return nil
	case <-cctx.Done():
		return errors.Wrap(cctx.Err(), "health check timeout")
	}
}
