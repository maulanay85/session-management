package main

import (
	"fmt"
	"scs-session/internal/config"
	"scs-session/internal/module"
)

func main() {
	conf, err := config.InitializeConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println("setup module ..")
	r := module.Init(*conf)

	err = r.Run(":" + conf.ServicePort)
	if err != nil {
		panic(err)
	}
}
