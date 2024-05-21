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
	S3Region           string `envconfig:"S3_REGION"`
	S3BucketName       string `envconfig:"S3_BUCKET_NAME"`
	S3Endpoint         string `envconfig:"S3_ENDPOINT"`
	S3AccessKey        string `envconfig:"S3_ACCESS_KEY"`
	S3SecretAccessKey  string `envconfig:"S3_SECRET_ACCESS_KEY"`
}

func ConfigEnv() *Env {
	var env Env
	err := envconfig.Process("PB_META", &env)
	if err != nil {
		log.Fatal(err.Error())
	}

	return &env
}
