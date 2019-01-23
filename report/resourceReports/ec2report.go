package resourceReports

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Appliscale/cloud-security-audit/configuration"
	"github.com/Appliscale/cloud-security-audit/environment"
	"github.com/Appliscale/cloud-security-audit/report"
	"github.com/Appliscale/cloud-security-audit/resource"
	"os"
	"sort"
	"strconv"
	"strings"
	"html/template"
)

type Ec2Report struct {
	VolumeReport      *VolumeReport `json:"volume_report"`
	InstanceID        string        `json:"instance_id"`
	SortableTags      *SortableTags `json:"sortable_tags"`
	SecurityGroupsIDs []string      `json:"security_groups_ids"`
	AvailabilityZone  string        `json:"availability_zone"`
}

func NewEc2Report(instanceID string) *Ec2Report {
	return &Ec2Report{
		InstanceID:   instanceID,
		VolumeReport: &VolumeReport{},
		SortableTags: NewSortableTags(),
	}
}

type Ec2Reports []*Ec2Report

type Ec2ReportRequiredResources struct {
	Ec2s             *resource.Ec2s
	KMSKeys          *resource.KMSKeys
	Volumes          *resource.Volumes
	SecurityGroups   *resource.SecurityGroups
	AvailabilityZone string
}

func (e Ec2Reports) GetJsonReport() []byte {
	output, err := json.Marshal(e)
	if err == nil {
		return output
	}
	report.ReportLogger.Error("Error generating Json report")
	os.Exit(1)
	return []byte{}
}

func (e Ec2Reports) PrintHtmlReport(outputFile *os.File) error {
	data := e.GetJsonReport()
	reportTemplate := GetHtmlTemplate()

	tmpl, err := template.New("report_template").Parse(reportTemplate)
	if err != nil {
		report.ReportLogger.Error("Can not parse reportTemplate")
	}

	err = tmpl.Execute(outputFile, map[string]string{"EC2_JSON_PLACEHOLDER": string(data)})
	if err != nil {
		report.ReportLogger.Error("Can not execute reportTemplate")
	}

	return err
}

func (e Ec2Reports) GetCsvReport() []byte {
	const internalSep = ";"
	const externalSep = ","

	csv := []string{strings.Join([]string{
		"\"Availability Zone\"",
		"\"EC2\"",
		"\"Volumes encrypton\"",
		"\"Open Security Groups\"",
		"\"EC2 Tags\"",
	}, externalSep)}

	for _, row := range e {
		rowStr := make([]string, 0)

		rowStr = append(rowStr, strings.Join(*row.VolumeReport, internalSep))
		rowStr = append(rowStr, row.InstanceID)

		for tag, val := range row.SortableTags.Tags {
			rowStr = append(rowStr, fmt.Sprintf("%s%s%s", tag, internalSep, val))
		}

		rowStr = append(rowStr, strings.Join(row.SecurityGroupsIDs, internalSep))

		csv = append(csv, strings.Join(rowStr, externalSep))
	}

	return []byte(strings.Join(csv, "\n"))
}

func (e *Ec2Reports) GetTableHeaders() []string {
	return []string{"Availability\nZone", "EC2", "Volumes\n(None) - not encrypted\n(DKMS) - encrypted with default KMSKey", "Security\nGroups\n(Incoming CIDR = 0\x2E0\x2E0\x2E0/0)\nID : PROTOCOL : PORT", "EC2 Tags"}
}

func (e *Ec2Reports) FormatDataToTable() [][]string {
	data := [][]string{}

	for _, ec2Report := range *e {
		row := []string{
			ec2Report.AvailabilityZone,
			ec2Report.InstanceID,
			ec2Report.VolumeReport.ToTableData(),
			SliceOfStringsToString(ec2Report.SecurityGroupsIDs),
			ec2Report.SortableTags.ToTableData(),
		}
		data = append(data, row)
	}
	sortedData := sortTableData(data)
	return sortedData
}

func (e *Ec2Reports) GenerateReport(r *Ec2ReportRequiredResources) {
	for _, ec2 := range *r.Ec2s {
		ec2Report := NewEc2Report(*ec2.InstanceId)
		ec2OK := true
		for _, blockDeviceMapping := range ec2.BlockDeviceMappings {
			volume := r.Volumes.FindById(*blockDeviceMapping.Ebs.VolumeId)
			if !*volume.Encrypted {
				ec2OK = false
				ec2Report.VolumeReport.AddEBS(*volume.VolumeId, report.NONE)
			} else {
				kmskey := r.KMSKeys.FindByKeyArn(*volume.KmsKeyId)
				if !kmskey.Custom {
					ec2OK = false
					ec2Report.VolumeReport.AddEBS(*volume.VolumeId, report.DKMS)
				}
			}
		}

		for _, sg := range ec2.SecurityGroups {
			ipPermissions := r.SecurityGroups.GetIpPermissionsByID(*sg.GroupId)
			if ipPermissions != nil {
				for _, ipPermission := range ipPermissions {
					for _, ipRange := range ipPermission.IpRanges {
						if *ipRange.CidrIp == "0.0.0.0/0" {
							ec2Report.SecurityGroupsIDs = append(ec2Report.SecurityGroupsIDs, *sg.GroupId+" : "+*ipPermission.IpProtocol+" : "+strconv.FormatInt(*ipPermission.ToPort, 10))
							ec2OK = false
						}
					}
				}
			}
		}
		if !ec2OK {
			ec2Report.SortableTags.Add(ec2.Tags)
			*e = append(*e, ec2Report)
		}
		ec2Report.AvailabilityZone = *ec2.Placement.AvailabilityZone
	}
}

// GetResources : Initialize and loads required resources to create ec2 report
func (e *Ec2Reports) GetResources(config *configuration.Config) (*Ec2ReportRequiredResources, error) {
	resources := &Ec2ReportRequiredResources{
		KMSKeys:          resource.NewKMSKeys(),
		Ec2s:             &resource.Ec2s{},
		Volumes:          &resource.Volumes{},
		SecurityGroups:   &resource.SecurityGroups{},
		AvailabilityZone: "zone",
	}

	for _, region := range *config.Regions {
		err := resource.LoadResources(
			config,
			region,
			resources.Ec2s,
			resources.KMSKeys,
			resources.Volumes,
			resources.SecurityGroups,
		)
		if err != nil {
			return nil, err
		}
	}
	return resources, nil
}

func SliceOfStringsToString(slice []string) string {
	n := len(slice)
	if n == 0 {
		return ""
	}
	var buffer bytes.Buffer
	for _, s := range slice[:n-1] {
		buffer.WriteString(s + "\n")
	}
	buffer.WriteString(slice[n-1])
	return buffer.String()
}

func sortTableData(data [][]string) [][]string {
	if data[0][0] == "" {
		return data
	}
	var regions []string
	var sortedData [][]string
	for _, regs := range data {
		reg := regs[0][:len(regs[0])-1]
		regions = append(regions, reg)
	}
	sort.Strings(regions)
	uniqueregions := environment.UniqueNonEmptyElementsOf(regions)
	for _, unique := range uniqueregions {
		for _, b := range data {
			if strings.Contains(b[0], unique) {
				sortedData = append(sortedData, b)
			}
		}
	}
	return sortedData
}
