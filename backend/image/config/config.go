package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	GrpcPort             string `envconfig:"GRPC_PORT"`
	PostgresDBUser       string `envconfig:"POSTGRES_DBUSER"`
	PostgresDBPassword   string `envconfig:"POSTGRES_DBPASSWORD"`
	PostgresDBName       string `envconfig:"POSTGRES_DBNAME"`
	PostgresDBPort       string `envconfig:"POSTGRES_DBPORT"`
	PostgresDBHost       string `envconfig:"POSTGRES_DBHOST"`
	PostgresDBSSLMode    string `envconfig:"POSTGRES_DBSSLMODE"`
	RekognitionRegion    string `envconfig:"REKOGNIION_REGION"`
	RekognitionAccessKey string `envconfig:"REKOGNIION_SECRET_ACCESS_KEY"`
	RekognitionSecret    string `envconfig:"REKOGNIION_ACCESS_KEY"`
	S3BucketName         string `envconfig:"S3_BUCKET_NAME"`
	RabbitMQHost         string `envconfig:"RABBITMQ_HOST"`
	RabbitMQPort         string `envconfig:"RABBITMQ_PORT"`
	RabbitMQUser         string `envconfig:"RABBITMQ_USER"`
}

func ConfigEnv() *Env {
	var env Env
	err := envconfig.Process("PB_IMAGE", &env)
	if err != nil {
		log.Fatal(err.Error())
	}

	return &env
}
