package scanner

import (
	"fmt"
	"strings"

	"github.com/Appliscale/cloud-security-audit/configuration"
	"github.com/Appliscale/cloud-security-audit/report/resourceReports"
)

func Run(config *configuration.Config) error {

	for _, service := range *config.Services {
		switch strings.ToLower(service) {
		case "ec2":
			config.Logger.Info("Gathering information about EC2s...")
			ec2Reports := resourceReports.Ec2Reports{}
			resources, err := ec2Reports.GetResources(config)
			if err != nil {
				return err
			}
			ec2Reports.GenerateReport(resources)
			config.PrintFormat(&ec2Reports)
		case "s3":
			config.Logger.Info("Gathering information about S3s...")
			s3BucketReports := resourceReports.S3BucketReports{}
			resources, err := s3BucketReports.GetResources(config)
			if err != nil {
				return err
			}
			s3BucketReports.GenerateReport(resources)
			config.PrintFormat(&s3BucketReports)
		default:
			return fmt.Errorf("Wrong service name: %s", service)
		}
	}
	return nil
}
