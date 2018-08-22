package sessionfactory

import (
	"github.com/Appliscale/tyr/tyrsession"

	"github.com/aws/aws-sdk-go/aws/session"
)

// GetSession returns cached AWS session.
func (factory *SessionFactory) GetSession(config tyrsession.SessionConfig) (*session.Session, error) {
	factory.mutex.Lock()
	defer factory.mutex.Unlock()

	if sess, ok := factory.regionToSession[config.Region]; ok {
		return sess, nil
	}

	return factory.NewSession(config)
}

// NewSession creates a new session and caches it.
func (factory *SessionFactory) NewSession(config tyrsession.SessionConfig) (*session.Session, error) {
	sess, err := tyrsession.CreateSession(config)
	if err != nil {
		return nil, err
	}

	factory.regionToSession[config.Region] = sess
	return sess, nil
}
