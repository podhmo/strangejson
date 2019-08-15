package item2

import "github.com/podhmo/strangejson/_examples/manypackages04/item"

// FormatCheck : (generated from github.com/podhmo/strangejson/_examples/manypackages04/item2.Item2)
func (x *Item2) FormatCheck() error {
	return nil
}

// UnmarshalJSON : (generated from github.com/podhmo/strangejson/_examples/manypackages04/item2.Item2)
func (x *Item2) UnmarshalJSON(b []byte) error {
	return (*item.Item)(x).UnmarshalJSON(b)
}
