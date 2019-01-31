package resourceReports

import (
	"encoding/json"
	"fmt"
	"github.com/Appliscale/cloud-security-audit/configuration"
	"github.com/Appliscale/cloud-security-audit/report"
	"github.com/Appliscale/cloud-security-audit/resource"
	"github.com/aws/aws-sdk-go/service/iam"
	"os"
	"strconv"
	"strings"
)

type IAMReport struct {
	UserName       string
	Groups         string
	InlinePolicies int
}

func NewIAMReport(u iam.UserDetail) *IAMReport {

	return &IAMReport{
		getUserName(u),
		getUserGroups(u),
		getNoOfUserInlinePolicies(u),
	}
}

type IAMReports []*IAMReport

type IAMReportRequiredResources struct {
	IAMInfo *resource.IAMInfo
}

func (i *IAMReports) GetTableHeaders() []string {
	return []string{"User name", "Groups", "# of Inline\npolicies"}
}

func (i *IAMReports) FormatDataToTable() [][]string {
	data := [][]string{}

	for _, iamReport := range *i {
		row := []string{
			iamReport.UserName,
			iamReport.Groups,
			fmt.Sprintf("%d", iamReport.InlinePolicies),
		}
		data = append(data, row)
	}

	return data
}

func (i *IAMReports) GetJsonReport() []byte {
	output, err := json.Marshal(i)
	if err == nil {
		return output
	}
	report.ReportLogger.Error("Error generating Json report")
	os.Exit(1)
	return []byte{}
}

func (i *IAMReports) PrintHtmlReport(*os.File) error {
	//	TODO:
	return nil
}

func (i IAMReports) GetCsvReport() []byte {
	const externalSep = ","

	csv := []string{strings.Join([]string{
		"\"UserName\"",
		"\"Groups\"",
		"\"Inline Policies\""}, externalSep)}

	for _, row := range i {
		s := strings.Join([]string{
			row.UserName,
			row.Groups,
			strconv.FormatInt(int64(row.InlinePolicies), 10)}, externalSep)

		csv = append(csv, s)
	}

	return []byte(strings.Join(csv, "\n"))
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

func getUserName(u iam.UserDetail) string {
	return *u.UserName
}

func getUserGroups(u iam.UserDetail) string {
	var grouplist []string
	if len(u.GroupList) != 0 {
		for _, group := range u.GroupList {
			grouplist = append(grouplist, *group)
		}
	}

	groups := "None"

	if len(grouplist) > 0 {
		groups = strings.Join(grouplist, ",")
	}

	return groups
}

func getNoOfUserInlinePolicies(u iam.UserDetail) int {
	return len(u.UserPolicyList)
}
