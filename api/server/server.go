package server

import (
	"github.com/gin-gonic/gin"

	"github.com/sun-yryr/recoto/api/server/router"
)

func NewServer(routers ...router.Router) *gin.Engine {
	engine := gin.Default()

	// v1 APIグループを作成
	v1 := engine.Group("/api/v1")

	for _, router := range routers {
		for _, route := range router.Routes() {
			v1.Handle(route.Method(), route.Path(), route.Handler())
		}
	}

	return engine
}
