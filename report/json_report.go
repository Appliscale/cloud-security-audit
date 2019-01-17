package report

type JsonReport interface {
	GetJsonReport() ([]byte, error)
}
