package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Credentials struct {
	PostgreSQL struct {
		Host     string `yaml:"Host"`
		Port     string `yaml:"Port"`
		Password string `yaml:"Password"`
		Name     string `yaml:"Name"`
		User     string `yaml:"User"`
	}
}

func (c *Credentials) GetPostgreDBConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		c.PostgreSQL.User, c.PostgreSQL.Password, c.PostgreSQL.Host, c.PostgreSQL.Port, c.PostgreSQL.Name)
}

func NewCredentials() *Credentials {
	var creds Credentials

	viper.AddConfigPath(".")
	viper.SetConfigName("credentials")

	err := viper.ReadInConfig()
	if err != nil {
		logrus.WithError(err).Fatal("Read credentials error")
	}

	err = viper.UnmarshalKey("PostgreSQL", &creds.PostgreSQL)
	if err != nil {
		logrus.Fatalf("Read config PostgreSQL error err_003, %v", err)
	}

	return &creds
}
