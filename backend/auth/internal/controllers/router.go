package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine, services Services, configs Configs) *gin.Engine {
	// config.InitCorsConfig()
	// r.Use(cors.New(config.CorsConfig))

	// Service health check
	r.GET("/api/auth/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "healthy"})
	})

	InitAuthRouter(r, services.AuthService, configs.AuthHandlerConfig)
	// InitUserHandler(r, services.UserService, configs.UserHandlerConfig)

	return r
}
