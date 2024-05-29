package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"photobox-meta/config"
	"photobox-meta/internal/services"
	"photobox-meta/internal/storage/postgres"
	"photobox-meta/internal/storage/s3"
	"photobox-meta/logger"
	mq "photobox-meta/pkg/mq/rabbitmq"
	"photobox-meta/proto"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	pb "google.golang.org/protobuf/proto"
)

func Run(env *config.Env) {
	// Logger
	l := logger.NewZap("")

	conn, err := mq.InitRabbitMQConnection(mq.RabbitMQConfig{Host: env.RabbitMQHost, Port: env.RabbitMQPort, User: env.RabbitMQUser})
	handleErr(err, ErrRabbitMQConn)
	defer conn.Close()

	ch, err := conn.Channel()
	handleErr(err, ErrRabbitMQOpenChan)
	defer ch.Close()

	mq := mq.InitRabbitMQ(ch)

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

	userClient := grpcUserClient(env.GrpcUserServiceHost, env.GrpcUserServicePort)

	storages := services.Storages{
		MetaRepo: postgres.InitMetaRepo(db),
		FileRepo: s3.InitFileRepo(s3Client, s3.FileRepoConfig{BucketName: env.S3BucketName}),
	}

	metaService := services.NewMetaService(storages.MetaRepo, storages.FileRepo, userClient, services.MetaServiceConfig{}, l, mq)

	// gRPC Server
	go grpcServer(env.GrpcPort, *metaService)

	err = mq.Consume("meta_upload", func(d amqp.Delivery) {
		var req proto.UplodaMetaRequest
		err := pb.Unmarshal(d.Body, &req)
		if err != nil {
			log.Printf("RabbitMQ [meta_upload] ERROR: %v", err)
		}
		log.Printf("RabbitMQ [meta_upload] Message: %+v", &req)

		_, err = metaService.UploadMeta(context.TODO(), &req)
		if err != nil {
			log.Printf("RabbitMQ [meta_upload] Message: %+v ERROR: %v", &req, err)
		}
	})
	handleErr(err, ErrRabbitMQConsume)
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

func grpcUserClient(host, port string) proto.UserServiceClient {
	conn, err := grpc.NewClient(fmt.Sprintf("%v:%v", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	handleErr(err, ErrTypeGrpcUserDial)

	client := proto.NewUserServiceClient(conn)
	return client
}

type AppErrType string

var (
	ErrTypePostgresInitDB AppErrType = "postgres init db"
	ErrTypeS3InitDB       AppErrType = "s3 init db"
	ErrTypeGinRouterRun   AppErrType = "gin router run"
	ErrTypeGrpcTcpListen  AppErrType = "grpc tcp listen"
	ErrTypeGrpcServe      AppErrType = "grpc serve"
	ErrTypeGrpcUserDial   AppErrType = "grpc dial user client"
	ErrRabbitMQConn       AppErrType = "rabbitmq failed to connect to RabbitMQ"
	ErrRabbitMQOpenChan   AppErrType = "rabbitmq failed to open a channel"
	ErrRabbitMQConsume    AppErrType = "rabbitmq consume failed"
)

func handleErr(err error, format AppErrType) {
	if err != nil {
		log.Fatalf("[%s] %v", format, err)
	}
}
