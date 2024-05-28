package controllers

import (
	"net/http"
	"photobox-api/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine, services Services, configs Configs, middles Middles, mq MQ) *gin.Engine {
	config.InitCorsConfig()
	r.Use(cors.New(config.CorsConfig))

	// API GATEWAY health check
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "api healthy"})
	})

	// /api/auth
	InitAuthHandler(r, services.AuthClient)

	// /api/users
	InitUserHandler(r, services.UserClient, configs.UserConfig, middles.Middleware)

	// /api/meta
	InitMetaHandler(r, services.MetaClient, middles.Middleware, mq)

	return r
}
