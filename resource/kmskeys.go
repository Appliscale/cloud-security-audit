package resource

import (
	"fmt"
	"strings"
	"sync"
	"tyr/configuration"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

type KMSKey struct {
	AliasArn  string
	AliasName string
	Custom    bool
	KeyId     string // the same as TargetKeyId in AliasListEntry
}

type KMSKeys struct {
	Values map[string]*KMSKey
	sync.RWMutex
}

// NewKMSKeys : Initialize KMS Keys struct with map of keys
func NewKMSKeys() *KMSKeys {
	return &KMSKeys{Values: make(map[string]*KMSKey)}
}

type KMSKeyAliases []*kms.AliasListEntry

type KMSKeysListEntries []*kms.KeyListEntry

func getRegionMapOfKMSAPIs(sess *session.Session, config *configuration.Config) (map[string]*kms.KMS, error) {
	regions := []string{
		"us-east-2",
		"us-east-1",
		"us-west-1",
		"us-west-2",
		"ap-northeast-1",
		"ap-northeast-2",
		"ap-northeast-3",
		"ap-south-1",
		"ap-southeast-1",
		"ap-southeast-2",
		"ca-central-1",
		"eu-central-1",
		"eu-west-1",
		"eu-west-2",
		"eu-west-3",
		"sa-east-1",
	}
	regionSessions := make(map[string]*kms.KMS)
	for _, region := range regions {
		sess, err := session.NewSessionWithOptions(
			session.Options{
				Config: aws.Config{
					Region: &region,
				},
				Profile: config.Profile,
			},
		)
		if err == nil {
			regionSessions[region] = kms.New(sess)
		} else {
			return nil, err
		}
	}
	return regionSessions, nil
}

// LoadAllFromAWS : Load KMS Keys from all regions
func (k *KMSKeys) LoadAllFromAWS(sess *session.Session, config *configuration.Config) error {
	regionSessions, err := getRegionMapOfKMSAPIs(sess, config)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	n := len(regionSessions) * 2
	done := make(chan bool, n)
	errc := make(chan error, n)
	wg.Add(n)

	go func() {
		wg.Wait()
		close(done)
		close(errc)
	}()

	kmsKeyAliases := &KMSKeyAliases{}
	kmsKeyListEntries := &KMSKeysListEntries{}
	for _, kmsAPI := range regionSessions {
		go loadKeyListEntries(kmsAPI, kmsKeyListEntries, done, errc, &wg)
		go loadKeyAliases(kmsAPI, kmsKeyAliases, done, errc, &wg)
	}
	for i := 0; i < n; i++ {
		select {
		case <-done:
		case err := <-errc:
			return err
		}
	}

	k.loadValuesToMap(kmsKeyAliases, kmsKeyListEntries)

	return nil
}

func loadKeyListEntries(kmsAPI *kms.KMS, keyListEntries *KMSKeysListEntries, done chan bool, errc chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	q := &kms.ListKeysInput{}
	for {
		result, err := kmsAPI.ListKeys(q)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case "SubscriptionRequiredException":
					done <- true
				default:
					errc <- fmt.Errorf("[AWS-ERROR] Error Msg: %s", aerr.Error())
				}
			} else {
				errc <- fmt.Errorf("[ERROR] Error Msg: %s", err.Error())
			}
			return
		}
		if len(result.Keys) == 0 {
			done <- true
			return
		}

		*keyListEntries = append(*keyListEntries, result.Keys...)
		if !*result.Truncated {
			done <- true
			break
		}
		q.Marker = result.NextMarker
	}
}

func loadKeyAliases(kmsAPI *kms.KMS, aliases *KMSKeyAliases, done chan bool, errc chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	listAliasesInput := &kms.ListAliasesInput{}
	for {
		result, err := kmsAPI.ListAliases(listAliasesInput)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case "SubscriptionRequiredException":
					done <- true
				default:
					errc <- fmt.Errorf("[AWS-ERROR] Error Msg: %s", aerr.Error())
				}
			} else {
				errc <- fmt.Errorf("[ERROR] Error Msg: %s", err.Error())
			}
			return
		}
		*aliases = append(*aliases, result.Aliases...)

		if !*result.Truncated {
			done <- true
			break
		}
		listAliasesInput.Marker = result.NextMarker
	}
}

func (k *KMSKeys) LoadFromAWS(sess *session.Session) error {
	kmsAPI := kms.New(sess)

	var wg sync.WaitGroup
	n := 2
	done := make(chan bool, n)
	errc := make(chan error, n)
	wg.Add(n)

	go func() {
		wg.Wait()
		close(done)
		close(errc)
	}()

	kmsKeyAliases := &KMSKeyAliases{}
	kmsKeyListEntries := &KMSKeysListEntries{}

	go loadKeyListEntries(kmsAPI, kmsKeyListEntries, done, errc, &wg)
	go loadKeyAliases(kmsAPI, kmsKeyAliases, done, errc, &wg)

	for i := 0; i < n; i++ {
		select {
		case <-done:
		case err := <-errc:
			return err
		}
	}

	k.loadValuesToMap(kmsKeyAliases, kmsKeyListEntries)

	return nil
}

func (k *KMSKeys) loadValuesToMap(aliases *KMSKeyAliases, keyListEntries *KMSKeysListEntries) {
	for _, keyListEntry := range *keyListEntries {
		key := KMSKey{KeyId: *keyListEntry.KeyId}
		for _, alias := range *aliases {
			if alias.TargetKeyId != nil {
				if key.KeyId == *alias.TargetKeyId {
					key.AliasArn = *alias.AliasArn
					key.AliasName = *alias.AliasName
					if !strings.Contains(*alias.AliasName, "alias/aws/") {
						key.Custom = true
					}
					break
				}
			} else {
				key.Custom = true
			}
		}
		k.Values[*keyListEntry.KeyArn] = &key
	}
}

func (k *KMSKeys) FindByKeyArn(keyArn string) *KMSKey {
	kmsKey, ok := k.Values[keyArn]
	if ok {
		return kmsKey
	}
	return nil
}
