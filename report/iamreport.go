package report

import (
	"github.com/Appliscale/cloud-security-audit/configuration"
	"github.com/Appliscale/cloud-security-audit/resource"
)

type IAMReport struct {
	UserName string
}

func NewIAMReport(userName string) *IAMReport {
	return &IAMReport{UserName: userName}
}

type IAMReports []*IAMReport

type IAMReportRequiredResources struct {
	Users *resource.Users
}

func (i *IAMReports) GetHeaders() []string {
	return []string{"User name"}
}

func (i *IAMReports) FormatDataToTable() [][]string {
	data := [][]string{}

	for _, iamReport := range *i {
		row := []string{iamReport.UserName}
		data = append(data, row)
	}

	return data
}

func (i *IAMReports) GenerateReport(r *IAMReportRequiredResources) {
	for _, user := range *r.Users {
		iamReport := NewIAMReport(*user.UserName)
		*i = append(*i, iamReport)
	}
}

func (i *IAMReports) GetResources(config *configuration.Config) (*IAMReportRequiredResources, error) {
	resources := &IAMReportRequiredResources{Users: &resource.Users{}}

	return resources, nil
}
