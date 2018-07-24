package resource

import (
	"encoding/json"
	"io/ioutil"

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
