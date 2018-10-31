package clientfactory

import "github.com/aws/aws-sdk-go/service/kms"

type KmsClient interface {
	ListKeys(input *kms.ListKeysInput) (*kms.ListKeysOutput, error)
	ListAliases(input *kms.ListAliasesInput) (*kms.ListAliasesOutput, error)
}

type AWSKmsClient struct {
	api *kms.KMS
}

func (client AWSKmsClient) ListKeys(input *kms.ListKeysInput) (*kms.ListKeysOutput, error) {
	return client.api.ListKeys(input)
}
func (client AWSKmsClient) ListAliases(input *kms.ListAliasesInput) (*kms.ListAliasesOutput, error) {
	return client.api.ListAliases(input)
}
