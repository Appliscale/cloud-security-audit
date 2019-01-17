package resource

import (
	"github.com/Appliscale/cloud-security-audit/configuration"
	"github.com/Appliscale/cloud-security-audit/csasession"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
)

type Users []*iam.UserDetail

type IAMInfo struct {
	hasRootAccessKeys bool
	users             Users
}

func (iaminfo IAMInfo) GetUsers() Users {
	return iaminfo.users
}

func (iaminfo IAMInfo) HasRootAccessKeys() bool {
	return iaminfo.hasRootAccessKeys
}

func (iaminfo *IAMInfo) LoadFromAWS(config *configuration.Config, region string) error {
	iamAPI, err := config.ClientFactory.GetIAMClient(csasession.SessionConfig{Profile: config.Profile, Region: region})
	if err != nil {
		return err
	}

	q := &iam.GetAccountAuthorizationDetailsInput{}

	// get Users
	result, err := iamAPI.ListUsers(q)
	CheckError(region, config, err)
	(*iaminfo).users = result.UserDetailList
	//get access keys for root

	return nil
}

func CheckError(region string, config *configuration.Config, err error) error {
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
	return nil
}
