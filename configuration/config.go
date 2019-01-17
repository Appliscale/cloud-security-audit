package configuration

import (
	"github.com/Appliscale/cloud-security-audit/csasession/clientfactory"
	"github.com/Appliscale/cloud-security-audit/csasession/sessionfactory"
	"github.com/Appliscale/cloud-security-audit/logger"
)

type Config struct {
	Regions        *[]string
	Services       *[]string
	Profile        string
	SessionFactory *sessionfactory.SessionFactory
	ClientFactory  clientfactory.ClientFactory
	Logger         *logger.Logger
	Mfa            bool
	MfaDuration    int64
	Format         int
	OutputFile     string
}

func GetConfig() (config Config) {
	myLogger := logger.CreateDefaultLogger()
	config.Logger = &myLogger
	config.SessionFactory = sessionfactory.New()
	config.ClientFactory = clientfactory.New(config.SessionFactory)

	return config
}
