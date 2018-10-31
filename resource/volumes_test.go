package resource

import (
	"github.com/Appliscale/cloud-security-audit/configuration"
	"github.com/Appliscale/cloud-security-audit/csasession"
	"github.com/Appliscale/cloud-security-audit/csasession/clientfactory/mocks"
	"github.com/aws/aws-sdk-go/service/ec2"
	"testing"
)

func TestLoadVolumesFromAWS(t *testing.T) {
	config := configuration.GetTestConfig(t)
	defer config.ClientFactory.(*mocks.ClientFactoryMock).Destroy()

	ec2Client, _ := config.ClientFactory.GetEc2Client(csasession.SessionConfig{})
	ec2Client.(*mocks.MockEC2Client).
		EXPECT().
		DescribeVolumes(&ec2.DescribeVolumesInput{}).
		Times(1).
		Return(&ec2.DescribeVolumesOutput{}, nil)

	LoadResource(&Volumes{}, &config, "region")
}
