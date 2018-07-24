package resource

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Volumes []*ec2.Volume

func (v *Volumes) LoadFromAWS(sess *session.Session) error {
	ec2API := ec2.New(sess)
	q := &ec2.DescribeVolumesInput{}
	for {
		result, err := ec2API.DescribeVolumes(q)
		if err != nil {
			return err
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
