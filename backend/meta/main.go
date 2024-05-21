package main

import (
	"photobox-meta/cmd"
	"photobox-meta/config"
)

func main() {
	env := config.ConfigEnv()

	cmd.Run(env)
}
