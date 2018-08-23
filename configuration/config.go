package configuration

import (
	"github.com/Appliscale/tyr/logger"
	"github.com/Appliscale/tyr/tyrsession/clientfactory"
	"github.com/Appliscale/tyr/tyrsession/sessionfactory"
)

type Config struct {
	Regions        *[]string
	Services       *[]string
	Profile        string
	SessionFactory *sessionfactory.SessionFactory
	ClientFactory  *clientfactory.ClientFactory
	Logger         *logger.Logger
}

func GetConfig() (config Config) {
	myLogger := logger.CreateDefaultLogger()
	config.Logger = &myLogger
	config.SessionFactory = sessionfactory.New()
	config.ClientFactory = clientfactory.New(config.SessionFactory)

	return config
}
