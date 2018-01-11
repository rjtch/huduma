package main

import (
	"log"

	"github.com/huduma/internal/config"
	"github.com/huduma/internal/mongo"
)

var (
	globalConfig *config.Config
	configFile   string
	session      *mongo.BooksDB
)

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

func initDB() {

}
