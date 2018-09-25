package resource

import (
	"github.com/Appliscale/tyr/configuration"
	"github.com/Appliscale/tyr/tyrsession"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Volumes []*ec2.Volume

func (v *Volumes) LoadFromAWS(config *configuration.Config, region string) error {
	ec2API, err := config.ClientFactory.GetEc2Client(tyrsession.SessionConfig{Profile: config.Profile, Region: region})
	if err != nil {
		return err
	}
	q := &ec2.DescribeVolumesInput{}
	for {
		result, err := ec2API.DescribeVolumes(q)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case "OptInRequired":
					break
				default:
					return err
				}
			} else {
				return err
			}
		}
		*v = append(*v, result.Volumes...)
		if result.NextToken == nil {
			break
		}
		q.NextToken = result.NextToken

	}

	return nil
}

func (v *Volumes) FindById(id string) *ec2.Volume {
	for _, volume := range *v {
		if *volume.VolumeId == id {
			return volume
		}
	}
	return nil
}
