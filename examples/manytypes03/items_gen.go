package model

import (
	"encoding/json"
	"errors"
)

// FormatCheck : (generated from github.com/podhmo/strangejson/examples/manytypes03.Item)
func (x *Item) FormatCheck() error {
	return nil
}

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
	return x.FormatCheck()
}

// FormatCheck : (generated from github.com/podhmo/strangejson/examples/manytypes03.Item2)
func (x *Item2) FormatCheck() error {
	return nil
}

// UnmarshalJSON : (generated from github.com/podhmo/strangejson/examples/manytypes03.Item2)
func (x *Item2) UnmarshalJSON(b []byte) error {
	return (*Item)(x).UnmarshalJSON(b)
}

// FormatCheck : (generated from github.com/podhmo/strangejson/examples/manytypes03.Item4)
func (x *Item4) FormatCheck() error {
	return nil
}

// UnmarshalJSON : (generated from github.com/podhmo/strangejson/examples/manytypes03.Item4)
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
	return x.FormatCheck()
}

// FormatCheck : (generated from github.com/podhmo/strangejson/examples/manytypes03.Item5)
func (x *Item5) FormatCheck() error {
	return nil
}

// UnmarshalJSON : (generated from github.com/podhmo/strangejson/examples/manytypes03.Item5)
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
	return x.FormatCheck()
}
