package recordings

import (
	"go.uber.org/zap"

	"github.com/sun-yryr/recoto/api/server/router"
	"github.com/sun-yryr/recoto/internal/broker"
)

type RecordingsRouter struct {
	routes []router.Route
	broker broker.Broker
	logger *zap.Logger
}

func NewRecordingsRouter(broker broker.Broker, logger *zap.Logger) *RecordingsRouter {
	r := &RecordingsRouter{
		broker: broker,
		logger: logger,
	}
	r.initRoutes()

	return r
}

func (r *RecordingsRouter) Routes() []router.Route {
	return r.routes
}

func (r *RecordingsRouter) initRoutes() {
	r.routes = []router.Route{}
}
