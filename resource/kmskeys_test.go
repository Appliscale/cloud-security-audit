package resource

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/stretchr/testify/assert"
)

func TestKMSKeysloadValuesToMapSetCustomToTrueWhenAliasAWSIsNotPresent(t *testing.T) {
	prefix := "arn:aws:kms:us-east-1:123456789101:"
	aliasName := "alias/some-alias"
	keyID := "abcdefgh-1234-5678-1234-abcdefghijkl"
	keyArn := prefix + "key/" + keyID
	aliases := &KMSKeyAliases{
		&kms.AliasListEntry{
			TargetKeyId: &keyID,
			AliasArn:    aws.String(prefix + aliasName),
			AliasName:   &aliasName,
		},
	}
	keyListEntries := &KMSKeysListEntries{
		&kms.KeyListEntry{
			KeyArn: &keyArn,
			KeyId:  &keyID,
		},
	}
	kmsKeys := NewKMSKeys()
	kmsKeys.loadValuesToMap(aliases, keyListEntries)
	assert.True(t, kmsKeys.Values[keyArn].Custom)
}

func TestKMSKeysloadValuesToMapSetCustomToFalseWhenAliasAWSIsPresent(t *testing.T) {
	prefix := "arn:aws:kms:us-east-1:123456789101:"
	aliasName := "alias/aws/ebs"
	keyID := "abcdefgh-1234-5678-1234-abcdefghijkl"
	keyArn := prefix + "key/" + keyID
	aliases := &KMSKeyAliases{
		&kms.AliasListEntry{
			TargetKeyId: &keyID,
			AliasArn:    aws.String(prefix + aliasName),
			AliasName:   &aliasName,
		},
	}
	keyListEntries := &KMSKeysListEntries{
		&kms.KeyListEntry{
			KeyArn: &keyArn,
			KeyId:  &keyID,
		},
	}
	kmsKeys := NewKMSKeys()
	kmsKeys.loadValuesToMap(aliases, keyListEntries)
	assert.False(t, kmsKeys.Values[keyArn].Custom)
}
