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
	CreateUser(context.Context, *proto.CreateUserRequest) (*proto.UserResponse, error)
	GetUser(context.Context, *proto.GetUserRequest) (*proto.UserResponse, error)
	GetAllUsers(context.Context, *proto.GetAllUsersRequest) (*proto.GetAllUsersResponse, error)
	UpdateUser(context.Context, *proto.UpdateUserRequest) (*proto.UserResponse, error)
	DeleteUser(context.Context, *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error)
}

func InitUserHandler(r *gin.Engine, usrSvc UserService, cfg UserHandlerConfig) {
	handler := UserHandler{usrSvc: usrSvc, cfg: cfg}

	middle := middlewares.InitMiddleware(middlewares.MiddlewareConfig{
		JwtSecret: handler.cfg.JwtSecret,
	})

	r.POST("/api/users", handler.createUser)
	users := r.Group("/api/users", middle.Protect())
	{
		users.GET("", handler.getAllUsers)
		users.GET("/me", handler.me)
		users.GET("/:user_id", handler.getUser)
		users.PATCH("/:user_id", handler.updateUser)
		users.DELETE("/:user_id", handler.deleteUser)
	}
}

func (h UserHandler) createUser(c *gin.Context) {
	var req proto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("create invalid data: %s", err)})
		return
	}

	res, err := h.usrSvc.CreateUser(c, &req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": fmt.Sprintf("create: %s", err.Error())})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": res})
}

func (h UserHandler) me(c *gin.Context) {
	userId := c.Keys[config.PayloadUserId].(string)

	res, err := h.usrSvc.GetUser(c, &proto.GetUserRequest{Id: userId})
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

	res, err := h.usrSvc.GetUser(c, &proto.GetUserRequest{Id: user_id})
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
	req.Id = user_id

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

	res, err := h.usrSvc.DeleteUser(c, &proto.DeleteUserRequest{Id: user_id})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}
