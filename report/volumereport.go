package report

import (
	"bytes"
	"fmt"
)

type VolumeReport []string

func (v *VolumeReport) AddEBS(volumeID string, encryptionType EncryptionType) {
	*v = append(*v, volumeID+fmt.Sprintf("[%s]", encryptionType.String()))
}

func (v *VolumeReport) ToTableData() string {
	if len(*v) == 0 {
		return ""
	}
	var buffer bytes.Buffer
	n := len(*v)
	for _, volumeID := range (*v)[:n-1] {
		buffer.WriteString(volumeID + "\n")
	}
	buffer.WriteString((*v)[n-1])
	return buffer.String()
}
