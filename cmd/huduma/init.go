package main

import (
	"log"
	"os"

	"github.com/huduma/api"
	"github.com/huduma/config"
)

var (
	globalConfig *config.Config
	configFile   string
)

func initAPIConfig() {
	logger, err := config.Log(&globalConfig.LogConf)

	if err != nil {
		log.Fatal("Unable to configure logging fo API" + err.Error())
	}

	logger.Info("API's starting on port: ", globalConfig.Port)
	apis := api.NewAPI(Version, globalConfig)

	err = apis.Serve()

	if err != nil {
		logger.WithError(err).Error("Error while listen to API:", err)
		os.Exit(1)
	}
	logger.Info("API shutdown")
}

func initConf() {
	globalConfig, err := config.Load(configFile)
	if err != nil {
		log.Fatal("Could not load config" + err.Error())
	}
	logger, err := config.Log(&globalConfig.LogConf)

	if err != nil {
		log.Fatal("failed to configur logging" + " " + err.Error())
	}
	logger.Info("starting config with", globalConfig)
}
