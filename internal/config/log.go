package config

import (
	"bufio"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

//LogConfig is the structur of all logging
type LogConfig struct {
	Level string `json:"level"`
	File  string `json:"file"`
}

//Log configures all logging sets
func Log(conf *LogConfig) (*logrus.Entry, error) {
	host, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	if conf.File != "" {
		file, errOpen := os.OpenFile(conf.File, os.O_RDWR|os.O_APPEND, 0660)
		if errOpen != nil {
			return nil, errOpen
		}
		logrus.SetOutput(bufio.NewWriter(file))
	}
	level, err := logrus.ParseLevel(strings.ToUpper(conf.Level))
	if err != nil {
		return nil, err
	}
	logrus.SetLevel(level)

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		DisableTimestamp: false,
	})
	return logrus.StandardLogger().WithField("host", host), nil
}
