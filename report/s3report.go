package report

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Appliscale/cloud-security-audit/configuration"
	"github.com/Appliscale/cloud-security-audit/resource"
	"github.com/aws/aws-sdk-go/service/s3"
)

const action = "Action"
const effect = "Effect"
const principal = "Principal"

type S3BucketReport struct {
	Name string				`json:"name"`
	EncryptionType			`json:"encryption_type"`
	LoggingEnabled bool		`json:"logging_enabled"`
	ACLIsPublic    string	`json:"acl_is_public"`
	PolicyIsPublic string	`json:"policy_is_public"`
}

type S3BucketReports []*S3BucketReport

type S3ReportRequiredResources struct {
	KMSKeys   *resource.KMSKeys
	S3Buckets *resource.S3Buckets
}

func (s3brs S3BucketReports) GetJsonReport() ([]byte, error) {
	return json.Marshal(s3brs)
}

func (s3brs S3BucketReports) GetCsvReport() ([]byte, error) {
	output := make([]byte, 0)

	for _, row := range s3brs {
		s := fmt.Sprintf("%s,%s,%t,%s,%s\n", row.Name, row.EncryptionType, row.LoggingEnabled, row.ACLIsPublic,
			row.PolicyIsPublic)
		output = append(output, []byte(s)...)
	}

	return output, nil
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
	return []string{"Bucket Name", "Default\nSSE", "Logging\nEnabled", "ACL\nis public\nR - Read\nW - Write\nD - Delete", "Policy\nis public\nR - Read\nW - Write\nD - Delete"}
}

func (s3brs *S3BucketReports) FormatDataToTable() [][]string {
	data := [][]string{}

	for _, s3br := range *s3brs {
		row := []string{
			s3br.Name,
			s3br.EncryptionType.String(),
			strconv.FormatBool(s3br.LoggingEnabled),
			s3br.ACLIsPublic,
			s3br.PolicyIsPublic,
		}
		data = append(data, row)
	}
	return data
}

func isBucketACLPublic(s3Bucket *resource.S3Bucket) (bool, string) {

	grants := s3Bucket.ACL.Grants
	ownerID := s3Bucket.ACL.Owner.ID
	accessTypeACL := ""
	isPublic := false

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
		isPublic = false
		granteeURI := grant.Grantee.URI
		granteePermission := grant.Permission
		granteeID := grant.Grantee.ID
		if granteeID != ownerID {
			if granteeID != ownerID {
				if (granteeURI != nil && isStringInArray(*granteeURI, uriGroups)) &&
					(granteePermission != nil && isStringInArray(*granteePermission, permissionsACL)) {
					if !strings.Contains(accessTypeACL, getTypeOfAccessACL(*granteePermission)) {
						accessTypeACL = accessTypeACL + getTypeOfAccessACL(*granteePermission)
					}
					isPublic = true
				}
			}
		}
	}
	return isPublic, accessTypeACL

}

func getTypeOfAccessACL(granteePermission string) string {
	if granteePermission == "FULL_CONTROL" {
		return "RWD"
	}
	accessTypeACL := string(granteePermission[0])
	return accessTypeACL
}

func isStringInArray(element string, array []string) bool {
	for _, arrayElement := range array {
		if arrayElement == element {
			return true
		}
	}
	return false
}

func isBucketPolicyPublic(s3Bucket *resource.S3Bucket) (bool, string) {
	isPublic := make(map[string]bool)
	accessTypePolicy := ""
	if s3Bucket.S3Policy != nil {
		bucketPolicy := s3Bucket.S3Policy
		stat := bucketPolicy.Statements

		for _, element := range stat {
			isPublic[effect] = false
			isPublic[action] = false
			isPublic[principal] = false
			accessTypePolicy = ""

			//Effect
			if element.Effect == "Allow" {
				isPublic[effect] = true
			}
			//Action
			if len(element.Actions) > 0 {
				isPublic[action] = true
				accessTypePolicy = getTypeOfAccessPolicy(element.Actions)
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
			return true, accessTypePolicy
		}
	}
	return false, ""
}

func getTypeOfAccessPolicy(actions resource.Actions) string {
	var types string
	for _, action := range actions {
		if strings.Contains(action, "Get") || strings.Contains(action, "List") {
			types = types + "R"
		} else if strings.Contains(action, "Delete") {
			types = types + "D"
		} else if strings.Contains(action, "Put") || strings.Contains(action, "Create") {
			types = types + "W"
		}
	}
	types = "[" + types + "]"
	return types
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
		policy, policyTypes := isBucketPolicyPublic(s3Bucket)
		s3BucketReport.PolicyIsPublic = strconv.FormatBool(policy) + " " + policyTypes
		acl, aclTypes := isBucketACLPublic(s3Bucket)
		if len(aclTypes) == 0 {
			s3BucketReport.ACLIsPublic = strconv.FormatBool(acl)
		} else {
			s3BucketReport.ACLIsPublic = strconv.FormatBool(acl) + " " + "[" + aclTypes + "]"
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
