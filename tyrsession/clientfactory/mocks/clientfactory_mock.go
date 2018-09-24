package mocks

import (
	"github.com/Appliscale/tyr/tyrsession"
	"github.com/Appliscale/tyr/tyrsession/clientfactory"
	"github.com/golang/mock/gomock"
	"testing"
)

type ClientFactoryMock struct {
	mockCtrl  *gomock.Controller
	kmsClient *MockKmsClient
	ec2Client *MockEC2Client
	s3Client  *MockS3Client
}

func NewClientFactoryMock(t *testing.T) ClientFactoryMock {
	clientMock := ClientFactoryMock{
		mockCtrl: gomock.NewController(t),
	}
	return clientMock
}

func (client *ClientFactoryMock) GetKmsClient(config tyrsession.SessionConfig) (clientfactory.KmsClient, error) {
	if client.kmsClient == nil {
		client.kmsClient = NewMockKmsClient(client.mockCtrl)
	}
	return client.kmsClient, nil
}
func (client *ClientFactoryMock) GetEc2Client(config tyrsession.SessionConfig) (clientfactory.EC2Client, error) {
	if client.ec2Client == nil {
		client.ec2Client = NewMockEC2Client(client.mockCtrl)
	}
	return client.ec2Client, nil
}
func (client *ClientFactoryMock) GetS3Client(config tyrsession.SessionConfig) (clientfactory.S3Client, error) {
	if client.s3Client == nil {
		client.s3Client = NewMockS3Client(client.mockCtrl)
	}
	return client.s3Client, nil
}

func (client *ClientFactoryMock) Destroy() {
	client.mockCtrl.Finish()
}
