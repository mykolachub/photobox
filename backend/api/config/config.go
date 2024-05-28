package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	HttpPort            string `envconfig:"HTTP_PORT"`
	JWTSecret           string `envconfig:"JWT_SECRET"`
	GrpcAuthServicePort string `envconfig:"GRPC_AUTH_SERVICE_PORT"`
	GrpcAuthServiceHost string `envconfig:"GRPC_AUTH_SERVICE_HOST"`
	GrpcUserServicePort string `envconfig:"GRPC_USER_SERVICE_PORT"`
	GrpcUserServiceHost string `envconfig:"GRPC_USER_SERVICE_HOST"`
	GrpcMetaServicePort string `envconfig:"GRPC_META_SERVICE_PORT"`
	GrpcMetaServiceHost string `envconfig:"GRPC_META_SERVICE_HOST"`
	RabbitMQHost        string `envconfig:"RABBITMQ_HOST"`
	RabbitMQPort        string `envconfig:"RABBITMQ_PORT"`
	RabbitMQUser        string `envconfig:"RABBITMQ_USER"`
}

func ConfigEnv() *Env {
	var env Env
	err := envconfig.Process("PB_API", &env)
	if err != nil {
		log.Fatal(err.Error())
	}

	return &env
}
