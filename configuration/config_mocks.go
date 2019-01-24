package configuration

import (
	"github.com/Appliscale/cloud-security-audit/csasession/clientfactory/mocks"
	"github.com/Appliscale/cloud-security-audit/csasession/sessionfactory"
	"github.com/Appliscale/cloud-security-audit/logger"
	"testing"
)

func GetTestConfig(t *testing.T) (config Config) {
	myLogger := logger.CreateQuietLogger()
	config.Logger = &myLogger
	clientFactory := mocks.NewClientFactoryMock(t)
	config.ClientFactory = &clientFactory
	config.SessionFactory = sessionfactory.New()

	return config
}
