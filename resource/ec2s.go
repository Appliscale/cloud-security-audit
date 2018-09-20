package resource

import (
	"github.com/Appliscale/tyr/configuration"
	"github.com/Appliscale/tyr/tyrsession"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Ec2s []*ec2.Instance

func (e *Ec2s) LoadFromAWS(config *configuration.Config, region string) error {
	ec2API, err := config.ClientFactory.GetEc2Client(tyrsession.SessionConfig{Profile: config.Profile, Region: region})
	if err != nil {
		return err
	}

	q := &ec2.DescribeInstancesInput{}
	for {
		result, err := ec2API.DescribeInstances(q)
		if err != nil {
			return err
		}
		for _, reservation := range result.Reservations {
			*e = append(*e, reservation.Instances...)
		}
		if result.NextToken == nil {
			break
		}
		q.NextToken = result.NextToken
	}

	return nil
}
