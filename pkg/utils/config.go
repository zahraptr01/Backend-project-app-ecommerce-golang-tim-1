package utils

import (
	"github.com/spf13/viper"
)

type Configuration struct {
	AppName             string
	Port                string
	Debug               bool
	DB                  DatabaseConfig
	Limit               int
	PathLogger          string
	PathUpload          string
	RedisAddress        string
	MailersendApiKey    string
	MailersendFromEmail string
	SMTPHost            string
	SMTPPort            int
	SMTPEmail           string
	SMTPPassword        string
}

type DatabaseConfig struct {
	Name         string
	Username     string
	Password     string
	Host         string
	Port         string
	TimeZone     string
	Logging      bool
	MaxIdleConns int
	MaxOpenConns int
	MaxIdleTime  int
	MaxLifeTime  int
}

func ReadConfiguration() (Configuration, error) {
	viper.SetConfigFile(".env") // read file .env
	viper.SetConfigType("env")
	err := viper.ReadInConfig()
	if err != nil {
		return Configuration{}, err
	}
	viper.AutomaticEnv() // read env os

	return Configuration{
		AppName:             viper.GetString("APP_NAME"),
		Port:                viper.GetString("PORT"),
		Debug:               viper.GetBool("DEBUG"),
		Limit:               viper.GetInt("LIMIT"),
		PathLogger:          viper.GetString("PATH_LOGGER"),
		PathUpload:          viper.GetString("PATH_UPLOAD"),
		RedisAddress:        viper.GetString("REDIS_ADDRESS"),
		MailersendApiKey:    viper.GetString("MAILERSEND_API_KEY"),
		MailersendFromEmail: viper.GetString("MAILERSEND_FROM_EMAIL"),
		DB: DatabaseConfig{
			Name:         viper.GetString("DATABASE_NAME"),
			Username:     viper.GetString("DATABASE_USER"),
			Password:     viper.GetString("DATABASE_PASSWORD"),
			Host:         viper.GetString("DATABASE_HOST"),
			Port:         viper.GetString("DATABASE_PORT"),
			TimeZone:     viper.GetString("DATABASE_TIME_ZONE"),
			MaxIdleConns: viper.GetInt("DB_MAX_IDLE_CONNS"),
			MaxOpenConns: viper.GetInt("DB_MAX_OPEN_CONNS"),
			MaxIdleTime:  viper.GetInt("DB_MAX_IDLE_TIME"),
			MaxLifeTime:  viper.GetInt("DB_MAX_LIFE_TIME"),
		},
		SMTPHost:     viper.GetString("SMTPHost"),
		SMTPPort:     viper.GetInt("SMTPPort"),
		SMTPEmail:    viper.GetString("SMTPEmail"),
		SMTPPassword: viper.GetString("SMTPPassword"),
	}, nil
}
