package main

import (
	"flag"
	"photobox-meta/cmd"
	"photobox-meta/config"
)

func main() {
	env := config.ConfigEnv()

	flag.StringVar(&env.GrpcPort, "port", env.GrpcPort, "grpc port")
	flag.Parse()

	cmd.Run(env)
}
