package report

type CsvReport interface {
	GetCsvReport() ([]byte, error)
}
