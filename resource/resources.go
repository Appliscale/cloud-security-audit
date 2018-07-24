package resource

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws/session"
)

type Resources struct {
	*Ec2s
	*Images
	*KMSKeys
	*Volumes
	*Snapshots
	*S3Buckets
}

type AWSResult struct {
	Error error
	Resource
}

func NewResources() *Resources {
	return &Resources{
		Ec2s:      &Ec2s{},
		Images:    &Images{},
		KMSKeys:   NewKMSKeys(),
		Volumes:   &Volumes{},
		Snapshots: &Snapshots{},
		S3Buckets: &S3Buckets{},
	}
}

func LoadResources(sess *session.Session, resources ...Resource) error {

	var wg sync.WaitGroup
	n := len(resources)
	wg.Add(n)
	errs := make(chan error, n)

	go func() {
		wg.Wait()
		close(errs)
	}()

	for _, r := range resources {
		go func(r Resource) {
			defer wg.Done()
			errs <- r.LoadFromAWS(sess)
		}(r)
	}

	for err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}
