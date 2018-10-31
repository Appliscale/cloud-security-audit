package clientfactory

import "github.com/aws/aws-sdk-go/service/s3"

type S3Client interface {
	GetBucketPolicy(input *s3.GetBucketPolicyInput) (*s3.GetBucketPolicyOutput, error)
	GetBucketEncryption(input *s3.GetBucketEncryptionInput) (*s3.GetBucketEncryptionOutput, error)
	GetBucketLogging(input *s3.GetBucketLoggingInput) (*s3.GetBucketLoggingOutput, error)
	GetBucketAcl(input *s3.GetBucketAclInput) (*s3.GetBucketAclOutput, error)
	ListBuckets(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error)
	GetBucketLocation(input *s3.GetBucketLocationInput) (*s3.GetBucketLocationOutput, error)
}

type AWSS3Client struct {
	api *s3.S3
}

func (client AWSS3Client) GetBucketPolicy(input *s3.GetBucketPolicyInput) (*s3.GetBucketPolicyOutput, error) {
	return client.api.GetBucketPolicy(input)
}
func (client AWSS3Client) GetBucketEncryption(input *s3.GetBucketEncryptionInput) (*s3.GetBucketEncryptionOutput, error) {
	return client.api.GetBucketEncryption(input)
}
func (client AWSS3Client) GetBucketLogging(input *s3.GetBucketLoggingInput) (*s3.GetBucketLoggingOutput, error) {
	return client.api.GetBucketLogging(input)
}
func (client AWSS3Client) GetBucketAcl(input *s3.GetBucketAclInput) (*s3.GetBucketAclOutput, error) {
	return client.api.GetBucketAcl(input)
}
func (client AWSS3Client) ListBuckets(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
	return client.api.ListBuckets(input)
}

func (client AWSS3Client) GetBucketLocation(input *s3.GetBucketLocationInput) (*s3.GetBucketLocationOutput, error) {
	return client.api.GetBucketLocation(input)
}
