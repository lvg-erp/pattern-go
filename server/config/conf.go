package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Logger            logger
	WebServer         webSever
	BackGroundWorkers BackGroundWorkersStruct `yaml:"BackGroundWorkers"`
}

type logger struct {
	Level logrus.Level `yaml:"Level"`
	Type  string
}

type webSever struct {
	ListenPorts struct {
		Http  string `yaml:"Http"`
		Https string `yaml:"Https"`
	} `yaml:"ListenPorts"`
	HttpsSettings struct {
		Enabled    bool     `yaml:"Enabled"`
		AdminEmail string   `yaml:"AdminEmail"`
		Domains    []string `yaml:"Domains"`
	} `yaml:"HttpsSettings"`
}

type BackGroundWorkersStruct struct {
	ATrucksData struct {
		UpdateOnStart bool `yaml:"UpdateOnStart"`
		AutoUpdate    struct {
			Enable      bool `yaml:"Enable"`
			IntervalMin int  `yaml:"IntervalMin"`
		} `yaml:"AutoUpdate"`
	} `yaml:"ATrucksData"`
}

func NewConfig() *Config {
	var conf Config

	viper.AddConfigPath(".")
	viper.SetConfigName("conf")
	err := viper.ReadInConfig()

	if err != nil {
		logrus.Fatalf("Read config error, %v", err)
	}

	err = viper.UnmarshalKey("Logger", &conf.Logger)
	if err != nil {
		logrus.Fatalf("Read config in section Logger error, %v", err)
	}

	err = viper.UnmarshalKey("WebServer", &conf.WebServer)
	if err != nil {
		logrus.Fatalf("Read config in section WebServer error, %v", err)
	}

	err = viper.UnmarshalKey("BackGroundWorkers", &conf.BackGroundWorkers)
	if err != nil {
		logrus.Fatalf("Read config in section BackGroundWorkers error, %v", err)
	}

	return &conf
}
