package cmd

import (
	"fmt"
	"log"
	"net"
	"photobox-meta/config"
	"photobox-meta/internal/services"
	"photobox-meta/internal/storage/postgres"
	"photobox-meta/internal/storage/s3"
	"photobox-meta/logger"
	"photobox-meta/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Run(env *config.Env) {
	// Logger
	l := logger.NewZap("")

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

	s3Client, err := s3.InitS3Connection(s3.S3BucketConfig{
		Region:    env.S3Region,
		Name:      env.S3BucketName,
		AccessKey: env.S3AccessKey,
		Secret:    env.S3SecretAccessKey,
		Endpoint:  env.S3Endpoint,
	})
	handleErr(err, ErrTypeS3InitDB)

	storages := services.Storages{
		MetaRepo: postgres.InitMetaRepo(db),
		FileRepo: s3.InitFileRepo(s3Client, s3.FileRepoConfig{BucketName: env.S3BucketName}),
	}

	metaService := services.NewMetaService(storages.MetaRepo, storages.FileRepo, services.MetaServiceConfig{}, l)

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
	ErrTypeS3InitDB       AppErrType = "s3 init db"
	ErrTypeGinRouterRun   AppErrType = "gin router run"
	ErrTypeGrpcTcpListen  AppErrType = "grpc tcp listen"
	ErrTypeGrpcServe      AppErrType = "grpc serve"
)

func handleErr(err error, format AppErrType) {
	if err != nil {
		log.Fatalf("[%s] %v", format, err)
	}
}
