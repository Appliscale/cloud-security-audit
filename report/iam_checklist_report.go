package report

import (
	"fmt"
	"github.com/Appliscale/cloud-security-audit/configuration"
	"github.com/Appliscale/cloud-security-audit/resource"
	//"github.com/aws/aws-sdk-go/service/iam"
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

func (i *IAMChecklist) GetHeaders() []string {
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
