package csasession

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// SessionConfig provides session configuration for clients.
type SessionConfig struct {
	Profile string
	Region  string
}

// CreateSession returns new AWS session.
func CreateSession(config SessionConfig) (*session.Session, error) {

	sess, err := session.NewSessionWithOptions(
		session.Options{
			Config: aws.Config{
				Region: &config.Region,
			},
			Profile: config.Profile,
		})

	return sess, err
}

// GetAvailableRegions returns list of available regions.
func GetAvailableRegions() *[]string {
	return &[]string{
		"us-east-2",
		"us-east-1",
		"us-west-1",
		"us-west-2",
		"ap-northeast-1",
		"ap-northeast-2",
		"ap-south-1",
		"ap-southeast-1",
		"ap-southeast-2",
		"ca-central-1",
		"eu-central-1",
		"eu-west-1",
		"eu-west-2",
		"eu-west-3",
		"sa-east-1",
	}
}
