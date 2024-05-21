package controllers

import (
	"photobox-api/internal/middlewares"
	"photobox-api/proto"
)

type Services struct {
	AuthClient proto.AuthServiceClient
	UserClient proto.UserServiceClient
	MetaClient proto.MetaServiceClient
}

type Configs struct {
	UserConfig UserHandlerConfig
}

type Middles struct {
	Middleware middlewares.Middleware
}
