package controllers

import (
	"fmt"
	"net/http"
	"photobox-api/config"
	"photobox-api/internal/middlewares"
	"photobox-api/proto"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userClient proto.UserServiceClient
	cfg        UserHandlerConfig
	middle     middlewares.Middleware
}

type UserHandlerConfig struct {
	JwtSecret string
}

func InitUserHandler(r *gin.Engine, userClient proto.UserServiceClient, cfg UserHandlerConfig, middle middlewares.Middleware) {
	handler := UserHandler{userClient: userClient, cfg: cfg, middle: middle}

	users := r.Group("/api/users", handler.middle.Protect())
	{
		users.POST("", handler.createUser)
		users.GET("", handler.getAllUsers)
		users.GET("/me", handler.getMe)
		users.GET("/:user_id", handler.getUser)
		users.PATCH("/:user_id", handler.updateUser)
		users.DELETE("/:user_id", handler.deleteUser)

		users.PATCH("/storage", handler.updateStorage)
	}
}

func (h UserHandler) createUser(c *gin.Context) {
	var req proto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("create user invalid data: %s", err)})
		return
	}

	res, err := h.userClient.CreateUser(c, &proto.CreateUserRequest{})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h UserHandler) getAllUsers(c *gin.Context) {
	res, err := h.userClient.GetAllUsers(c, &proto.GetAllUsersRequest{})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h UserHandler) getMe(c *gin.Context) {
	userId := c.Keys[config.PayloadUserId].(string)

	res, err := h.userClient.GetUser(c, &proto.GetUserRequest{Id: userId})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	res, err := h.userClient.GetUser(c, &proto.GetUserRequest{})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	res, err := h.userClient.UpdateUser(c, &req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h UserHandler) updateStorage(c *gin.Context) {
	user_id := c.Keys[config.PayloadUserId].(string)

	var req proto.UpdateStorageUsedRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("update invalid data: %s", err)})
		return
	}
	req.Id = user_id

	res, err := h.userClient.UpdateStorageUsed(c, &req)
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

	res, err := h.userClient.DeleteUser(c, &proto.DeleteUserRequest{Id: user_id})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}
