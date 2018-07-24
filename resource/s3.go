package resource

import (
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Bucket struct {
	*s3.Bucket
	S3Policy S3Policy
	Region   *string
	*s3.ServerSideEncryptionConfiguration
	*s3.LoggingEnabled
}

type S3Buckets []*S3Bucket

type S3Policy interface{}

// type S3Policy struct {
// 	Version   string //`json:",omitempty"`
// 	Id        string //`json:",omitempty"`
// 	Statement interface{}
// }

func (b *S3Buckets) LoadRegions(sess *session.Session) error {
	sess.Handlers.Unmarshal.PushBackNamed(s3.NormalizeBucketLocationHandler)
	s3API := s3.New(sess)

	wg := sync.WaitGroup{}
	n := len(*b)
	wg.Add(n)
	done := make(chan bool, n)
	cerrs := make(chan error, n)

	go func() {
		wg.Wait()
		close(done)
		close(cerrs)
	}()

	for _, bucket := range *b {
		go func(s3Bucket *S3Bucket) {
			result, err := s3API.GetBucketLocation(&s3.GetBucketLocationInput{Bucket: s3Bucket.Name})
			if err != nil {
				cerrs <- err
				return
			}
			s3Bucket.Region = result.LocationConstraint
			done <- true
		}(bucket)
	}
	for i := 0; i < n; i++ {
		select {
		case <-done:
		case err := <-cerrs:
			return err
		}
	}

	return nil
}

// LoadNames : Get All S3 Bucket names
func (b *S3Buckets) LoadNames(sess *session.Session) error {
	s3API := s3.New(sess)

	result, err := s3API.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return err
	}
	for _, bucket := range result.Buckets {
		*b = append(*b, &S3Bucket{Bucket: bucket})
	}
	return nil
}

func getRegionMapOfS3APIs(s3Buckets S3Buckets) (map[string]*s3.S3, error) {
	regionS3APIs := make(map[string]*s3.S3)
	for _, bucket := range s3Buckets {
		if _, ok := regionS3APIs[*bucket.Region]; !ok {
			if sess, err := session.NewSession(&aws.Config{Region: bucket.Region}); err == nil {
				regionS3APIs[*bucket.Region] = s3.New(sess)
			} else {
				return nil, err
			}
		}
		// TODO : Add some check to stop iteration
		// if len(regionS3APIs) >= 17 {
		// 	break
		// }
	}
	return regionS3APIs, nil
}

func (b *S3Buckets) LoadFromAWS(sess *session.Session) error {

	err := b.LoadNames(sess)
	if err != nil {
		return err
	}

	err = b.LoadRegions(sess)
	if err != nil {
		return err
	}

	regionS3APIs, err := getRegionMapOfS3APIs(*b)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	n := 2 * len(*b)
	done := make(chan bool, n)
	errs := make(chan error, n)
	wg.Add(n)

	go func() {
		wg.Wait()
		close(done)
		close(errs)
	}()

	for _, s3Bucket := range *b {
		regionS3API := regionS3APIs[*s3Bucket.Region]
		// TODO : Need to add struct for s3 bucket policy
		// go getPolicy(s3Bucket, regionS3API, done, errs, &wg)
		go getEncryption(s3Bucket, regionS3API, done, errs, &wg)
		go getBucketLogging(s3Bucket, regionS3API, done, errs, &wg)
	}
	for i := 0; i < n; i++ {
		select {
		case <-done:
		case err := <-errs:
			return err
		}
	}
	return nil
}

func getPolicy(s3Bucket *S3Bucket, s3API *s3.S3, done chan bool, errs chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	result, err := s3API.GetBucketPolicy(&s3.GetBucketPolicyInput{
		Bucket: s3Bucket.Name,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NoSuchBucketPolicy":
				done <- true
			default:
				errs <- fmt.Errorf("[AWS-ERROR] Bucket: %s  Error Msg: %s", *s3Bucket.Name, aerr.Error())
			}
		} else {
			errs <- fmt.Errorf("[ERROR] %s: %s", *s3Bucket.Name, err.Error())
		}
		return
	}
	if result.Policy != nil {
		s3Bucket.S3Policy = *result.Policy
	}
	done <- true
}

func getEncryption(s3Bucket *S3Bucket, s3API *s3.S3, done chan bool, errs chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	result, err := s3API.GetBucketEncryption(&s3.GetBucketEncryptionInput{Bucket: s3Bucket.Name})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "ServerSideEncryptionConfigurationNotFoundError":
				done <- true
			default:
				errs <- fmt.Errorf("[AWS-ERROR] \nBucket: %s \n Error Msg: %s", *s3Bucket.Name, aerr.Error())
			}
		} else {
			errs <- fmt.Errorf("[ERROR] %s: %s", *s3Bucket.Name, err.Error())
		}
		return
	}

	if result.ServerSideEncryptionConfiguration != nil {
		s3Bucket.ServerSideEncryptionConfiguration = result.ServerSideEncryptionConfiguration
	}
	done <- true
}

func getBucketLogging(s3Bucket *S3Bucket, s3API *s3.S3, done chan bool, errs chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	result, err := s3API.GetBucketLogging(&s3.GetBucketLoggingInput{Bucket: s3Bucket.Name})
	if err != nil {
		errs <- fmt.Errorf("[ERROR] %s: %s", *s3Bucket.Name, err.Error())
		return
	}
	s3Bucket.LoggingEnabled = result.LoggingEnabled
	done <- true
}
