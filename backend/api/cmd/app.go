package cmd

import (
	"fmt"
	"log"
	"photobox-api/config"
	"photobox-api/internal/controllers"
	"photobox-api/internal/middlewares"
	mq "photobox-api/pkg/mq/rabbitmq"
	"photobox-api/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Run(env *config.Env) {
	// RabbitMQ Producer
	conn, err := mq.InitRabbitMQConnection(mq.RabbitMQConfig{Host: env.RabbitMQHost, Port: env.RabbitMQPort, User: env.RabbitMQUser})
	handleErr(err, ErrRabbitMQConn)
	defer conn.Close()

	ch, err := conn.Channel()
	handleErr(err, ErrRabbitMQOpenChan)
	defer ch.Close()

	mq := mq.InitRabbitMQ(ch)

	userClient := grpcUserClient(env.GrpcUserServiceHost, env.GrpcUserServicePort)
	authClient := grpcAuthClient(env.GrpcAuthServiceHost, env.GrpcAuthServicePort)
	metaClient := grpcMetaClient(env.GrpcMetaServiceHost, env.GrpcMetaServicePort)

	services := controllers.Services{
		AuthClient: authClient,
		UserClient: userClient,
		MetaClient: metaClient,
	}

	configs := controllers.Configs{
		UserConfig: controllers.UserHandlerConfig{
			JwtSecret: env.JWTSecret,
		},
	}

	middles := controllers.Middles{
		Middleware: middlewares.InitMiddleware(middlewares.MiddlewareConfig{
			JwtSecret: env.JWTSecret,
		}),
	}

	httpServer(env.HttpPort, services, configs, middles, mq)
}

func httpServer(port string, services controllers.Services, configs controllers.Configs, middles controllers.Middles, mq controllers.MQ) {
	router := controllers.InitRouter(gin.Default(), services, configs, middles, mq)
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
	handleErr(err, ErrTypeGrpcAuthDial)

	client := proto.NewAuthServiceClient(conn)
	return client
}

func grpcMetaClient(host, port string) proto.MetaServiceClient {
	conn, err := grpc.NewClient(fmt.Sprintf("%v:%v", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	handleErr(err, ErrTypeGrpcMetaDial)

	client := proto.NewMetaServiceClient(conn)
	return client
}

type AppErrType string

var (
	ErrTypeGinRouterRun  AppErrType = "gin router run"
	ErrTypeGrpcTcpListen AppErrType = "grpc tcp listen"
	ErrTypeGrpcServe     AppErrType = "grpc serve"
	ErrTypeGrpcUserDial  AppErrType = "gprc user client dial"
	ErrTypeGrpcMetaDial  AppErrType = "gprc meta client dial"
	ErrTypeGrpcAuthDial  AppErrType = "gprc auth client dial"
	ErrRabbitMQConn      AppErrType = "rabbitmq failed to connect to RabbitMQ"
	ErrRabbitMQOpenChan  AppErrType = "rabbitmq failed to open a channel"
)

func handleErr(err error, format AppErrType) {
	if err != nil {
		log.Fatalf("[%s] %v", format, err)
	}
}
