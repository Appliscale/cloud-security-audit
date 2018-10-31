package clientfactory

import (
	"github.com/Appliscale/cloud-security-audit/csasession/sessionfactory"
)

// ClientFactory provides methods for creation and management of service clients.
type ClientFactoryAWS struct {
	sessionFactory *sessionfactory.SessionFactory
}

// New creates a new instance of the ClientFactory.
func New(sessionFactory *sessionfactory.SessionFactory) ClientFactory {
	factory := &ClientFactoryAWS{
		sessionFactory: sessionFactory,
	}

	return factory
}
