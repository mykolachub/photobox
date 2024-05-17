package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	HttpPort            string `envconfig:"HTTP_PORT"`
	JWTSecret           string `envconfig:"JWT_SECRET"`
	GrpcPort            string `envconfig:"GRPC_PORT"`
	GrpcUserServicePort string `envconfig:"GRPC_USER_SERVICE_PORT"`
	GrpcUserServiceHost string `envconfig:"GRPC_USER_SERVICE_HOST"`
	GoogleClientId      string `envconfig:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret  string `envconfig:"GOOGLE_CLIENT_SECRET"`
	PostgresDBUser      string `envconfig:"POSTGRES_DBUSER"`
	PostgresDBPassword  string `envconfig:"POSTGRES_DBPASSWORD"`
	PostgresDBName      string `envconfig:"POSTGRES_DBNAME"`
	PostgresDBPort      string `envconfig:"POSTGRES_DBPORT"`
	PostgresDBHost      string `envconfig:"POSTGRES_DBHOST"`
	PostgresDBSSLMode   string `envconfig:"POSTGRES_DBSSLMODE"`
}

func ConfigEnv() *Env {
	var env Env
	err := envconfig.Process("PB_AUTH", &env)
	if err != nil {
		log.Fatal(err.Error())
	}

	return &env
}
