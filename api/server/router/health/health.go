package health

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/sun-yryr/recoto/api/server/router"
	"github.com/sun-yryr/recoto/internal/broker"
)

type HealthChecker interface {
	GetName() string
	Check(ctx context.Context) error
}

type HealthRouter struct {
	routes         []router.Route
	broker         broker.Broker
	logger         *zap.Logger
	healthCheckers []HealthChecker
}

func NewHealthRouter(
	broker broker.Broker,
	logger *zap.Logger,
	healthCheckers ...HealthChecker,
) *HealthRouter {
	r := &HealthRouter{
		broker:         broker,
		logger:         logger,
		healthCheckers: healthCheckers,
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
	type checkResult struct {
		Name  string `json:"name"`
		Error string `json:"error,omitempty"`
		OK    bool   `json:"ok"`
	}

	results := make([]checkResult, 0, len(r.healthCheckers))

	resultCh := make(chan checkResult, len(r.healthCheckers))
	defer close(resultCh)

	// 全てのヘルスチェックを非同期で実行
	for _, checker := range r.healthCheckers {
		go func(c HealthChecker) {
			result := checkResult{
				Name: c.GetName(),
				OK:   true,
			}

			if err := c.Check(ctx); err != nil {
				result.OK = false
				result.Error = err.Error()
				r.logger.Error(
					"Health check failed",
					zap.String("checker", c.GetName()),
					zap.Error(err),
				)
			}

			resultCh <- result
		}(checker)
	}

	// 全ての結果を待機・収集
	for range len(r.healthCheckers) {
		results = append(results, <-resultCh)
	}

	// 全体の状態を判断
	allOK := true

	for _, result := range results {
		if !result.OK {
			allOK = false

			break
		}
	}

	if allOK {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"results": results,
		})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"results": results,
		})
	}
}
