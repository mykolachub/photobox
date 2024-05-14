package main

import (
	"photobox-user/cmd"
	"photobox-user/config"
)

func main() {
	env := config.ConfigEnv()

	cmd.Run(env)
}
