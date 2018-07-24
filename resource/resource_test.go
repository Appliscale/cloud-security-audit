package resource

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/stretchr/testify/assert"
)

type MockResource struct {
	DummyProperty string
}

func (m *MockResource) LoadFromAWS(sess *session.Session) error {
	return errors.New("Dummy error")
}

func TestResourcesReturnsErrorFromLoadFromAWSFunc(t *testing.T) {
	sess := &session.Session{}
	mockResource := &MockResource{}
	assert.Error(t, LoadResources(sess, mockResource))
}
