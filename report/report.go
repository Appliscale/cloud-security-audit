package report

import (
	"github.com/olekukonko/tablewriter"
	"os"
	"strings"
)

type Report interface {
	FormatDataToTable() [][]string
	GetHeaders() []string
}

func PrintTable(r Report) {

	data := r.FormatDataToTable()

	table := tablewriter.NewWriter(os.Stdout)
	// Configure Headers
	table.SetReflowDuringAutoWrap(false)
	table.SetAutoFormatHeaders(false)
	table.SetHeader(customFormatHeaders(r.GetHeaders()))
	// Configure rows&cells
	table.SetRowSeparator("-")
	table.SetRowLine(true)
	table.SetAutoWrapText(false)
	table.AppendBulk(data)
	table.Render()
}

func customFormatHeaders(headers []string) []string {
	for i, header := range headers {
		headers[i] = Title(header)
	}
	return headers
}

func Title(name string) string {
	origLen := len(name)
	name = strings.Replace(name, "_", " ", -1)
	//name = strings.Replace(name, ".", " ", -1)
	name = strings.TrimSpace(name)
	if len(name) == 0 && origLen > 0 {
		// Keep at least one character. This is important to preserve
		// empty lines in multi-line headers/footers.
		name = " "
	}
	return strings.ToUpper(name)
}

type EncryptionType int

const (
	// NONE : No Encryption
	NONE EncryptionType = 0
	// AES256 : 256-bit Advanced Encryption Standard
	AES256 EncryptionType = 1
	// DKMS : Encrypted with default KMS Key
	DKMS EncryptionType = 2
	// CKMS : Encrypted with Custom KMS Key
	CKMS EncryptionType = 3
)

func (et EncryptionType) String() string {
	types := [...]string{"NONE", "AES256", "DKMS", "CKMS"}
	if et < NONE || et > CKMS {
		return "Unknown"
	}
	return types[et]
}
