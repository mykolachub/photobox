package main

import (
	"photobox-api/cmd"
	"photobox-api/config"
)

func main() {
	env := config.ConfigEnv()

	cmd.Run(env)
}
