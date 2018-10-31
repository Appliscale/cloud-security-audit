package resource

import (
	"testing"

	"errors"
	"github.com/Appliscale/cloud-security-audit/configuration"
	"github.com/stretchr/testify/assert"
)

type MockResource struct {
	DummyProperty string
}

func (m *MockResource) LoadFromAWS(config *configuration.Config, region string) error {
	return errors.New("Dummy error")
}

func TestResourcesReturnsErrorFromLoadFromAWSFunc(t *testing.T) {
	mockResource := &MockResource{}
	config := configuration.GetTestConfig(t)
	assert.Error(t, LoadResources(&config, "region", mockResource))
}
