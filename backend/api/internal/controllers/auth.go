package controllers

import (
	"net/http"
	"photobox-api/proto"

	"github.com/gin-gonic/gin"
)

type ApiHandler struct {
	authClient proto.AuthServiceClient
}

func InitAuthHandler(r *gin.Engine, authClient proto.AuthServiceClient) {
	handler := ApiHandler{authClient: authClient}

	auth := r.Group("/api/auth")
	{
		auth.POST("/signup/google", handler.signupGoogle)
		auth.POST("/login/google", handler.loginGoogle)
	}
}

func (h ApiHandler) signupGoogle(c *gin.Context) {
	res, err := h.authClient.GoogleSignup(c, &proto.GoogleSignupRequest{})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h ApiHandler) loginGoogle(c *gin.Context) {
	code := c.Query("code")

	res, err := h.authClient.GoogleLogin(c, &proto.GoogleLoginRequest{Code: code})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}
