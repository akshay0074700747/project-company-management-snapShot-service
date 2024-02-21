package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	EndPoint        string
	AcessKey        string
	SecretAccessKey string
	MongoUrl        string
}

func LoadConfigurations() (Config, error) {

	if err := godotenv.Load(".env"); err != nil {
		return Config{}, err
	}

	var conf Config

	conf.EndPoint = os.Getenv("endpoint")
	conf.AcessKey = os.Getenv("accessKey")
	conf.SecretAccessKey = os.Getenv("secretKey")
	conf.MongoUrl = os.Getenv("mongourl")

	return conf, nil
}
