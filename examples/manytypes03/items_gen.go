package model

import (
	"encoding/json"
	"errors"
)

// UnmarshalJSON : (generated from github.com/podhmo/strangejson/examples/manytypes03.Item)
func (x *Item) UnmarshalJSON(b []byte) error {
	type internal struct {
		Name *string
	}

	var p internal
	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	if p.Name == nil {
		return errors.New("Name is required")
	}
	x.Name = *p.Name
	return nil
}

func (x *Item2) UnmarshalJSON(b []byte) error {
	return (*Item)(x).UnmarshalJSON(b)
}

func (x *Item4) UnmarshalJSON(b []byte) error {
	type internal struct {
		Name *string
	}

	var p internal
	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	if p.Name == nil {
		return errors.New("Name is required")
	}
	x.Name = *p.Name
	return nil
}

func (x *Item5) UnmarshalJSON(b []byte) error {
	type internal struct {
		Name *string `required:"false"`
	}

	var p internal
	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	if p.Name != nil {
		x.Name = *p.Name
	}
	return nil
}
