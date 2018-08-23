package resource

import (
	"encoding/json"
	"io/ioutil"
	"sync"

	"github.com/aws/aws-sdk-go/aws/session"
)

type Resource interface {
	LoadFromAWS(sess *session.Session) error
}

func LoadResource(r Resource, sess *session.Session) error {
	err := r.LoadFromAWS(sess)
	if err != nil {
		return err
	}

	return nil
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

func SaveToFile(r Resource, filename string) error {

	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0644)
}

func LoadFromFile(r Resource, filename string) error {

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, r)
}

func GetAvailableServices() *[]string {
	return &[]string{
		"ec2",
		"s3",
	}
}
