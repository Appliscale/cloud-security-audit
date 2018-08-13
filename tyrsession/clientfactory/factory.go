package clientfactory

import (
	"github.com/Appliscale/tyr/tyrsession/sessionfactory"
)

// ClientFactory provides methods for creation and management of service clients.
type ClientFactory struct {
	sessionFactory *sessionfactory.SessionFactory
}

// New creates a new instance of the ClientFactory.
func New(sessionFactory *sessionfactory.SessionFactory) *ClientFactory {
	return newFactory(sessionFactory)
}

func newFactory(sessionFactory *sessionfactory.SessionFactory) *ClientFactory {
	factory := &ClientFactory{
		sessionFactory: sessionFactory,
	}

	return factory
}
