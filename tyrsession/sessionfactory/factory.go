package sessionfactory

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws/session"
)

// SessionFactory provides methods for creation and management of service clients and sessions.
type SessionFactory struct {
	regionToSession map[string]*session.Session
	mutex           sync.Mutex
}

// New creates a new instance of the SessionFactory.
func New() *SessionFactory {
	return newFactory()
}

func newFactory() *SessionFactory {
	factory := &SessionFactory{
		regionToSession: make(map[string]*session.Session),
	}

	return factory
}
