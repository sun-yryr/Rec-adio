package health

import (
	"github.com/gin-gonic/gin"
	"github.com/sun-yryr/recoto/api/server/router"
	"github.com/sun-yryr/recoto/internal/broker"
	"go.uber.org/zap"
)

type HealthRouter struct {
	routes []router.Route
	broker broker.Broker
	logger *zap.Logger
}

func NewHealthRouter(broker broker.Broker, logger *zap.Logger) *HealthRouter {
	r := &HealthRouter{
		broker: broker,
		logger: logger,
	}
	r.initRoutes()

	return r
}

func (r *HealthRouter) Routes() []router.Route {
	return r.routes
}

func (r *HealthRouter) initRoutes() {
	r.routes = []router.Route{
		router.NewGetRouter("/health", r.healthCheck),
	}
}

func (r *HealthRouter) healthCheck(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"status": "ok",
	})
}
