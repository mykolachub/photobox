package main

import (
	"fmt"
	"photobox-user/cmd"
	"photobox-user/config"
)

func main() {
	env := config.ConfigEnv()
	fmt.Printf("env: %v\n", env)

	cmd.Run(env)
}
