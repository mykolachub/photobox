package controllers

import (
	"context"
	"fmt"
	"net/http"
	"photobox-user/proto"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	usrSvc UserService
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

func InitUserHandler(r *gin.Engine, usrSvc UserService) {
	handler := UserHandler{usrSvc: usrSvc}

	// middleware

	r.POST("/api/login", handler.login)
	r.POST("/api/signup", handler.signup)

	users := r.Group("/api/users", handler.getAllUsers)
	{
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
}

func (h UserHandler) getUser(c *gin.Context) {

}

func (h UserHandler) getAllUsers(c *gin.Context) {

}

func (h UserHandler) updateUser(c *gin.Context) {

}

func (h UserHandler) deleteUser(c *gin.Context) {

}
