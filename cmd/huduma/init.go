package main

import (
	"log"

	"github.com/huduma/internal/config"
	"github.com/huduma/internal/mongo"
	"github.com/spf13/cobra"
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

func init() {
	cobra.OnInitialize(initConf)
	rootCommand.PersistentFlags().StringP("config", "C", "", "An explicit config file to use")
	rootCommand.Flags().IntP("port", "p", 0, "the port to do things on")
	rootCommand.AddCommand(&versionCmd)

}
