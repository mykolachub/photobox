package main

import (
	"photobox/image/cmd"
	"photobox/image/config"
)

func main() {
	env := config.ConfigEnv()

	cmd.Run(env)
}
