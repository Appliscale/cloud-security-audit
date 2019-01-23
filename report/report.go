package report

import (
	"github.com/Appliscale/cloud-security-audit/logger"
	"github.com/olekukonko/tablewriter"
	"os"
	"strings"
)

var ReportLogger = logger.CreateDefaultLogger()

type Report interface {
	FormatDataToTable() [][]string
	GetTableHeaders() []string
	GetJsonReport() []byte
	PrintHtmlReport(*os.File) error
	GetCsvReport() []byte
}

func handleOutput(output []byte, file *os.File) {
	out := append(output, []byte("\n\n")...)
	_, err := file.Write(out)
	if err != nil {
		panic(err)
	}
	file.Sync()
}

func PrintJsonReport(r Report, filename *os.File) {
	handleOutput(r.GetJsonReport(), filename)
}

func PrintHtmlReport(r Report, filename *os.File) {
	r.PrintHtmlReport(filename)
}

func PrintCSVReport(r Report, filename *os.File) {
	handleOutput(r.GetCsvReport(), filename)
}

func PrintTable(r Report, file *os.File) {

	data := r.FormatDataToTable()

	table := tablewriter.NewWriter(file)
	// Configure Headers
	table.SetReflowDuringAutoWrap(false)
	table.SetAutoFormatHeaders(false)
	table.SetHeader(customFormatHeaders(r.GetTableHeaders()))
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
