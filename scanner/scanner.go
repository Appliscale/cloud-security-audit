package scanner

import (
	"fmt"
	"strings"

	"github.com/Appliscale/cloud-security-audit/configuration"
	"github.com/Appliscale/cloud-security-audit/report"
)

func Run(config *configuration.Config) error {

	for _, service := range *config.Services {
		switch strings.ToLower(service) {
		case "ec2":
			config.Logger.Info("Gathering information about EC2s...")
			ec2Reports := report.Ec2Reports{}
			resources, err := ec2Reports.GetResources(config)
			if err != nil {
				return err
			}
			ec2Reports.GenerateReport(resources)
			report.PrintTable(&ec2Reports)
		case "s3":
			config.Logger.Info("Gathering information about S3s...")
			s3BucketReports := report.S3BucketReports{}
			resources, err := s3BucketReports.GetResources(config)
			if err != nil {
				return err
			}
			s3BucketReports.GenerateReport(resources)
			report.PrintTable(&s3BucketReports)
		case "iam":
			config.Logger.Info("Gathering information about IAM...")
			iamReports := report.IAMReports{}
			resources, err := iamReports.GetResources(config)
			if err != nil {
				return err
			}
			iamReports.GenerateReport(resources)
			report.PrintTable(&iamReports)

			iamChecklist := report.IAMChecklist{}
			iamChecklist.GenerateReport(resources)
			report.PrintTable(&iamChecklist)

		default:
			return fmt.Errorf("Wrong service name: %s", service)
		}
	}
	return nil
}
