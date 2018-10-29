package sessionfactory

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws/session"
)

// sessionFactory provides methods for creation and management of service clients and sessions.
type SessionFactory struct {
	regionToSession map[string]*session.Session
	mutex           sync.Mutex
}

// New creates a new instance of the sessionFactory.
func New() *SessionFactory {
	factory := &SessionFactory{
		regionToSession: make(map[string]*session.Session),
	}

	return factory
}
