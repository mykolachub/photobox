package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	GrpcPort           string `envconfig:"GRPC_PORT"`
	PostgresDBUser     string `envconfig:"POSTGRES_DBUSER"`
	PostgresDBPassword string `envconfig:"POSTGRES_DBPASSWORD"`
	PostgresDBName     string `envconfig:"POSTGRES_DBNAME"`
	PostgresDBPort     string `envconfig:"POSTGRES_DBPORT"`
	PostgresDBHost     string `envconfig:"POSTGRES_DBHOST"`
	PostgresDBSSLMode  string `envconfig:"POSTGRES_DBSSLMODE"`
}

func ConfigEnv() *Env {
	var env Env
	err := envconfig.Process("PB_META", &env)
	if err != nil {
		log.Fatal(err.Error())
	}

	return &env
}
