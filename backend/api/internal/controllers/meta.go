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
	pb "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MetaHandler struct {
	metaClient proto.MetaServiceClient
	middle     middlewares.Middleware
	mq         MQ
}

func InitMetaHandler(r *gin.Engine, metaClient proto.MetaServiceClient, middle middlewares.Middleware, mq MQ) {
	handler := MetaHandler{metaClient: metaClient, middle: middle, mq: mq}

	meta := r.Group("/api/meta")
	{
		meta.POST("", handler.middle.Protect(), handler.uploadMeta)
		meta.GET("", handler.middle.Protect(), handler.getMeta)
		meta.GET("/files", handler.middle.Protect(), handler.getFile)
		meta.PATCH("/:id", handler.updateMeta)
		meta.DELETE("/:id", handler.deleteMeta)
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

	fileWidth, err := strconv.ParseInt(c.PostForm("fileWidth"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileHeight, err := strconv.ParseInt(c.PostForm("fileHeight"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Preparing body
	protoReqBody := &proto.UplodaMetaRequest{
		UserId:           userId,
		Filename:         fileHeader.Filename,
		File:             fileBytes,
		FileLastModified: timestamppb.New(fileLastModified),
		FileWidth:        int32(fileWidth),
		FileHeight:       int32(fileHeight),
	}

	bytesReqBody, err := pb.Marshal(protoReqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal proto request body"})
		return
	}

	err = h.mq.Publish(c, "meta_upload", bytesReqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send image to queue"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "image is uploading..."})
}

func (h MetaHandler) getMeta(c *gin.Context) {
	userId := c.Keys[config.PayloadUserId].(string)

	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user unathorized"})
		return
	}

	res, err := h.metaClient.GetMetaByUser(c, &proto.GetMetaByUserRequest{UserId: userId})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h MetaHandler) getFile(c *gin.Context) {
	userId := c.Keys[config.PayloadUserId].(string)

	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user unathorized"})
		return
	}

	key := c.Query("file_location")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing file key parameter"})
		return
	}

	res, err := h.metaClient.GetFileByKey(c, &proto.GetFileByKeyRequest{UserId: userId, Key: key})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contentType := http.DetectContentType(res.File)
	c.Data(http.StatusOK, contentType, res.File)
}

func (h MetaHandler) deleteMeta(c *gin.Context) {
	id := c.Param("id")
	if len(id) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id parameter"})
		return
	}

	res, err := h.metaClient.DeleteMetaById(c, &proto.DeleteMetaByIdRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h MetaHandler) updateMeta(c *gin.Context) {

}
