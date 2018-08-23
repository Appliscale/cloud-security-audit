package scanner

import (
	"fmt"
	"strings"

	"github.com/Appliscale/tyr/configuration"
	"github.com/Appliscale/tyr/report"
)

func Run(config *configuration.Config) error {

	switch strings.ToLower(config.Service) {
	case "ec2":
		ec2Reports := report.Ec2Reports{}
		resources, err := ec2Reports.GetResources(config)
		if err != nil {
			return err
		}
		ec2Reports.GenerateReport(resources)
		report.PrintTable(&ec2Reports)
	case "s3":
		s3BucketReports := report.S3BucketReports{}
		resources, err := s3BucketReports.GetResources(config)
		if err != nil {
			return err
		}
		s3BucketReports.GenerateReport(resources)
		report.PrintTable(&s3BucketReports)
	default:
		return fmt.Errorf("Wrong service name: %s", config.Service)
	}
	return nil
}
