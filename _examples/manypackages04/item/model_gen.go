package item

import (
	"encoding/json"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

// FormatCheck : (generated from github.com/podhmo/strangejson/_examples/manypackages04/item.Item)
func (x *Item) FormatCheck() error {
	return nil
}

// UnmarshalJSON : (generated from github.com/podhmo/strangejson/_examples/manypackages04/item.Item)
func (x *Item) UnmarshalJSON(b []byte) error {
	type internal struct {
		Name *string
	}

	var p internal
	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	var merr *multierror.Error
	if p.Name == nil {
		merr = multierror.Append(merr, errors.New("Name is required"))
	} else {
		x.Name = *p.Name
	}
	return multierror.Append(merr, x.FormatCheck()).ErrorOrNil()
}
