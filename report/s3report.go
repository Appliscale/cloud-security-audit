package report

import (
	"strconv"

	"github.com/Appliscale/tyr/configuration"
	"github.com/Appliscale/tyr/resource"
	"github.com/aws/aws-sdk-go/service/s3"
)

const action = "Action"
const effect = "Effect"
const principal = "Principal"

type S3BucketReport struct {
	Name string
	EncryptionType
	LoggingEnabled bool
	ACLIsPublic    bool
	PolicyIsPublic bool
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
	return []string{"Bucket Name", "Default SSE", "Logging Enabled", "ACL - is public", "Policy - is public"}
}

func (s3brs *S3BucketReports) FormatDataToTable() [][]string {
	data := [][]string{}

	for _, s3br := range *s3brs {
		row := []string{
			s3br.Name,
			s3br.EncryptionType.String(),
			strconv.FormatBool(s3br.LoggingEnabled),
			strconv.FormatBool(s3br.ACLIsPublic),
			strconv.FormatBool(s3br.PolicyIsPublic),
		}
		data = append(data, row)
	}
	return data
}

func isBucketACLPublic(s3Bucket *resource.S3Bucket) bool {

	grants := s3Bucket.ACL.Grants
	ownerID := s3Bucket.ACL.Owner.ID

	var uriGroups = []string{
		"http://acs.amazonaws.com/groups/global/AuthenticatedUsers",
		"http://acs.amazonaws.com/groups/global/AllUsers",
	}
	var permissionsACL = []string{
		"READ",
		"WRITE",
		"READ_ACP",
		"WRITE_ACP",
		"FULL_CONTROL",
	}

	for _, grant := range grants {
		granteeURI := grant.Grantee.URI
		granteePermission := grant.Permission
		granteeID := grant.Grantee.ID
		if granteeID != ownerID {
			if granteeID != ownerID {
				if (granteeURI != nil && isStringInArray(*granteeURI, uriGroups)) &&
					(granteePermission != nil && isStringInArray(*granteePermission, permissionsACL)) {
					return true
				}
			}
		}
	}
	return false

}

func isStringInArray(element string, array []string) bool {
	for _, arrayElement := range array {
		if arrayElement == element {
			return true
		}
	}
	return false
}

func isBucketPolicyPublic(s3Bucket *resource.S3Bucket) bool {
	isPublic := make(map[string]bool)
	if s3Bucket.S3Policy != nil {
		bucketPolicy := s3Bucket.S3Policy
		stat := bucketPolicy.Statements

		for _, element := range stat {
			isPublic[effect] = false
			isPublic[action] = false
			isPublic[principal] = false

			//Effect
			if element.Effect == "Allow" {
				isPublic[effect] = true
			}
			//Action
			if len(element.Actions) > 0 {
				isPublic[action] = true
			}
			//Principal
			if element.Principal.Wildcard != "" && element.Principal.Wildcard == "*" {
				isPublic[principal] = true
			} else if len(element.Principal.Map) > 0 {
				for _, array := range element.Principal.Map {
					for _, principal := range array {
						if principal == "*" {
							isPublic[principal] = true
						}
					}
				}
			}
		}
		if isPublic[action] && isPublic[effect] && isPublic[principal] {
			return true
		}
	}
	return false
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
		s3BucketReport.PolicyIsPublic = isBucketPolicyPublic(s3Bucket)
		s3BucketReport.ACLIsPublic = isBucketACLPublic(s3Bucket)

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

	//sess, err := tyrsession.CreateSession(
	//	tyrsession.SessionConfig{
	//		Region:  (*config.Regions)[0],
	//		Profile: config.Profile,
	//	})
	//if err != nil {
	//	return nil, err
	//}

	err := resources.S3Buckets.LoadFromAWS(config, (*config.Regions)[0])
	if err != nil {
		return nil, err
	}
	err = resources.KMSKeys.LoadAllFromAWS(config)
	if err != nil {
		return nil, err
	}

	return resources, nil
}
