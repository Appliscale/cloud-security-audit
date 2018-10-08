package resource

import (
	"encoding/json"
	"errors"
)

type S3Policy struct {
	Version    string
	Id         string      `json:",omitempty"`
	Statements []Statement `json:"Statement"`
}

func NewS3Policy(s string) (*S3Policy, error) {
	b := []byte(s)
	s3Policy := &S3Policy{}
	err := json.Unmarshal(b, s3Policy)
	if err != nil {
		return nil, err
	}
	return s3Policy, nil
}

type Statement struct {
	Effect    string
	Principal Principal `json:"Principal"`
	Actions   Actions   `json:"Action"`
	Resource  Resources `json:"Resource"`
	Condition Condition `json:",omitempty"`
}

type Condition struct {
	Bool map[string]string `json:",omitempty"`
	Null map[string]string `json:",omitempty"`
}

type Actions []string

func (a *Actions) UnmarshalJSON(b []byte) error {

	array := []string{}
	err := json.Unmarshal(b, &array)
	/*
		if error is: "json: cannot unmarshal string into Go value of type []string"
		then fallback to unmarshaling string
	*/
	if err != nil {
		s := ""
		err = json.Unmarshal(b, &s)
		if err != nil {
			return err
		}
		*a = append(*a, s)
		return nil
	}
	for _, action := range array {
		*a = append(*a, action)
	}
	return nil
}

// Principal : Specifies user, account, service or other
// 			   entity that is allowed or denied access to resource
type Principal struct {
	Map      map[string][]string // Values in Map: https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_principal.html
	Wildcard string              // Values: *
}

func (p *Principal) UnmarshalJSON(b []byte) error {
	p.Map = make(map[string][]string)
	s := ""
	err := json.Unmarshal(b, &s)
	if err != nil {
		m := make(map[string]interface{})

		err = json.Unmarshal(b, &m)
		if err != nil {
			return err
		}
		for key, value := range m {
			switch t := value.(type) {
			case string:
				p.Map[key] = append(p.Map[key], value.(string))
			case []interface{}:
				for _, elem := range value.([]interface{}) {
					p.Map[key] = append(p.Map[key], elem.(string))
				}
			default:
				return errors.New("Unsupported type " + t.(string))

			}
		}
	}
	p.Wildcard = s
	return nil
}

type Resources []string

func (r *Resources) UnmarshalJSON(b []byte) error {

	array := []string{}
	err := json.Unmarshal(b, &array)
	/*
		if error is: "json: cannot unmarshal string into Go value of type []string"
		then fallback to unmarshaling string
	*/
	if err != nil {
		s := ""
		err = json.Unmarshal(b, &s)
		if err != nil {
			return err
		}
		*r = append(*r, s)
		return nil
	}
	for _, resource := range array {
		*r = append(*r, resource)
	}
	return nil
}
