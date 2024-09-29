package main

import (
	"fmt"
	"github/chetasp/kafka/app"
	"github/chetasp/kafka/config"
)

func main() {

	config, err := config.GetConfig()
	if err != nil {
		fmt.Print("error getting config")
	}

	app, err := app.InitializeApplication(config)
	if err != nil {
		fmt.Print("error initializing application")
	}

	app.Run(config)
}
