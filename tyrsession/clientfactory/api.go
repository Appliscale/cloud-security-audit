package clientfactory

import (
	"github.com/Appliscale/tyr/tyrsession"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/s3"
)

type AWSApi interface {
	GetKmsClient(config tyrsession.SessionConfig) (*KmsClient, error)
	GetEc2Client(config tyrsession.SessionConfig) (*EC2Client, error)
	GetS3Client(config tyrsession.SessionConfig) (*S3Client, error)
}

// GetKmsClient creates a new KMS client from cached session.
func (factory *ClientFactory) GetKmsClient(config tyrsession.SessionConfig) (KmsClient, error) {
	sess, err := factory.sessionFactory.GetSession(config)
	if err != nil {
		return nil, err
	}

	client := kms.New(sess)
	return AWSKmsClient{api: client}, nil
}

// GetEc2Client creates a new EC2 client from cached session.
func (factory *ClientFactory) GetEc2Client(config tyrsession.SessionConfig) (EC2Client, error) {
	sess, err := factory.sessionFactory.GetSession(config)
	if err != nil {
		return nil, err
	}

	client := ec2.New(sess)
	return AWSEC2Client{api: client}, nil
}

// GetS3Client creates a new S3 client from cached session.
func (factory *ClientFactory) GetS3Client(config tyrsession.SessionConfig) (S3Client, error) {
	sess, err := factory.sessionFactory.GetSession(config)
	if err != nil {
		return nil, err
	}

	client := s3.New(sess)
	return AWSS3Client{api: client}, nil
}
