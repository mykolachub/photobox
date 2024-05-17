package controllers

import (
	"context"
	"net/http"
	"photobox-auth/proto"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authSvc AuthService
	cfg     AuthHandlerConfig
}

type AuthHandlerConfig struct {
}

type AuthService interface {
	GoogleSignup(ctx context.Context, in *proto.GoogleSignupRequest) (*proto.GoogleSignupResponse, error)
	GoogleLogin(ctx context.Context, in *proto.GoogleLoginRequest) (*proto.GoogleLoginResponse, error)
}

func InitAuthRouter(r *gin.Engine, authSvc AuthService, cfg AuthHandlerConfig) {
	handler := AuthHandler{authSvc: authSvc, cfg: cfg}

	r.POST("/api/auth/signup/google", handler.signupGoogle)
	r.GET("/api/auth/login/google", handler.loginGoogle)

}

func (h AuthHandler) signupGoogle(c *gin.Context) {
	res, err := h.authSvc.GoogleSignup(c, &proto.GoogleSignupRequest{})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h AuthHandler) loginGoogle(c *gin.Context) {
	code := c.Query("code")

	res, err := h.authSvc.GoogleLogin(c, &proto.GoogleLoginRequest{Code: code})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}
