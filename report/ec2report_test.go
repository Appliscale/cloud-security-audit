package report

import (
	"sort"
	"strings"
	"testing"

	"github.com/Appliscale/tyr/resource"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/assert"
)

func TestEC2Report_IfVolumeEncryptedWithDefaultKMSLineHasSuffixDKMS(t *testing.T) {
	kmsKeyArn := "arn:aws:kms:us-east-1:123456789101:key/abcdefgh-1234-5678-1234-abcdefghijkl"
	volumeID := "vol-111"
	volume := &ec2.Volume{
		VolumeId:  &volumeID,
		Encrypted: aws.Bool(true),
		KmsKeyId:  &kmsKeyArn,
	}
	instanceBlockDeviceMapping := &ec2.InstanceBlockDeviceMapping{
		Ebs: &ec2.EbsInstanceBlockDevice{
			VolumeId: &volumeID,
		},
	}
	instance := &ec2.Instance{
		InstanceId:          &volumeID,
		BlockDeviceMappings: []*ec2.InstanceBlockDeviceMapping{instanceBlockDeviceMapping},
	}
	kmsKeys := resource.NewKMSKeys()
	kmsKeys.Values[kmsKeyArn] = &resource.KMSKey{Custom: false}
	r := &Ec2ReportRequiredResources{
		Ec2s:             &resource.Ec2s{instance},
		Volumes:          &resource.Volumes{volume},
		KMSKeys:          kmsKeys,
		AvailabilityZone: "Zone",
	}
	ec2Reports := &Ec2Reports{}
	ec2Reports.GenerateReport(r)

	assert.Equal(t, volumeID+"[DKMS]", (*(*ec2Reports)[0].VolumeReport)[0])
}

func TestEC2Report_IfVolumeNotEncryptedLineHasSuffixNone(t *testing.T) {

	volumeID := "vol-111"
	volume := &ec2.Volume{
		VolumeId:  &volumeID,
		Encrypted: aws.Bool(false),
	}
	instanceBlockDeviceMapping := &ec2.InstanceBlockDeviceMapping{
		Ebs: &ec2.EbsInstanceBlockDevice{
			VolumeId: &volumeID,
		},
	}
	instance := &ec2.Instance{
		InstanceId:          &volumeID,
		BlockDeviceMappings: []*ec2.InstanceBlockDeviceMapping{instanceBlockDeviceMapping},
	}
	r := &Ec2ReportRequiredResources{
		Ec2s:    &resource.Ec2s{instance},
		Volumes: &resource.Volumes{volume},
		KMSKeys: resource.NewKMSKeys(),
	}
	ec2Reports := &Ec2Reports{}
	ec2Reports.GenerateReport(r)

	assert.Equal(t, volumeID+"[NONE]", (*(*ec2Reports)[0].VolumeReport)[0])
}

func TestEC2Report_IfTagsInTableDataAreFormattedCorrectly(t *testing.T) {

	instanceID := "i-1"
	ec2Tags := []*ec2.Tag{
		{
			Key:   aws.String("key1"),
			Value: aws.String("value"),
		},
		{
			Key:   aws.String("key2"),
			Value: aws.String("some value2"),
		},
	}

	ec2Report1 := NewEc2Report(instanceID)
	ec2Report1.SortableTags.Add(ec2Tags)

	ec2Reports := Ec2Reports{ec2Report1}
	tableData := ec2Reports.FormatDataToTable()

	formattedTags := strings.Split(tableData[0][3], "\n")
	sort.Strings(formattedTags)

	assert.True(t, len(formattedTags) == 2)
	for i, tag := range ec2Tags {
		assert.EqualValues(t, *tag.Key+":"+*tag.Value, formattedTags[i])
	}

}

func TestVolumeIfSliceIsEmptyDoNotPanic(t *testing.T) {
	v := &VolumeReport{}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code did panic")
		}
	}()

	v.ToTableData()
}
