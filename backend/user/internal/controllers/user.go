package controllers

import (
	"context"
	"fmt"
	"net/http"
	"photobox-user/config"
	"photobox-user/internal/middlewares"
	"photobox-user/proto"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	usrSvc UserService
	cfg    UserHandlerConfig
}

type UserHandlerConfig struct {
	JwtSecret string
}

type UserService interface {
	DeleteUser(context.Context, *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error)
	GetAllUsers(context.Context, *proto.GetAllUsersRequest) (*proto.GetAllUsersResponse, error)
	GetMe(context.Context, *proto.GetMeRequest) (*proto.UserResponse, error)
	GetUser(context.Context, *proto.GetUserRequest) (*proto.UserResponse, error)
	Login(context.Context, *proto.LoginRequest) (*proto.LoginResponse, error)
	Signup(context.Context, *proto.SignupRequest) (*proto.SignupResponse, error)
	UpdateUser(context.Context, *proto.UpdateUserRequest) (*proto.UserResponse, error)
}

func InitUserHandler(r *gin.Engine, usrSvc UserService, cfg UserHandlerConfig) {
	handler := UserHandler{usrSvc: usrSvc, cfg: cfg}

	middle := middlewares.InitMiddleware(middlewares.MiddlewareConfig{
		JwtSecret: handler.cfg.JwtSecret,
	})

	r.POST("/api/login", handler.login)
	r.POST("/api/signup", handler.signup)

	users := r.Group("/api/users", middle.Protect())
	{
		users.GET("", handler.getAllUsers)
		users.GET("/me", handler.me)
		users.GET("/:user_id", handler.getUser)
		users.PATCH("/:user_id", handler.updateUser)
		users.DELETE("/:user_id", handler.deleteUser)
	}
}

func (h UserHandler) login(c *gin.Context) {
	var req proto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("login invalid data: %s", err)})
		return
	}

	res, err := h.usrSvc.Login(c, &req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": fmt.Sprintf("login: %s", err.Error())})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": res})
}

func (h UserHandler) signup(c *gin.Context) {
	var req proto.SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("signup invalid data: %s", err)})
		return
	}

	res, err := h.usrSvc.Signup(c, &req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": fmt.Sprintf("signup: %s", err.Error())})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": res})
}

func (h UserHandler) me(c *gin.Context) {
	userId := c.Keys[config.PayloadUserId].(string)

	res, err := h.usrSvc.GetUser(c, &proto.GetUserRequest{UserId: userId})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h UserHandler) getUser(c *gin.Context) {
	user_id := c.Param("user_id")
	if len(user_id) == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "user id parameter missing"})
		return
	}

	res, err := h.usrSvc.GetUser(c, &proto.GetUserRequest{UserId: user_id})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h UserHandler) getAllUsers(c *gin.Context) {
	res, err := h.usrSvc.GetAllUsers(c, &proto.GetAllUsersRequest{})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h UserHandler) updateUser(c *gin.Context) {
	user_id := c.Param("user_id")
	if len(user_id) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing user id parameter"})
		return
	}

	var req proto.UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("update invalid data: %s", err)})
		return
	}
	req.UserId = user_id

	res, err := h.usrSvc.UpdateUser(c, &req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})

}

func (h UserHandler) deleteUser(c *gin.Context) {
	user_id := c.Param("user_id")
	if len(user_id) == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "user id parameter missing"})
		return
	}

	res, err := h.usrSvc.DeleteUser(c, &proto.DeleteUserRequest{UserId: user_id})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}
