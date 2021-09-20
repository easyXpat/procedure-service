package config

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

const (
	EnvPrefix = "XPAT_"
)

var DB *gorm.DB

type Configuration struct {
	ServerAddress              string
	DBHost                     string
	DBName                     string
	DBUser                     string
	DBPass                     string
	DBPort                     string
	DatabaseURL                     string
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
	c := &Configuration{
		ServerAddress: viper.GetString("SERVER_ADDRESS"),
		DBHost:                     viper.GetString("DB_HOST"),
		DBName:                     viper.GetString("DB_NAME"),
		DBUser:                     viper.GetString("DB_USER"),
		DBPass:                     viper.GetString("DB_PASS"),
		DBPort:                     viper.GetString("DB_PORT"),
	}
	c.SetPGConnectionString()
	return c
}

func(c *Configuration) SetPGConnectionString() {
	if viper.GetString("DATABASE_URL") != "" {
		c.DatabaseURL = viper.GetString("DATABASE_URL")
		return
	}
	c.DatabaseURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?", c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName)
}