package report

import (
	"strconv"

	"github.com/Appliscale/tyr/configuration"
	"github.com/Appliscale/tyr/resource"
	"github.com/Appliscale/tyr/tyrsession"

	"github.com/aws/aws-sdk-go/service/s3"
)

type S3BucketReport struct {
	Name string
	EncryptionType
	LoggingEnabled bool
}

type S3BucketReports []*S3BucketReport

type S3ReportRequiredResources struct {
	KMSKeys   *resource.KMSKeys
	S3Buckets *resource.S3Buckets
}

// CheckEncryptionType : Returns Encryption Type (AES256, CKMS, DKMS, NONE)
func (s3br *S3BucketReport) CheckEncryptionType(s3EncryptionType s3.ServerSideEncryptionByDefault, kmsKeys *resource.KMSKeys) {

	switch *s3EncryptionType.SSEAlgorithm {
	case "AES256":
		s3br.EncryptionType = AES256
	case "aws:kms":
		kmsKey := kmsKeys.FindByKeyArn(*s3EncryptionType.KMSMasterKeyID)
		if kmsKey.Custom {
			s3br.EncryptionType = CKMS
		} else {
			s3br.EncryptionType = DKMS
		}
	default:
		s3br.EncryptionType = NONE
	}
}

func (s3brs *S3BucketReports) GetHeaders() []string {
	return []string{"Bucket Name", "Default SSE", "Logging Enabled"}
}

func (s3brs *S3BucketReports) FormatDataToTable() [][]string {
	data := [][]string{}

	for _, s3br := range *s3brs {
		row := []string{
			s3br.Name,
			s3br.EncryptionType.String(),
			strconv.FormatBool(s3br.LoggingEnabled),
		}
		data = append(data, row)
	}
	return data
}

func (s3brs *S3BucketReports) GenerateReport(r *S3ReportRequiredResources) {

	for _, s3Bucket := range *r.S3Buckets {
		s3BucketReport := &S3BucketReport{Name: *s3Bucket.Name}
		ok := true
		if v := s3Bucket.ServerSideEncryptionConfiguration; v != nil {
			s3BucketReport.CheckEncryptionType(*v.Rules[0].ApplyServerSideEncryptionByDefault, r.KMSKeys)
		} else {
			s3BucketReport.EncryptionType = NONE
			ok = false
		}
		if s3Bucket.LoggingEnabled != nil {
			s3BucketReport.LoggingEnabled = true
		} else {
			ok = false
		}
		if !ok {
			*s3brs = append(*s3brs, s3BucketReport)
		}
	}
}

func (s3brs *S3BucketReports) GetResources(config *configuration.Config) (*S3ReportRequiredResources, error) {
	resources := &S3ReportRequiredResources{
		KMSKeys:   resource.NewKMSKeys(),
		S3Buckets: &resource.S3Buckets{},
	}

	sess, err := tyrsession.CreateSession(
		tyrsession.SessionConfig{
			Region:  (*config.Regions)[0],
			Profile: config.Profile,
		})
	if err != nil {
		return nil, err
	}

	err = resources.S3Buckets.LoadFromAWS(sess, config)
	if err != nil {
		return nil, err
	}
	err = resources.KMSKeys.LoadAllFromAWS(sess, config)
	if err != nil {
		return nil, err
	}

	return resources, nil
}
