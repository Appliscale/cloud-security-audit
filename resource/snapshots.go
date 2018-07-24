package resource

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Snapshots []*ec2.Snapshot

func (s *Snapshots) LoadFromAWS(sess *session.Session) error {
	ec2API := ec2.New(sess)
	q := &ec2.DescribeSnapshotsInput{
		OwnerIds: []*string{aws.String("self")},
	}
	for {
		result, err := ec2API.DescribeSnapshots(q)
		if err != nil {
			return err
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
