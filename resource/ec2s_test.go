package resource

import (
	"github.com/Appliscale/tyr/configuration"
	"github.com/Appliscale/tyr/tyrsession"
	"github.com/Appliscale/tyr/tyrsession/clientfactory/mocks"
	"github.com/aws/aws-sdk-go/service/ec2"
	"testing"
)

func TestLoadEC2sFromAWS(t *testing.T) {
	config := configuration.GetTestConfig(t)
	defer config.ClientFactory.(*mocks.ClientFactoryMock).Destroy()

	ec2Client, _ := config.ClientFactory.GetEc2Client(tyrsession.SessionConfig{})
	ec2Client.(*mocks.MockEC2Client).
		EXPECT().
		DescribeInstances(&ec2.DescribeInstancesInput{}).
		Times(1).
		Return(&ec2.DescribeInstancesOutput{
			Reservations: []*ec2.Reservation{},
			NextToken:    nil,
		}, nil)

	LoadResource(&Ec2s{}, &config, "region")
}
