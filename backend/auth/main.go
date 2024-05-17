package main

import (
	"photobox-auth/cmd"
	"photobox-auth/config"
)

func main() {
	env := config.ConfigEnv()

	cmd.Run(env)
}
