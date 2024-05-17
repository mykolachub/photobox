package cmd

import (
	"fmt"
	"log"
	"photobox-api/config"
	"photobox-api/internal/controllers"
	"photobox-api/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Run(env *config.Env) {
	userClient := grpcUserClient(env.GrpcUserServiceHost, env.GrpcUserServicePort)
	authClient := grpcAuthClient(env.GrpcAuthServiceHost, env.GrpcAuthServicePort)

	services := controllers.Services{
		AuthClient: authClient,
		UserClient: userClient,
	}

	configs := controllers.Configs{
		UserConfig: controllers.UserHandlerConfig{
			JwtSecret: env.JWTSecret,
		},
	}

	httpServer(env.HttpPort, services, configs)
}

func httpServer(port string, services controllers.Services, configs controllers.Configs) {
	router := controllers.InitRouter(gin.Default(), services, configs)
	router.Run(fmt.Sprintf(":%v", port))
}

func grpcUserClient(host, port string) proto.UserServiceClient {
	conn, err := grpc.NewClient(fmt.Sprintf("%v:%v", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	handleErr(err, ErrTypeGrpcUserDial)

	client := proto.NewUserServiceClient(conn)
	return client
}

func grpcAuthClient(host, port string) proto.AuthServiceClient {
	conn, err := grpc.NewClient(fmt.Sprintf("%v:%v", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	handleErr(err, ErrTypeGrpcUserDial)

	client := proto.NewAuthServiceClient(conn)
	return client
}

type AppErrType string

var (
	ErrTypeGinRouterRun  AppErrType = "gin router run"
	ErrTypeGrpcTcpListen AppErrType = "grpc tcp listen"
	ErrTypeGrpcServe     AppErrType = "grpc serve"
	ErrTypeGrpcUserDial  AppErrType = "gprc user client dial"
)

func handleErr(err error, format AppErrType) {
	if err != nil {
		log.Fatalf("[%s] %v", format, err)
	}
}
