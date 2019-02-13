package configuration

import (
	"github.com/Appliscale/cloud-security-audit/csasession/clientfactory"
	"github.com/Appliscale/cloud-security-audit/csasession/sessionfactory"
	"github.com/Appliscale/cloud-security-audit/logger"
	"github.com/Appliscale/cloud-security-audit/report"
	"os"
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
	PrintFormat    func(report.Report, *os.File)
	OutputFile     *os.File
}

func GetConfig() (config Config) {
	myLogger := logger.CreateDefaultLogger()
	config.Logger = &myLogger
	config.SessionFactory = sessionfactory.New()
	config.ClientFactory = clientfactory.New(config.SessionFactory)

	return config
}
