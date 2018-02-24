package item

import (
	"encoding/json"
	"errors"
)

// FormatCheck : (generated from github.com/podhmo/strangejson/examples/manypackages04/item.Item)
func (x *Item) FormatCheck() error {
	return nil
}

// UnmarshalJSON : (generated from github.com/podhmo/strangejson/examples/manypackages04/item.Item)
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
