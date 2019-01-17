package resourceReports

type CsvReport interface {
	GetCsvReport() ([]byte, error)
}
