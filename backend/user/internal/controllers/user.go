package controllers

import (
	"context"
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
	GetUser(context.Context, *proto.GetUserRequest) (*proto.GetUserResponse, error)
	Login(context.Context, *proto.LoginRequest) (*proto.LoginResponse, error)
	Signup(context.Context, *proto.SignupRequest) (*proto.SignupResponse, error)
	UpdateUser(context.Context, *proto.UpdateUserRequest) (*proto.GetUserResponse, error)
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
}

func (h UserHandler) signup(c *gin.Context) {
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
