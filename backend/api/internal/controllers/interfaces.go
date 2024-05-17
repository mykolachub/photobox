package controllers

import "photobox-api/proto"

type Services struct {
	AuthClient proto.AuthServiceClient
	UserClient proto.UserServiceClient
}

type Configs struct {
	UserConfig UserHandlerConfig
}
