package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config/")
	viper.AutomaticEnv()

	// 允许通过命令行参数覆盖
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.name", "DB_NAME")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Using config from environment or flags (no config file found)")
	}

	return &Config{
		DBUser:     viper.GetString("database.user"),
		DBPassword: viper.GetString("database.password"),
		DBHost:     viper.GetString("database.host"),
		DBPort:     viper.GetString("database.port"),
		DBName:     viper.GetString("database.name"),
	}

}
