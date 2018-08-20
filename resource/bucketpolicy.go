package resource

type BucketPolicy struct {
	Version   string      `json:"Version,omitempty"`
	ID        string      `json:"Id,omitempty"`
	Statement []Statement `json:"Statement"`
}

type Statement struct {
	Sid       string      `json:"Sid"`
	Effect    string      `json:"Effect"`
	Principal interface{} `json:"Principal"`
	Action    interface{} `json:"Action"`
	Resource  string      `json:"Resource"`
	Condition Condition   `json:"Condition,omitempty"`
}

type Condition struct {
	Bool interface{} `json:",omitempty"`
	Null interface{} `json:",omitempty"`
}
