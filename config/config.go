package config

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/hashicorp/go-hclog"
)

const (
	EnvPrefix = "XPAT_"
)

type Configuration struct {
	ServerAddress              string
	DBHost                     string
	DBName                     string
	DBUser                     string
	DBPass                     string
	DBPort                     string
}

func setDefaults(logger hclog.Logger) {
	viper.SetDefault("SERVER_ADDRESS", "0.0.0.0:9090")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_NAME", "postgres")
	viper.SetDefault("DB_USER", "username")
	viper.SetDefault("DB_PASS", "password")
	viper.SetDefault("DB_PORT", "5432")
}

func NewConfiguration(logger hclog.Logger) *Configuration {
	logger.Debug("Fetch default configuration")
	viper.AutomaticEnv()
	setDefaults(logger)
	return &Configuration{
		ServerAddress: viper.GetString("SERVER_ADDRESS"),
		DBHost:                     viper.GetString("DB_HOST"),
		DBName:                     viper.GetString("DB_NAME"),
		DBUser:                     viper.GetString("DB_USER"),
		DBPass:                     viper.GetString("DB_PASS"),
		DBPort:                     viper.GetString("DB_PORT"),
	}
}

func(c *Configuration) GetPGConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?", c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName)
}