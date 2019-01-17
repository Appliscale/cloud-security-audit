package report

import (
	"fmt"
	"github.com/Appliscale/cloud-security-audit/configuration"
	"github.com/Appliscale/cloud-security-audit/resource"
	"github.com/aws/aws-sdk-go/service/iam"
)

type IAMReport struct {
	UserName       string
	InlinePolicies int
}

func NewIAMReport(u iam.UserDetail) *IAMReport {
	return &IAMReport{
		*u.UserName,
		len(u.UserPolicyList),
	}
}

type IAMReports []*IAMReport

type IAMReportRequiredResources struct {
	IAMInfo *resource.IAMInfo
}

func (i *IAMReports) GetHeaders() []string {
	return []string{"User name", "Inline policies"}
}

func (i *IAMReports) FormatDataToTable() [][]string {
	data := [][]string{}

	for _, iamReport := range *i {
		row := []string{
			iamReport.UserName,
			fmt.Sprintf("%d", iamReport.InlinePolicies),
		}
		data = append(data, row)
	}

	return data
}

func (i *IAMReports) GenerateReport(r *IAMReportRequiredResources) {
	for _, user := range (*r.IAMInfo).GetUsers() {
		iamReport := NewIAMReport(*user)
		*i = append(*i, iamReport)
	}

}

func (i *IAMReports) GetResources(config *configuration.Config) (*IAMReportRequiredResources, error) {
	resources := &IAMReportRequiredResources{IAMInfo: &resource.IAMInfo{}}

	for _, region := range *config.Regions {
		err := resource.LoadResources(
			config,
			region,
			resources.IAMInfo,
		)

		if err != nil {
			return nil, err
		}
	}

	return resources, nil
}
