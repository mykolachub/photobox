package cmd

import (
	"fmt"
	"log"
	"net"
	"photobox-auth/config"
	"photobox-auth/internal/controllers"
	"photobox-auth/internal/services"
	"photobox-auth/proto"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func Run(env *config.Env) {
	// Google Oauth Configuration
	googleOauthConfig := &oauth2.Config{
		RedirectURL:  "http://localhost:8081/api/auth/login/google", // TODO: set callback to frontend
		ClientID:     env.GoogleClientId,
		ClientSecret: env.GoogleClientSecret,
		Scopes:       config.GoogleOauthScopes,
		Endpoint:     google.Endpoint,
	}

	// gRPC Clients
	userClient := grpcUserClient(env.GrpcUserServiceHost, env.GrpcUserServicePort)

	authSvc := services.NewAuthService(googleOauthConfig, services.AuthServiceConfig{JWTSecret: env.JWTSecret}, userClient)

	services := controllers.Services{
		AuthService: authSvc,
	}

	configs := controllers.Configs{}

	go grpcServer(env.GrpcPort, *authSvc)

	httpServer(env.HttpPort, services, configs)
}

func httpServer(port string, services controllers.Services, configs controllers.Configs) {
	router := controllers.InitRouter(gin.Default(), services, configs)
	router.Run(fmt.Sprintf(":%v", port))
}

func grpcServer(port string, authService services.AuthService) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	handleErr(err, ErrTypeGrpcTcpListen)

	s := grpc.NewServer()
	proto.RegisterAuthServiceServer(s, &authService)
	reflection.Register(s)

	log.Printf("gRPC server listening at %v\n", lis.Addr())
	handleErr(s.Serve(lis), ErrTypeGrpcServe)
}

func grpcUserClient(host, port string) proto.UserServiceClient {
	conn, err := grpc.NewClient(fmt.Sprintf("%v:%v", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	handleErr(err, ErrTypeGrpcUserDial)

	client := proto.NewUserServiceClient(conn)
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
