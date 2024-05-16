package cmd

import (
	"fmt"
	"log"
	"net"
	"photobox-user/config"
	"photobox-user/internal/controllers"
	svcs "photobox-user/internal/services"
	"photobox-user/internal/storage/postgres"
	"photobox-user/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Run(env *config.Env) {
	// Database
	db, err := postgres.InitDBConnection(postgres.PostgresConfig{
		DBUser:     env.PostgresDBUser,
		DBName:     env.PostgresDBName,
		DBPassword: env.PostgresDBPassword,
		DBPort:     env.PostgresDBPort,
		DBSSLMode:  env.PostgresDBSSLMode,
		DBHost:     env.PostgresDBHost,
	})
	handleErr(err, ErrTypePostgresInitDB)

	storages := svcs.Storages{
		UserRepo: postgres.InitUserRepo(db),
	}

	userSvc := svcs.NewUserService(storages.UserRepo, svcs.UserServiceConfig{
		JwtSecret: env.JWTSecret,
	})
	services := controllers.Services{
		UserService: userSvc,
	}

	configs := controllers.Configs{
		UserHandlerConfig: controllers.UserHandlerConfig{
			JwtSecret: env.JWTSecret,
		},
	}

	// gRPC Server
	go grpcServer(env.GrpcPort, *userSvc)

	// HTTP Server
	httpServer(env.HttpPort, services, configs)
}

func httpServer(port string, services controllers.Services, configs controllers.Configs) {
	router := controllers.InitRouter(gin.Default(), services, configs)
	addr := fmt.Sprintf(":%v", port)
	handleErr(router.Run(addr), ErrTypeGinRouterRun)
}

func grpcServer(port string, userService svcs.UserService) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	handleErr(err, ErrTypeGrpcTcpListen)

	s := grpc.NewServer()
	proto.RegisterUserServiceServer(s, &userService)
	reflection.Register(s)

	log.Printf("gRPC server listening at %v\n", lis.Addr())
	handleErr(s.Serve(lis), ErrTypeGrpcServe)
}

type AppErrType string

var (
	ErrTypePostgresInitDB AppErrType = "postgres init db"
	ErrTypeGinRouterRun   AppErrType = "gin router run"
	ErrTypeGrpcTcpListen  AppErrType = "grpc tcp listen"
	ErrTypeGrpcServe      AppErrType = "grpc serve"
)

func handleErr(err error, format AppErrType) {
	if err != nil {
		log.Fatalf("[%s] %v", format, err)
	}
}
