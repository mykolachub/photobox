package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"photobox/image/config"
	processors "photobox/image/internal/processors/rekognition"
	"photobox/image/internal/services"
	"photobox/image/internal/storages/postgres"
	mq "photobox/image/pkg/mq/rabbitmq"
	"photobox/image/proto"

	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "google.golang.org/protobuf/proto"
)

func Run(env *config.Env) {
	// New RabbitMQ Connection
	conn, err := mq.NewRabbitMQConnection(mq.RabbitMQConfig{Host: env.RabbitMQHost, Port: env.RabbitMQPort, User: env.RabbitMQUser})
	handleErr(err, ErrRabbitMQConn)
	defer conn.Close()

	ch, err := conn.Channel()
	handleErr(err, ErrRabbitMQOpenChan)
	defer ch.Close()

	// New Postgres Connection
	db, err := postgres.NewPostgresConnection(postgres.PostgresConfig{
		DBUser:     env.PostgresDBUser,
		DBName:     env.PostgresDBName,
		DBPassword: env.PostgresDBPassword,
		DBSSLMode:  env.PostgresDBSSLMode,
		DBPort:     env.PostgresDBPort,
		DBHost:     env.PostgresDBHost,
	})
	if err != nil {
		log.Fatal(err)
	}

	// New AWS Rekognition Connection
	rekClient, err := processors.NewRekognitionConnection(processors.RekognitionConfig{
		Region:    env.RekognitionRegion,
		Secret:    env.RekognitionSecret,
		AccessKey: env.RekognitionAccessKey,
	})
	if err != nil {
		log.Fatal(err)
	}

	storages := services.Storages{
		LabesRepo: postgres.InitLabelsRepo(db),
	}

	processors := services.Processors{
		RekognitionRepo: processors.InitRekognition(rekClient, processors.RekognitionRepoConfig{
			BuckerName: env.S3BucketName,
		}),
	}

	imageService := services.InitImageService(processors.RekognitionRepo, storages.LabesRepo)

	// gRPC Server
	go grpcServer(env.GrpcPort, *imageService)

	// RabbitMQ Consuming
	rabbit := mq.InitRabbitMQ(ch)
	err = rabbit.Consume("image_detect", func(d amqp091.Delivery) {
		var req proto.DetectImageLabelsRequest
		err := pb.Unmarshal(d.Body, &req)
		log.Print(err)

		log.Printf("RabbitMQ [image_detect] Message: %+v", &req)

		_, err = imageService.DetectImageLabels(context.Background(), &req)
		if err != nil {
			log.Printf("RabbitMQ [image_detect] Message: %+v ERROR: %v", &req, err)
		}
	})
	handleErr(err, ErrRabbitMQConsume)
}

func grpcServer(port string, imageService services.ImageService) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	handleErr(err, ErrTypeGrpcTcpListen)

	s := grpc.NewServer()
	proto.RegisterImageServiceServer(s, &imageService)
	reflection.Register(s)

	log.Printf("gRPC User server listening at %v\n", lis.Addr())
	handleErr(s.Serve(lis), ErrTypeGrpcServe)
}

type AppErrType string

var (
	ErrTypeS3InitDB      AppErrType = "s3 init db"
	ErrTypeGrpcTcpListen AppErrType = "grpc tcp listen"
	ErrTypeGrpcServe     AppErrType = "grpc serve"
	ErrTypeGrpcUserDial  AppErrType = "grpc dial user client"
	ErrRabbitMQConn      AppErrType = "rabbitmq failed to connect to RabbitMQ"
	ErrRabbitMQOpenChan  AppErrType = "rabbitmq failed to open a channel"
	ErrRabbitMQConsume   AppErrType = "rabbitmq consume failed"
)

func handleErr(err error, format AppErrType) {
	if err != nil {
		log.Fatalf("[%s] %v", format, err)
	}
}
