package configuration

import (
	"github.com/Appliscale/tyr/logger"
	"github.com/Appliscale/tyr/tyrsession/clientfactory/mocks"
	"github.com/Appliscale/tyr/tyrsession/sessionfactory"
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
