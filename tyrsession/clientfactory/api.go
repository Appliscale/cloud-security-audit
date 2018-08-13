package clientfactory

import (
	"github.com/Appliscale/tyr/tyrsession"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/s3"
)

// GetKmsClient creates a new KMS client from cached session.
func (factory *ClientFactory) GetKmsClient(config tyrsession.SessionConfig) (*kms.KMS, error) {
	sess, err := factory.sessionFactory.GetSession(config)
	if err != nil {
		return nil, err
	}

	client := kms.New(sess)
	return client, nil
}

// GetEc2Client creates a new EC2 client from cached session.
func (factory *ClientFactory) GetEc2Client(config tyrsession.SessionConfig) (*ec2.EC2, error) {
	sess, err := factory.sessionFactory.GetSession(config)
	if err != nil {
		return nil, err
	}

	client := ec2.New(sess)
	return client, nil
}

// GetS3Client creates a new S3 client from cached session.
func (factory *ClientFactory) GetS3Client(config tyrsession.SessionConfig) (*s3.S3, error) {
	sess, err := factory.sessionFactory.GetSession(config)
	if err != nil {
		return nil, err
	}

	client := s3.New(sess)
	return client, nil
}
