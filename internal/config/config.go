package config

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

var ErrConfig = errors.New("config error")
var logger log.Logger
var RuntimeViper = viper.New()
var consulAddr = flag.String("consulAddress", "localhost:8500", "Consul URL")
var consulKey = flag.String("consulKey", "reseller-nettv", "Consul Key")

type RabbitMQ struct {
	HostName     string  `json:"hostname" mapstructure:"hostname"`
	Port         string  `json:"port" mapstructure:"port"`
	Name         string  `json:"name" mapstructure:"name"`
	Username     string  `json:"username" mapstructure:"username"`
	Password     string  `json:"password" mapstructure:"password"`
	UseAuth      bool    `json:"useAuth" mapstructure:"useAuth"`
	ExchangeName string  `json:"exchangeName" mapstructure:"exchangeName"`
	QueueName    string  `json:"queueName" mapstructure:"queueName"`
	RoutingKey   *string `json:"routingKey" mapstructure:"routingKey"`
}

type Config struct {
	Debug struct {
		Enable bool `yaml:"debug"`
	}
	RabbitMQ RabbitMQ
}

func NewConfig() Config {

	RuntimeViper.AddRemoteProvider("consul", *consulAddr, *consulKey)
	RuntimeViper.SetConfigType("json")

	if err := RuntimeViper.ReadRemoteConfig(); err != nil {

		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			level.Error(logger).Log(err, err)
			os.Exit(1)
		} else {
			fmt.Println(err)
			level.Error(logger).Log(err, ErrConfig)
		}

	}

	var C Config

	err := RuntimeViper.Unmarshal(&C)

	if err != nil {

		level.Error(logger).Log("err", err)
	}
	return C
}
