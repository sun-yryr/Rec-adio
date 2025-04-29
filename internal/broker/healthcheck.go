package broker

import (
	"context"
	"time"
)

const timeout = 3 * time.Second

type BrokerHealthCheck struct {
	Broker Broker
}

// BrokerHealthCheck は、Brokerの健康状態を確認するための構造体.
func NewBrokerHealthCheck(broker Broker) *BrokerHealthCheck {
	return &BrokerHealthCheck{
		Broker: broker,
	}
}

func (h *BrokerHealthCheck) GetName() string {
	return "BrokerHealthCheck"
}

// Check は、BrokerのPublish, Subscribeが正常に動作することを確認する関数.
func (h *BrokerHealthCheck) Check(ctx context.Context) error {
	cctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	healthCheckSubject := "recoto.health.check.v1"

	resultCh := make(chan struct{}, 1)
	defer close(resultCh)

	unsubscribe, err := h.Broker.Subscribe(cctx, healthCheckSubject, func(message []byte) {
		resultCh <- struct{}{}
	})
	if err != nil {
		return err
	}

	defer func() {
		if unsubscribe != nil {
			unsubscribe()
		}
	}()

	err = h.Broker.Publish(cctx, healthCheckSubject, []byte("ping"))
	if err != nil {
		return err
	}

	select {
	case <-resultCh:
		return nil
	case <-cctx.Done():
		return cctx.Err()
	}
}
