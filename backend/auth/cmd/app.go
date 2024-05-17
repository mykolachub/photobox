package cmd

import (
	"fmt"
	"log"
	"net"
	"photobox-auth/config"
	"photobox-auth/internal/services"
	"photobox-auth/proto"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func Run(env *config.Env) {
	// Google Oauth Configuration
	googleOauthConfig := &oauth2.Config{
		RedirectURL:  env.GoogleRedirectURL, // TODO: set callback to frontend
		ClientID:     env.GoogleClientId,
		ClientSecret: env.GoogleClientSecret,
		Scopes:       config.GoogleOauthScopes,
		Endpoint:     google.Endpoint,
	}

	// gRPC Clients
	userClient := grpcUserClient(env.GrpcUserServiceHost, env.GrpcUserServicePort)

	authSvc := services.NewAuthService(googleOauthConfig, services.AuthServiceConfig{JWTSecret: env.JWTSecret}, userClient)

	grpcServer(env.GrpcPort, *authSvc)
}

func grpcServer(port string, authService services.AuthService) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	handleErr(err, ErrTypeGrpcTcpListen)

	s := grpc.NewServer()
	proto.RegisterAuthServiceServer(s, &authService)
	reflection.Register(s)

	log.Printf("gRPC Auth server listening at %v\n", lis.Addr())
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
