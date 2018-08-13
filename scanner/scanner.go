package scanner

import (
	"fmt"
	"strings"
	"github.com/Appliscale/tyr/configuration"
	"github.com/Appliscale/tyr/report"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/aws/session"
)

func Run(config *configuration.Config) error {

	sess, err := session.NewSessionWithOptions(
		session.Options{
			Config: aws.Config{
				Region: &config.Region,
			},
			Profile: config.Profile,
		},
	)
	if err != nil {
		return err
	}

	switch strings.ToLower(config.Service) {
	case "ec2":
		ec2Reports := report.Ec2Reports{}
		resources, err := ec2Reports.GetResources(sess)
		if err != nil {
			return err
		}
		ec2Reports.GenerateReport(resources)
		report.PrintTable(&ec2Reports)
	case "s3":
		s3BucketReports := report.S3BucketReports{}
		resources, err := s3BucketReports.GetResources(sess, config)
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
