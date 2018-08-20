package report

import (
	"fmt"
	"github.com/Appliscale/tyr/configuration"
	"github.com/Appliscale/tyr/resource"
	"strconv"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/s3"
)

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
	return []string{"Bucket Name", "Default SSE", "Logging Enabled", "ACL is public", "Policy is public"}
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

	var uriFlag bool
	var permissionFlag bool

	var uriGroups [2]string
	uriGroups[0] = "http://acs.amazonaws.com/groups/global/AuthenticatedUsers"
	uriGroups[1] = "http://acs.amazonaws.com/groups/global/AllUsers"

	var permissionsACL [5]string
	permissionsACL[0] = "READ"
	permissionsACL[1] = "WRITE"
	permissionsACL[2] = "READ_ACP"
	permissionsACL[3] = "WRITE_ACP"
	permissionsACL[4] = "FULL_CONTROL"

	for _, grant := range grants {
		granteeURI := grant.Grantee.URI
		granteePermission := grant.Permission
		granteeID := grant.Grantee.ID
		if granteeID != ownerID {
			for _, group := range uriGroups {
				if granteeURI != nil && *granteeURI == group {
					uriFlag = true
				}
			}
			for _, permission := range permissionsACL {
				if granteePermission != nil && *granteePermission == permission {
					permissionFlag = true
				}
			}
		}
		if uriFlag == true && permissionFlag == true {
			return true
		}
	}
	return false
}

func isBucketPolicyPublic(s3Bucket *resource.S3Bucket) bool {
	isPublic := make(map[string]bool)
	bucketPolicy := s3Bucket.S3Policy
	stat := bucketPolicy.Statement

	if len(stat) > 0 {
		for _, b := range stat {
			isPublic["Effect"] = false
			isPublic["Action"] = false
			isPublic["Principal"] = false
			actionType := checkType(b.Action)
			principalType := checkType(b.Principal)

			if actionType == "string" {
				if b.Action.(string) != "" {
					isPublic["Action"] = true
				}
			} else if actionType == "map" {
				if len(b.Action.(map[string]interface{})) > 0 {
					isPublic["Action"] = true
				}
			}
			if b.Effect == "Allow" {
				isPublic["Effect"] = true
			}
			if principalType == "string" {
				if b.Principal.(string) == "*" {
					isPublic["Principal"] = true
				}
			} else if principalType == "map" {

				pri := b.Principal.(map[string]interface{})
				if len(b.Principal.(map[string]interface{})) > 0 && pri["AWS"] == "*" {
					isPublic["Principal"] = true
				}
			}
		}
		counter := 0
		for _, value := range isPublic {
			if value == true {
				counter++
			}
		}
		if counter == 3 {
			return true
		}
	}
	return false
}

func checkType(i interface{}) string {
	switch v := i.(type) {
	case string:
		return "string"
	case map[string]interface{}:
		return "map"
	default:
		fmt.Printf("I don't know about type %T\n", v)
		return ""
	}
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

func (s3brs *S3BucketReports) GetResources(sess *session.Session, config *configuration.Config) (*S3ReportRequiredResources, error) {
	resources := &S3ReportRequiredResources{
		KMSKeys:   resource.NewKMSKeys(),
		S3Buckets: &resource.S3Buckets{},
	}

	err := resources.S3Buckets.LoadFromAWS(sess, config)
	if err != nil {
		return nil, err
	}
	err = resources.KMSKeys.LoadAllFromAWS(sess, config)
	if err != nil {
		return nil, err
	}

	return resources, nil
}
