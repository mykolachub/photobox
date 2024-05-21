package cmd

import (
	"fmt"
	"log"
	"net"
	"photobox-meta/config"
	"photobox-meta/internal/services"
	"photobox-meta/internal/storage/postgres"
	"photobox-meta/proto"

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

	storages := services.Storages{
		MetaRepo: postgres.InitMetaRepo(db),
	}

	metaService := services.NewMetaService(storages.MetaRepo, services.MetaServiceConfig{})

	// gRPC Server
	grpcServer(env.GrpcPort, *metaService)
}

func grpcServer(port string, metaService services.MetaService) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	handleErr(err, ErrTypeGrpcTcpListen)

	s := grpc.NewServer()
	proto.RegisterMetaServiceServer(s, &metaService)
	reflection.Register(s)

	log.Printf("gRPC User server listening at %v\n", lis.Addr())
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
