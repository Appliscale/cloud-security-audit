package report

import (
	"bytes"
	"fmt"
)

type VolumeReport []string

// type VolumeViolation int

// const (
// 	None VolumeViolation = 0
// 	NE   VolumeViolation = 1 // Not Encrypted
// 	DKMS VolumeViolation = 2 // Encrypted with Default KMS Key
// )

// func (vv VolumeViolation) String() string {
// 	types := [...]string{"None", "NE", "DKMS"}
// 	if vv < None || vv > DKMS {
// 		return "Unknown"
// 	}
// 	return types[vv]
// }

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
