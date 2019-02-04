package clientfactory

import (
	"github.com/aws/aws-sdk-go/service/iam"
)

type IAMClient interface {
	ListUsers(input *iam.GetAccountAuthorizationDetailsInput) (*iam.GetAccountAuthorizationDetailsOutput, error)
	ListAccessKeys(input *iam.ListAccessKeysInput) (*iam.ListAccessKeysOutput, error)
}

type AWSIAMClient struct {
	api *iam.IAM
}

func (client AWSIAMClient) ListUsers(input *iam.GetAccountAuthorizationDetailsInput) (*iam.GetAccountAuthorizationDetailsOutput, error) {
	return client.api.GetAccountAuthorizationDetails(input)
}

func (client AWSIAMClient) ListAccessKeys(input *iam.ListAccessKeysInput) (*iam.ListAccessKeysOutput, error) {
	return client.api.ListAccessKeys(input)
}
