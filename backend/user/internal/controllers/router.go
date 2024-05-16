package controllers

import (
	"net/http"
	"photobox-user/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine, services Services, configs Configs) *gin.Engine {
	config.InitCorsConfig()
	r.Use(cors.New(config.CorsConfig))

	// Service health check
	r.GET("/api/users/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "healthy"})
	})

	InitUserHandler(r, services.UserService, configs.UserHandlerConfig)

	return r
}
