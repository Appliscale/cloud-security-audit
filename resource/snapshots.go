package resource

import (
	"github.com/Appliscale/tyr/configuration"
	"github.com/Appliscale/tyr/tyrsession"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Snapshots []*ec2.Snapshot

func (s *Snapshots) LoadFromAWS(config *configuration.Config, region string) error {
	ec2API, err := config.ClientFactory.GetEc2Client(tyrsession.SessionConfig{Profile: config.Profile, Region: region})
	if err != nil {
		return err
	}
	q := &ec2.DescribeSnapshotsInput{
		OwnerIds: []*string{aws.String("self")},
	}
	for {
		result, err := ec2API.DescribeSnapshots(q)
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
		*s = append(*s, result.Snapshots...)
		if result.NextToken == nil {
			break
		}
		q.NextToken = result.NextToken
	}

	return nil
}

func (s *Snapshots) FindById(id string) *ec2.Snapshot {
	for _, snapshot := range *s {
		if id == *snapshot.SnapshotId {
			return snapshot
		}
	}
	return nil
}
