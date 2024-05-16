package middlewares

import (
	"net/http"
	"photobox-user/config"
	"photobox-user/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Middleware struct {
	cfg MiddlewareConfig
}

type MiddlewareConfig struct {
	JwtSecret string
}

func InitMiddleware(cfg MiddlewareConfig) Middleware {
	return Middleware{cfg: cfg}
}

func (m Middleware) Protect() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(config.AuthorizationHeader)
		if len(authHeader) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header missing"})
			return
		}

		accessToken, err := utils.ValidateBearerHeader(authHeader)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		token, err := utils.ParseAndValidateToken(accessToken, m.cfg.JwtSecret)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		userId := token.Claims.(jwt.MapClaims)[config.JWTClainsUserId]
		userEmail := token.Claims.(jwt.MapClaims)[config.JWTClainsUserEmail]

		ctx.Set(config.PayloadUserId, userId)
		ctx.Set(config.PayloadUserEmail, userEmail)

		ctx.Next()
	}
}
