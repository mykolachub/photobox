package controllers

import (
	"net/http"
	"photobox-api/proto"

	"github.com/gin-gonic/gin"
	pb "google.golang.org/protobuf/proto"
)

type ImageHandler struct {
	imageClient proto.ImageServiceClient
	mq          MQ
}

func InitImageHandler(r *gin.Engine, imageClient proto.ImageServiceClient, mq MQ) {
	handler := ImageHandler{imageClient: imageClient, mq: mq}

	image := r.Group("/api/image")
	{
		image.POST("/labels", handler.detectLabels)
	}
}

func (h ImageHandler) detectLabels(c *gin.Context) {
	fileLocation := c.Query("file_location")
	metaId := c.Query("meta_id")

	protoReqBody := &proto.DetectImageLabelsRequest{
		FileLocation: fileLocation,
		MetaId:       metaId,
	}

	bytesReqBody, err := pb.Marshal(protoReqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal proto request body"})
		return
	}

	err = h.mq.Publish(c, "image_detect", bytesReqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send image to queue"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "image being processed"})
}
