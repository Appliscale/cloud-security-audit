package resource

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type SecurityGroups map[string][]*ec2.IpPermission

func (s *SecurityGroups) LoadFromAWS(sess *session.Session) error {

	ec2API := ec2.New(sess)

	q := &ec2.DescribeSecurityGroupsInput{}

	for {
		result, err := ec2API.DescribeSecurityGroups(q)
		if err != nil {
			return err
		}
		for _, sg := range result.SecurityGroups {
			(*s)[*sg.GroupId] = append((*s)[*sg.GroupId], sg.IpPermissions...)
		}
		if result.NextToken == nil {
			break
		}
		q.NextToken = result.NextToken
	}
	return nil
}

func (s *SecurityGroups) GetIpPermissionsByID(groupID string) []*ec2.IpPermission {

	if sg, ok := (*s)[groupID]; ok {
		return sg
	}
	return nil
}
