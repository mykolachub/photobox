package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	JWTSecret           string `envconfig:"JWT_SECRET"`
	GrpcPort            string `envconfig:"GRPC_PORT"`
	GrpcUserServicePort string `envconfig:"GRPC_USER_SERVICE_PORT"`
	GrpcUserServiceHost string `envconfig:"GRPC_USER_SERVICE_HOST"`
	GoogleClientId      string `envconfig:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret  string `envconfig:"GOOGLE_CLIENT_SECRET"`
	GoogleRedirectURL   string `envconfig:"GOOGLE_REDIRECT_URL"`
}

func ConfigEnv() *Env {
	var env Env
	err := envconfig.Process("PB_AUTH", &env)
	if err != nil {
		log.Fatal(err.Error())
	}

	return &env
}
