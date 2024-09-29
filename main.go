package main

import (
	"fmt"
	"github/chetasp/crowdstrike/app"
	"github/chetasp/crowdstrike/config"
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
