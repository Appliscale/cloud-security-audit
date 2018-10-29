package resource

import (
	"encoding/json"
	"io/ioutil"
	"sync"

	"github.com/Appliscale/cloud-security-audit/configuration"
)

type Resource interface {
	LoadFromAWS(config *configuration.Config, region string) error
}

func LoadResource(r Resource, config *configuration.Config, region string) error {
	err := r.LoadFromAWS(config, region)
	if err != nil {
		return err
	}

	return nil
}

func LoadResources(config *configuration.Config, region string, resources ...Resource) error {

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
			errs <- r.LoadFromAWS(config, region)
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
