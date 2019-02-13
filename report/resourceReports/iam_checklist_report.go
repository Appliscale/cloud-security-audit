package resourceReports

import (
	"fmt"
	"github.com/Appliscale/cloud-security-audit/configuration"
	"github.com/Appliscale/cloud-security-audit/resource"
	//"github.com/aws/aws-sdk-go/service/iam"
	"encoding/json"
	"github.com/Appliscale/cloud-security-audit/report"
	"os"
	"strconv"
	"strings"
)

type IAMItem struct {
	Name  string
	Value bool
}

func NewIAMItem(name string, value bool) *IAMItem {
	return &IAMItem{
		name,
		value,
	}
}

type IAMChecklist []*IAMItem

type IAMChecklistRequiredResources struct {
	IAMInfo *resource.IAMInfo
}

func (i *IAMChecklist) GetTableHeaders() []string {
	return []string{"Guideline", "Status"}
}

func (i *IAMChecklist) FormatDataToTable() [][]string {
	data := [][]string{}

	for _, iamItem := range *i {
		row := []string{
			iamItem.Name,
			fmt.Sprintf("%t", iamItem.Value),
		}
		data = append(data, row)
	}

	return data
}

func (i *IAMChecklist) GetJsonReport() []byte {
	output, err := json.Marshal(i)
	if err == nil {
		return output
	}
	report.ReportLogger.Error("Error generating Json report")
	os.Exit(1)
	return []byte{}
}

func (i *IAMChecklist) PrintHtmlReport(*os.File) error {
	//	TODO:
	return nil
}

func (i IAMChecklist) GetCsvReport() []byte {
	const externalSep = ","

	csv := []string{strings.Join([]string{
		"\"Name\"",
		"\"Value\""}, externalSep)}

	for _, row := range i {
		s := strings.Join([]string{
			row.Name,
			strconv.FormatBool(row.Value)}, externalSep)

		csv = append(csv, s)
	}

	return []byte(strings.Join(csv, "\n"))
}

func (i *IAMChecklist) GenerateReport(r *IAMReportRequiredResources) {

	*i = append(*i, NewIAMItem("Root account's access keys locked away", !r.IAMInfo.HasRootAccessKeys()))

}

func (i *IAMChecklist) GetResources(config *configuration.Config) (*IAMReportRequiredResources, error) {
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
