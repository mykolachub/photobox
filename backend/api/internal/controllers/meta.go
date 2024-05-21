package controllers

import (
	"io"
	"net/http"
	"photobox-api/config"
	"photobox-api/internal/middlewares"
	"photobox-api/proto"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MetaHandler struct {
	metaClient proto.MetaServiceClient
	middle     middlewares.Middleware
}

func InitMetaHandler(r *gin.Engine, metaClient proto.MetaServiceClient, middle middlewares.Middleware) {
	handler := MetaHandler{metaClient: metaClient, middle: middle}

	meta := r.Group("/api/meta")
	{
		meta.POST("", handler.middle.Protect(), handler.uploadMeta)
		meta.GET("", handler.getMeta)
		meta.PATCH("/:id", handler.updateMeta)
		meta.DELETE("", handler.deleteMeta)
	}
}

func (h MetaHandler) uploadMeta(c *gin.Context) {
	userId := c.Keys[config.PayloadUserId].(string)

	// Process FormFile
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	// Getting []bytes
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Optional: Get file last_modified info
	createdAt, err := strconv.ParseInt(c.PostForm("lastModified"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fileLastModified := time.UnixMilli(createdAt)

	res, err := h.metaClient.UploadMeta(c, &proto.UplodaMetaRequest{
		UserId:           userId,
		Filename:         fileHeader.Filename,
		File:             fileBytes,
		FileLastModified: timestamppb.New(fileLastModified),
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h MetaHandler) getMeta(c *gin.Context) {
	var res any
	var err error

	id := c.Query("id")
	if id != "" {
		res, err = h.metaClient.GetMetaById(c, &proto.GetMetaByIdRequest{Id: id})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": res})
		return
	}

	userId := c.Query("user_id")
	if userId != "" {
		res, err = h.metaClient.GetMetaByUser(c, &proto.GetMetaByUserRequest{UserId: userId})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": res})
		return
	}

	res, err = h.metaClient.GetAllMeta(c, &proto.GetAllMetaRequest{})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h MetaHandler) deleteMeta(c *gin.Context) {

}

func (h MetaHandler) updateMeta(c *gin.Context) {

}
