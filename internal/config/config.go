package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

//Config is the structur of configutation
type Config struct {
	Port    int           `envconfig:"PORT"`
	Timeout time.Duration `envconfig:"TIMEOUT"`
	LogConf LogConfig
	Host    string `envconfig:"HOST"`
}

//Load loads all commands settings
func Load(configFile string) (*Config, error) {

	if configFile == "" {
		viper.SetConfigFile(configFile)
	}
	viper.SetConfigType("json")

	viper.SetEnvPrefix("HUDUMA")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetConfigName("huduma")
	viper.AddConfigPath("/etc/huduma")
	viper.AddConfigPath("../..etc/huduma")
	//viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)

		if !ok {
			return nil, errors.Wrap(err, "error when reading config from files")
		}
		return nil, err
	}

	fmt.Printf("port: %d\n", viper.GetInt("port"))

	var configur Config
	if err := viper.Unmarshal(&configur); err != nil {
		return nil, errors.Wrap(err, "Unmarschalling config")
	}
	return &configur, nil
}

//LoadEnv loads configuration from environment variables
func LoadEnv() (*Config, error) {

	var conf Config

	//Default loading
	err := viper.Unmarshal(&conf)
	if err != nil {
		return nil, err
	}

	err = envconfig.Process("", &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil

}
