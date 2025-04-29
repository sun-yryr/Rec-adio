package server

import (
	"github.com/gin-gonic/gin"

	"github.com/sun-yryr/recoto/api/server/router"
)

func NewServer(routers ...router.Router) *gin.Engine {
	engine := gin.Default()

	for _, router := range routers {
		for _, route := range router.Routes() {
			engine.Handle(route.Method(), route.Path(), route.Handler())
		}
	}

	return engine
}
