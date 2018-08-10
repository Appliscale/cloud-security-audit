package configuration

import (
	"github.com/Appliscale/tyr/logger"
)

type Config struct {
	Region  string
	Service string
	Profile string
	Logger  *logger.Logger
}

func GetConfig() (config Config) {
	myLogger := logger.CreateDefaultLogger()
	config = Config{
		Logger: &myLogger,
	}
	return config
}
