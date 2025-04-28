package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sun-yryr/recoto/api/server/router"
)

func NewServer(routers ...router.Router) *gin.Engine {
	r := gin.Default()

	for _, router := range routers {
		for _, route := range router.Routes() {
			r.Handle(route.Method(), route.Path(), route.Handler())
		}
	}

	return r
}
