package report

import (
	"tyr/resource"

	"github.com/aws/aws-sdk-go/aws/session"
)

type Ec2Report struct {
	VolumeReport *VolumeReport
	InstanceID   string
	SortableTags *SortableTags
}

func NewEc2Report(instanceID string) *Ec2Report {
	return &Ec2Report{
		InstanceID:   instanceID,
		VolumeReport: &VolumeReport{},
		SortableTags: NewSortableTags(),
	}
}

type Ec2Reports []*Ec2Report

type ec2ReportRequiredResources struct {
	Ec2s    *resource.Ec2s
	KMSKeys *resource.KMSKeys
	Volumes *resource.Volumes
}

func (e *Ec2Reports) GetHeaders() []string {
	return []string{"EC2", "Volumes\n(None) - not encrypted\n(DKMS) - encrypted with default KMSKey", "EC2 Tags"}
}
func (e *Ec2Reports) FormatDataToTable() [][]string {
	data := [][]string{}

	for _, ec2Report := range *e {
		row := []string{
			ec2Report.InstanceID,
			ec2Report.VolumeReport.ToTableData(),
			ec2Report.SortableTags.ToTableData(),
		}
		data = append(data, row)
	}
	return data
}

func (e *Ec2Reports) GenerateReport(r *ec2ReportRequiredResources) {
	for _, ec2 := range *r.Ec2s {
		ec2Report := NewEc2Report(*ec2.InstanceId)
		ec2OK := true
		for _, blockDeviceMapping := range ec2.BlockDeviceMappings {
			volume := r.Volumes.FindById(*blockDeviceMapping.Ebs.VolumeId)
			if !*volume.Encrypted {
				ec2OK = false
				ec2Report.VolumeReport.AddEBS(*volume.VolumeId, NONE)
			} else {
				kmskey := r.KMSKeys.FindByKeyArn(*volume.KmsKeyId)
				if !kmskey.Custom {
					ec2OK = false
					ec2Report.VolumeReport.AddEBS(*volume.VolumeId, DKMS)
				}
			}
		}
		if !ec2OK {
			ec2Report.SortableTags.Add(ec2.Tags)
			*e = append(*e, ec2Report)
		}
	}
}

// GetResources : Initialize and loads required resources to create ec2 report
func (e *Ec2Reports) GetResources(sess *session.Session) (*ec2ReportRequiredResources, error) {
	resources := &ec2ReportRequiredResources{
		KMSKeys: resource.NewKMSKeys(),
		Ec2s:    &resource.Ec2s{},
		Volumes: &resource.Volumes{},
	}
	err := resource.LoadResources(
		sess,
		resources.Ec2s,
		resources.KMSKeys,
		resources.Volumes,
	)
	if err != nil {
		return nil, err
	}

	return resources, nil
}
