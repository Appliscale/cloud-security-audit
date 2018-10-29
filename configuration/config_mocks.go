package configuration

import (
	"github.com/Appliscale/cloud-security-audit/logger"
	"github.com/Appliscale/cloud-security-audit/csasession/clientfactory/mocks"
	"github.com/Appliscale/cloud-security-audit/csasession/sessionfactory"
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
