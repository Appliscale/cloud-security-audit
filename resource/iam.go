package resource

import (
	"github.com/Appliscale/cloud-security-audit/configuration"
	"github.com/Appliscale/cloud-security-audit/csasession"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
)

type Users []*iam.User

func (u *Users) LoadFromAWS(config *configuration.Config, region string) error {
	iamAPI, err := config.ClientFactory.GetIAMClient(csasession.SessionConfig{Profile: config.Profile, Region: region})
	if err != nil {
		return err
	}

	q := &iam.ListUsersInput{}
	for {
		result, err := iamAPI.ListUsers(q)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case "OptInRequired":
					config.Logger.Warning("you are not subscribed to the IAM service in region: " + region)
					break
				default:
					return err
				}
			} else {
				return err
			}
		}
		*u = result.Users
	}
}
