package platio

import (
	"encoding/json"
)

type RecordId = string

type Record struct {
	Id     RecordId `json:"id"`
	Values Values   `json:"values"`
}

func (record *Record) Name() string {
	if record.Values.Name != nil {
		return record.Values.Name.Value
	} else {
		return ""
	}
}

func (record *Record) Age() int {
	if record.Values.Age != nil {
		return int(record.Values.Age.Value)
	} else {
		return 0
	}
}

type Values struct {
	Name *StringValue `json:"cd33ed98,omitempty"`
	Age  *NumberValue `json:"ce0f2361,omitempty"`
}

type StringValue struct {
	Value string
}

func (value *StringValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}{"String", value.Value})
}

type NumberValue struct {
	Value float64
}

func (value *NumberValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type  string  `json:"type"`
		Value float64 `json:"value"`
	}{"Number", value.Value})
}
