package controllers

import (
	"context"
	"photobox-api/internal/middlewares"
	"photobox-api/proto"
)

type Services struct {
	AuthClient  proto.AuthServiceClient
	UserClient  proto.UserServiceClient
	MetaClient  proto.MetaServiceClient
	ImageClient proto.ImageServiceClient
}

type Configs struct {
	UserConfig UserHandlerConfig
}

type Middles struct {
	Middleware middlewares.Middleware
}

type MQ interface {
	Publish(ctx context.Context, name string, body []byte) error
}
