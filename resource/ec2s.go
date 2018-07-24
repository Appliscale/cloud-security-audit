package resource

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Ec2s []*ec2.Instance

func (e *Ec2s) LoadFromAWS(sess *session.Session) error {
	ec2API := ec2.New(sess)
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
